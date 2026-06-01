package pine

import (
	"encoding/binary"
	"fmt"
)

// Deserializes the given input using the given opCode as identification of the target type.
func deserialize(opCode OpCode, input []byte) (any, error) {
	switch opCode {
	case MsgRead8:
		return toUint8(input)
	case MsgRead16:
		return toUint16(input)
	case MsgRead32:
		return toUint32(input)
	case MsgRead64:
		return toUint64(input)
	case MsgWrite8:
		return "", nil
	case MsgWrite16:
		return "", nil
	case MsgWrite32:
		return "", nil
	case MsgWrite64:
		return "", nil
	case MsgVersion:
		return toString(input), nil
	case MsgSaveState:
		return "", nil
	case MsgLoadState:
		return "", nil
	case MsgTitle:
		return toString(input), nil
	case MsgId:
		return toString(input), nil
	case MsgUuid:
		return toString(input), nil
	case MsgGameVersion:
		return toString(input), nil
	case MsgStatus:
		return toUint32(input)
	default:
		return nil, fmt.Errorf("unsupported opCode in response, opCode=%X", opCode)
	}
}

func toUint8(bytes []byte) (uint8, error) {
	if len(bytes) == 0 {
		return 0, fmt.Errorf("expected size 1 but is %v", len(bytes))
	}
	return bytes[0], nil
}

func toUint16(bytes []byte) (uint16, error) {
	if len(bytes) < 2 {
		return 0, fmt.Errorf("expected size 2 but is %v", len(bytes))
	}
	return binary.LittleEndian.Uint16(bytes), nil
}

func toUint32(bytes []byte) (uint32, error) {
	if len(bytes) < 4 {
		return 0, fmt.Errorf("expected size 4 but is %v", len(bytes))
	}
	return binary.LittleEndian.Uint32(bytes), nil
}

func toUint64(bytes []byte) (uint64, error) {
	if len(bytes) < 8 {
		return 0, fmt.Errorf("expected size 8 but is %v", len(bytes))
	}
	return binary.LittleEndian.Uint64(bytes), nil
}

// Returns the given bytes as string. Removes null-termination byte and all possible garbage after
// it.
//
// TODO: tests needed
func toString(bytes []byte) string {
	for i := 0; i < len(bytes); i++ {
		if bytes[i] == 0 {
			return string(bytes[:i])
		}
	}
	return string(bytes)
}
