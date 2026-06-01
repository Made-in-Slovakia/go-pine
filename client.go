package pine

import (
	"bytes"
	"errors"
	"fmt"
	"log"
)

// If set to true, library will print debugging messages to log.
//
// TODO: evaluate possibility of adding TRACE log level
// TODO: evaluate slog library
var DebugLogEnabled bool = false

// Maximum number of commands sent in a single batch.
const MaxBatchSize = 50000

// Main PINE client struct.
//
// Client is not thread-safe and does not provide any synchronization or locking.
type Client struct {
	socket *socket
}

// Returns a new PINE [Client].
func NewClient(port int) *Client {
	return &Client{
		socket: newSocket(port),
	}
}

// Connects to the emulator.
func (c *Client) Connect() error {
	if c.socket == nil {
		return errors.New("client is not properly initialized")
	}
	return c.socket.connect()
}

// Disconnects from the emulator.
func (c *Client) Disconnect() error {
	if c.socket == nil {
		return errors.New("client is not properly initialized")
	}
	return c.socket.disconnect()
}

// Reads 4 bytes (32 bits) from the given addresses and returns them as uint32.
func (c *Client) Read32(addresses []uint32) ([]uint32, error) {
	if c.socket == nil {
		return nil, errors.New("client is not properly initialized")
	}
	if len(addresses) == 0 {
		return []uint32{}, nil
	}

	commands := make([]Command, len(addresses))
	for i, address := range addresses {
		commands[i] = Read32Command(address)
	}

	answers, err := c.SendCommands(commands)
	if err != nil {
		return nil, err
	}

	results := make([]uint32, len(answers))
	for i, answer := range answers {
		value, err := answer.ContentAsUint32()
		if err != nil {
			return nil, err
		}
		results[i] = value
	}

	return results, nil
}

// Returns the emulator version.
//
// Content of version string is not specified in PINE standard, but in general, it should contain
// emulator name and version, example "PCSX2 v2.6.3".
func (c *Client) Version() (string, error) {
	if c.socket == nil {
		return "", errors.New("client is not properly initialized")
	}

	command := VersionCommand()

	if DebugLogEnabled {
		log.Printf("sending version command")
	}

	answer, err := c.SendCommand(command)
	if err != nil {
		return "", err
	}

	return answer.ContentAsString()
}

// Returns the emulator status. Possible values are 0:Running, 1:Paused, 2:Shutdown.
func (c *Client) Status() (uint32, error) {
	if c.socket == nil {
		return 0, errors.New("client is not properly initialized")
	}

	command := StatusCommand()

	if DebugLogEnabled {
		log.Printf("sending status command")
	}

	answer, err := c.SendCommand(command)
	if err != nil {
		return 0, err
	}
	return answer.ContentAsUint32()
}

// Sends multiple write [Command] needed to write the given string to the specified memory address.
// Returns the number of commands sent.
func (c *Client) SendWriteStringCommands(address uint32, text string) (int, error) {
	return c.SendWriteBytesCommands(address, []byte(text))
}

// Sends multiple write [Command] needed to write the given bytes to the specified memory address.
// Returns the number of commands sent.
func (c *Client) SendWriteBytesCommands(address uint32, input []byte) (int, error) {
	if len(input) == 0 {
		// We can safely ignore empty input and return command counter = 0.
		return 0, nil
	}

	commands := WriteBytesCommands(address, input)

	// Answers for write commands are always empty, we can ignore them.
	_, err := c.SendCommands(commands)
	return len(commands), err
}

// Sends the given [Command] to the connected emulator and returns [Answer].
func (c *Client) SendCommand(command Command) (Answer, error) {
	answers, err := c.SendCommands([]Command{command})
	if err != nil {
		return Answer{}, err
	}
	return answers[0], nil
}

// Sends given [Command] to connected emulator and returns [Answer].
func (c *Client) SendCommands(commands []Command) ([]Answer, error) {
	if c.socket == nil {
		return nil, errors.New("client is not properly initialized")
	}
	if len(commands) == 0 {
		return []Answer{}, nil
	}
	if len(commands) >= MaxBatchSize {
		return nil, errors.New("maximum command batch size is exceeded")
	}

	if DebugLogEnabled {
		log.Printf("sending commands, commands=%v", commands)
	}

	request, err := toBytes(commands)
	if err != nil {
		return nil, err
	}

	if DebugLogEnabled {
		log.Printf("request message created, message=[% X]", request)
	}

	// Write data
	bytesSent, err := c.socket.writeBytes(request)
	if err != nil {
		return nil, err
	}

	if DebugLogEnabled {
		log.Printf("message sent, bytesSent=%d", bytesSent)
	}

	// Read data
	// TODO: configurable size, or do it based on commands
	response, err := c.socket.readBytes(1024)
	if err != nil {
		return nil, err
	}

	if DebugLogEnabled {
		log.Printf("response message received, message=[% X]", response)
	}

	if len(response) < 5 {
		return nil, errors.New("response is too short")
	}

	if response[4] == 0xFF {
		return nil, errors.New("error in response")
	}

	buffer := bytes.NewBuffer(response)

	// First 4 bytes are response message size followed by one byte for result code We can skip
	// those when reading this buffer.
	//
	// TODO: we can create buffer with slice that has no first 5 bytes with 'response[5:]'
	buffer.Next(5)

	rawAnswers := make([][]byte, len(commands))
	for i, command := range commands {
		var rawAnswer []byte

		switch command.opCode {
		case MsgRead8:
			rawAnswer, err = nextBytes(buffer, 1)
		case MsgRead16:
			rawAnswer, err = nextBytes(buffer, 2)
		case MsgRead32:
			rawAnswer, err = nextBytes(buffer, 4)
		case MsgRead64:
			rawAnswer, err = nextBytes(buffer, 8)
		case MsgWrite8, MsgWrite16, MsgWrite32, MsgWrite64, MsgSaveState, MsgLoadState:
			// Empty response
			rawAnswer = nil
		case MsgVersion, MsgTitle, MsgId, MsgUuid, MsgGameVersion:
			rawAnswer, err = nextString(buffer)
		case MsgStatus:
			rawAnswer, err = nextBytes(buffer, 4)
		default:
			return nil, fmt.Errorf("unsupported opCode in response, opCode=%X", command.opCode)
		}

		if err != nil {
			return nil, err
		}
		rawAnswers[i] = rawAnswer
	}

	if DebugLogEnabled {
		log.Printf("response splitted to answers, rawAnswers=%v", rawAnswers)
	}

	answers := make([]Answer, len(rawAnswers))
	for i, rawAnswer := range rawAnswers {
		answer, err := deserialize(commands[i].opCode, rawAnswer)
		if err != nil {
			return nil, err
		}
		answers[i] = Answer{content: answer}
	}

	if DebugLogEnabled {
		log.Printf("raw answers converted, answers=%v", answers)
	}

	// This should never happen but we want be sure that caller always
	// receive array of same size as given commands array.
	if len(commands) != len(answers) {
		return nil, errors.New("number of commands and answers do not match")
	}
	return answers, nil
}
