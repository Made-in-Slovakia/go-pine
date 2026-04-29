package pine

import (
	"bytes"
	"errors"
	"fmt"
	"log"
)

// TODO evaluate possibility of adding TRACE log level or better log library
// Set to true to enable debugging messages in logs.
var DebugLogEnabled bool = false

// TODO rename this constant
// Maximum number of commands sent in a batch.
const MaxBatchSize = 50000

type Client struct {
	socket *socket
}

// Returns new PINE Client.
func NewClient(port int) *Client {
	c := &Client{}
	c.socket = &socket{port: port}
	return c
}

func (c *Client) Connect() error {
	return c.socket.connect()
}

func (c *Client) Disconnect() error {
	return c.socket.disconnect()
}

func (c *Client) Read32(addresses []uint32) ([]uint32, error) {
	if c.socket == nil {
		return nil, errors.New("client is not properly initialized")
	}

	commands := make([]Command, len(addresses))
	for i, a := range addresses {
		commands[i] = Read32Command(a)
	}

	answers, err := c.SendCommands(commands)
	if err != nil {
		return nil, err
	}

	results := make([]uint32, len(answers))
	for i, a := range answers {
		v, err := a.ContentAsUint32()
		if err != nil {
			return nil, err
		}
		results[i] = v
	}

	return results, nil
}

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

// TODO docs
func (c *Client) SendWriteStringCommands(address uint32, input string) (int, error) {
	return c.SendWriteBytesCommands(address, []byte(input))
}

// TODO docs
func (c *Client) SendWriteBytesCommands(address uint32, input []byte) (int, error) {
	// TODO it is possible ignore this as everything else has fallback to empty slice or nil?
	if len(input) == 0 {
		// We can safely ignore empty input and return command counter = 0.
		return 0, nil
	}

	commands := WriteBytesCommands(address, input)

	// Answers for write commands are always empty, we can ignore them.
	_, err := c.SendCommands(commands)
	return len(commands), err
}

// Sends given Command to connected emulator and returns Answer.
func (c *Client) SendCommand(command Command) (Answer, error) {
	answers, err := c.SendCommands([]Command{command})
	if err != nil {
		return Answer{}, err
	}
	return answers[0], nil
}

// Sends given Commands to connected emulator and returns Answers.
func (c *Client) SendCommands(commands []Command) ([]Answer, error) {
	if c.socket == nil {
		return nil, errors.New("client is not properly initialized")
	}

	// Input validations
	if commands == nil {
		// TODO test if it is even possible send nil
		// just defensive programming, this should never happen
		return nil, nil
	}
	// TODO test this, maybe nil slice will be better
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
		log.Printf("PINE request message created, message=[% X]", request)
	}

	// Write data
	b, err := c.socket.writeBytes(request)
	if err != nil {
		return nil, err
	}

	if DebugLogEnabled {
		log.Printf("PINE message sent, bytesSent=%d", b)
	}

	// Read data
	response, err := c.socket.readBytes(1024)
	if err != nil {
		return nil, err
	}

	if DebugLogEnabled {
		log.Printf("PINE response message received, message=[% X]", response)
	}

	if len(response) < 5 {
		return nil, errors.New("response is too short")
	}

	if response[4] == 0xFF {
		return nil, errors.New("error in response")
	}

	buf := bytes.NewBuffer(response)

	// First 4 bytes are response message size followed by one byte for result code.
	// We can skip those when reading this buffer.
	buf.Next(5)

	rawAnswers := make([][]byte, len(commands))
	for i, c := range commands {
		var rawAnswer []byte

		switch c.opCode {
		case msgRead8:
			rawAnswer, err = nextBytes(buf, 1)
		case msgRead16:
			rawAnswer, err = nextBytes(buf, 2)
		case msgRead32:
			rawAnswer, err = nextBytes(buf, 4)
		case msgRead64:
			rawAnswer, err = nextBytes(buf, 8)
		case msgWrite8, msgWrite16, msgWrite32, msgWrite64, msgSaveState, msgLoadState:
			// Empty response
			rawAnswer = nil
		case msgVersion, msgTitle, msgId, msgUuid, msgGameVersion:
			rawAnswer, err = nextString(buf)
		case msgStatus:
			rawAnswer, err = nextBytes(buf, 4)
		default:
			return nil, fmt.Errorf("unsupported opCode in response, opCode=%X", c.opCode)
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
		// TODO better message
		return nil, errors.New("wrong size of answers slice")
	}
	return answers, nil
}
