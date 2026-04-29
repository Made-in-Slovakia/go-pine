package pine

import (
	"encoding/binary"
	"fmt"
)

// Deserializes given input using given opCode as identification of target type.
func deserialize(opCode opCode, input []byte) (any, error) {
	switch opCode {
	case msgRead8:
		return toUint8(input)
	case msgRead16:
		return toUint16(input)
	case msgRead32:
		return toUint32(input)
	case msgRead64:
		return toUint64(input)
	case msgWrite8:
		return "", nil
	case msgWrite16:
		return "", nil
	case msgWrite32:
		return "", nil
	case msgWrite64:
		return "", nil
	case msgVersion:
		return toString(input), nil
	case msgSaveState:
		return "", nil
	case msgLoadState:
		return "", nil
	case msgTitle:
		return toString(input), nil
	case msgId:
		return toString(input), nil
	case msgUuid:
		return toString(input), nil
	case msgGameVersion:
		return toString(input), nil
	case msgStatus:
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

func toString(b []byte) string {
	// Remove null-termination byte and all possible garbage after it.
	// TODO tests needed
	for i := 0; i < len(b); i++ {
		if b[i] == 0 {
			return string(b[:i])
		}
	}
	return string(b)
}
