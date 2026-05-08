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

func toUint8(b []byte) (uint8, error) {
	if len(b) == 0 {
		return 0, fmt.Errorf("expected size 1 but is %v", len(b))
	}
	return b[0], nil
}

func toUint16(b []byte) (uint16, error) {
	if len(b) < 2 {
		return 0, fmt.Errorf("expected size 2 but is %v", len(b))
	}
	return binary.LittleEndian.Uint16(b), nil
}

func toUint32(b []byte) (uint32, error) {
	if len(b) < 4 {
		return 0, fmt.Errorf("expected size 4 but is %v", len(b))
	}
	return binary.LittleEndian.Uint32(b), nil
}

func toUint64(b []byte) (uint64, error) {
	if len(b) < 8 {
		return 0, fmt.Errorf("expected size 8 but is %v", len(b))
	}
	return binary.LittleEndian.Uint64(b), nil
}

// Returns the given bytes as string. Removes null-termination byte and all possible garbage after
// it.
//
// TODO: tests needed
func toString(b []byte) string {
	for i := 0; i < len(b); i++ {
		if b[i] == 0 {
			return string(b[:i])
		}
	}
	return string(b)
}
