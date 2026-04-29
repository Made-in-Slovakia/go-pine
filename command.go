package pine

import (
	"bytes"
	"fmt"
)

type Command struct {
	opCode   opCode
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

// Returns new Command with opCode=msRead8 and given argument.
func Read8Command(address uint32) Command {
	return Command{
		opCode:   msgRead8,
		argument: fromUint32(address),
	}
}

// Returns new Command with opCode=msgRead16 and given argument.
func Read16Command(address uint32) Command {
	return Command{
		opCode:   msgRead16,
		argument: fromUint32(address),
	}
}

// Returns new Command with opCode=msgRead32 and given argument.
func Read32Command(address uint32) Command {
	return Command{
		opCode:   msgRead32,
		argument: fromUint32(address),
	}
}

// Returns new Command with opCode=msgRead64 and given argument.
func Read64Command(address uint32) Command {
	return Command{
		opCode:   msgRead64,
		argument: fromUint32(address),
	}
}

// Returns new Command with opCode=msgWrite8 and given arguments.
func Write8Command(address uint32, value uint8) Command {
	return Command{
		opCode:   msgWrite8,
		argument: append(fromUint32(address), fromUint8(value)...),
	}
}

// TODO
func write8CommandRaw(address uint32, value byte) Command {
	return Command{
		opCode:   msgWrite8,
		argument: append(fromUint32(address), value),
	}
}

// Returns new Command with opCode=msgWrite16 and given arguments.
func Write16Command(address uint32, value uint16) Command {
	return Command{
		opCode:   msgWrite16,
		argument: append(fromUint32(address), fromUint16(value)...),
	}
}

// TODO docs and rename
func write16CommandRaw(address uint32, value1 byte, value2 byte) Command {
	return Command{
		opCode:   msgWrite16,
		argument: append(fromUint32(address), value1, value2),
	}
}

// Returns new Command with opCode=msgWrite32 and given arguments.
func Write32Command(address uint32, value uint32) Command {
	return Command{
		opCode:   msgWrite32,
		argument: append(fromUint32(address), fromUint32(value)...),
	}
}

// TODO docs and rename
func write32CommandRaw(address uint32, value []byte) Command {
	// TODO check
	return Command{
		opCode:   msgWrite32,
		argument: append(fromUint32(address), value...),
	}
}

// Returns new Command with opCode=msgWrite64 and given arguments.
func Write64Command(address uint32, value uint64) Command {
	return Command{
		opCode:   msgWrite64,
		argument: append(fromUint32(address), fromUint64(value)...),
	}
}

// Returns new Command with opCode=msgVersion.
func VersionCommand() Command {
	return Command{
		opCode:   msgVersion,
		argument: nil,
	}
}

// Returns new Command with opCode=msgSaveState and given argument.
func SaveStateCommand(saveState uint8) Command {
	return Command{
		opCode:   msgSaveState,
		argument: fromUint8(saveState),
	}
}

// Returns new Command with opCode=msgLoadState and given argument.
func LoadStateCommand(saveState uint8) Command {
	return Command{
		opCode:   msgLoadState,
		argument: fromUint8(saveState),
	}
}

// Returns new Command with opCode=msgTitle.
func TitleCommand() Command {
	return Command{
		opCode:   msgTitle,
		argument: nil,
	}
}

// Returns new Command with opCode=msgId.
func IdCommand() Command {
	return Command{
		opCode:   msgId,
		argument: nil,
	}
}

// Returns new Command with opCode=msgUuid.
func UuidCommand() Command {
	return Command{
		opCode:   msgUuid,
		argument: nil,
	}
}

// Returns new Command with opCode=msgGameVersion.
func GameVersionCommand() Command {
	return Command{
		opCode:   msgGameVersion,
		argument: nil,
	}
}

// Returns new Command with opCode=msgStatus.
func StatusCommand() Command {
	return Command{
		opCode:   msgStatus,
		argument: nil,
	}
}

// TODO docs
func WriteBytesCommands(address uint32, values []byte) []Command {
	if len(values) == 0 {
		return nil
	}

	buf := bytes.NewBuffer(values)

	// TODO evaluate possibility of refactoring to use Write64.

	// We create slice for commands with initial size calculated on how many Write32 commands we
	// will need. If input byte slice has size that can't be divided only to Write32 commands, we
	// will append rest as Write8 and Write16 commands at the end.
	commands := make([]Command, (len(values) / 4))
	for i := 0; buf.Len() > 0; i++ {
		b := buf.Next(4)

		switch len(b) {
		case 3:
			commands = append(commands, write16CommandRaw(address, b[0], b[1]), write8CommandRaw(address+2, b[2]))
			address += 3
		case 2:
			commands = append(commands, write16CommandRaw(address, b[0], b[1]))
			address += 2
		case 1:
			commands = append(commands, write8CommandRaw(address, b[0]))
			address += 1
		case 0:
			// No more bytes left in buffer. We should never reach this code, it is just failsafe.
		case 4:
			fallthrough
		default:
			commands[i] = write32CommandRaw(address, b)
			address += 4
		}
	}

	return commands
}
