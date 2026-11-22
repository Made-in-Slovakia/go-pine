package pine

import (
	"bytes"
	"fmt"
)

// Command is a struct that represents single command that can be send to the emulator.
type Command struct {
	opCode   OpCode
	argument []byte
}

func (c *Command) len() int {
	if c.argument == nil {
		return 1
	}
	return len(c.argument) + 1
}

func (c Command) String() string {
	// Example: {opCode:F argument:[AB CD EF 00]}
	return fmt.Sprintf("{opCode:%X argument:[% X]}", c.opCode, c.argument)
}

// Returns a new [Command] with [OpCode] = [MsgRead8] and the given argument.
func Read8Command(address uint32) Command {
	return Command{
		opCode:   MsgRead8,
		argument: fromUint32(address),
	}
}

// Returns a new [Command] with [OpCode] = [MsgRead16] and the given argument.
func Read16Command(address uint32) Command {
	return Command{
		opCode:   MsgRead16,
		argument: fromUint32(address),
	}
}

// Returns a new [Command] with [OpCode] = [MsgRead32] and the given argument.
func Read32Command(address uint32) Command {
	return Command{
		opCode:   MsgRead32,
		argument: fromUint32(address),
	}
}

// Returns a new [Command] with [OpCode] = [MsgRead64] and the given argument.
func Read64Command(address uint32) Command {
	return Command{
		opCode:   MsgRead64,
		argument: fromUint32(address),
	}
}

// Returns a new [Command] with [OpCode] = [MsgWrite8] and the given arguments.
func Write8Command(address uint32, value uint8) Command {
	return Command{
		opCode:   MsgWrite8,
		argument: append(fromUint32(address), fromUint8(value)...),
	}
}

func write8CommandRaw(address uint32, value byte) Command {
	return Command{
		opCode:   MsgWrite8,
		argument: append(fromUint32(address), value),
	}
}

// Returns a new [Command] with [OpCode] = [MsgWrite16] and the given arguments.
func Write16Command(address uint32, value uint16) Command {
	return Command{
		opCode:   MsgWrite16,
		argument: append(fromUint32(address), fromUint16(value)...),
	}
}

func write16CommandRaw(address uint32, value1 byte, value2 byte) Command {
	return Command{
		opCode:   MsgWrite16,
		argument: append(fromUint32(address), value1, value2),
	}
}

// Returns a new [Command] with [OpCode] = [MsgWrite32] and the given arguments.
func Write32Command(address uint32, value uint32) Command {
	return Command{
		opCode:   MsgWrite32,
		argument: append(fromUint32(address), fromUint32(value)...),
	}
}

func write32CommandRaw(address uint32, value []byte) Command {
	return Command{
		opCode:   MsgWrite32,
		argument: append(fromUint32(address), value...),
	}
}

// Returns a new [Command] with [OpCode] = [MsgWrite64] and the given arguments.
func Write64Command(address uint32, value uint64) Command {
	return Command{
		opCode:   MsgWrite64,
		argument: append(fromUint32(address), fromUint64(value)...),
	}
}

// Returns a new [Command] with [OpCode] = [MsgVersion].
func VersionCommand() Command {
	return Command{
		opCode:   MsgVersion,
		argument: nil,
	}
}

// Returns a new [Command] with [OpCode] = [MsgSaveState] and the given argument.
func SaveStateCommand(saveState uint8) Command {
	return Command{
		opCode:   MsgSaveState,
		argument: fromUint8(saveState),
	}
}

// Returns a new [Command] with [OpCode] = [MsgLoadState] and the given argument.
func LoadStateCommand(saveState uint8) Command {
	return Command{
		opCode:   MsgLoadState,
		argument: fromUint8(saveState),
	}
}

// Returns a new [Command] with [OpCode] = [MsgTitle].
func TitleCommand() Command {
	return Command{
		opCode:   MsgTitle,
		argument: nil,
	}
}

// Returns a new [Command] with [OpCode] = [MsgId].
func IdCommand() Command {
	return Command{
		opCode:   MsgId,
		argument: nil,
	}
}

// Returns a new [Command] with [OpCode] = [MsgUuid].
func UuidCommand() Command {
	return Command{
		opCode:   MsgUuid,
		argument: nil,
	}
}

// Returns a new [Command] with [OpCode] = [MsgGameVersion].
func GameVersionCommand() Command {
	return Command{
		opCode:   MsgGameVersion,
		argument: nil,
	}
}

// Returns a new [Command] with [OpCode] = [MsgStatus].
func StatusCommand() Command {
	return Command{
		opCode:   MsgStatus,
		argument: nil,
	}
}

// Returns a new write [Command] that can be used to write the given bytes to the given address.
func WriteBytesCommands(address uint32, values []byte) []Command {
	if len(values) == 0 {
		return []Command{}
	}

	buffer := bytes.NewBuffer(values)
	// We create slice for commands with initial size calculated on how many Write32 commands we
	// will need. If input byte slice has size that can't be divided only to Write32 commands, we
	// will append rest as Write8 and Write16 commands at the end.
	//
	// TODO: evaluate possibility of refactoring to use Write64 instead of Write32.
	commands := make([]Command, (len(values) / 4))
	for i := 0; buffer.Len() > 0; i++ {
		bytes := buffer.Next(4)

		switch len(bytes) {
		case 3:
			commands = append(commands, write16CommandRaw(address, bytes[0], bytes[1]), write8CommandRaw(address+2, bytes[2]))
			address += 3
		case 2:
			commands = append(commands, write16CommandRaw(address, bytes[0], bytes[1]))
			address += 2
		case 1:
			commands = append(commands, write8CommandRaw(address, bytes[0]))
			address += 1
		case 0:
			// No more bytes left in buffer. We should never reach this code, it is just failsafe.
		case 4:
			fallthrough
		default:
			commands[i] = write32CommandRaw(address, bytes)
			address += 4
		}
	}

	return commands
}

// Returns a new [Command] with [OpCode] = [MsgUnimplemented].
func EmptyCommand() Command {
	return Command{
		opCode:   MsgUnimplemented,
		argument: nil,
	}
}
