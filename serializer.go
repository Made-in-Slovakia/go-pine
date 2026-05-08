package pine

import (
	"bytes"
	"encoding/binary"
)

// TODO: Little and Big Endian should be configurable because it is emulator dependent.

// Serializes given input to []byte.
func fromUint8(input uint8) []byte {
	return []byte{input}
}

// Serializes given input to []byte.
func fromUint16(input uint16) []byte {
	bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(bytes, input)
	return bytes
}

// Serializes given input to []byte.
func fromUint32(input uint32) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, input)
	return bytes
}

// Serializes given input to []byte.
func fromUint64(input uint64) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, input)
	return bytes
}

func toBytes(commands []Command) ([]byte, error) {
	if len(commands) == 0 {
		return []byte{4, 0, 0, 0}, nil
	}

	// First 4 bytes are size, including these 4 bytes that defines size.
	messageSize := 4
	for _, c := range commands {
		messageSize += c.len()
	}

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, uint32(messageSize))

	for _, c := range commands {
		err := binary.Write(buf, binary.LittleEndian, c.opCode)
		if err != nil {
			return nil, err
		}
		if c.argument != nil {
			err := binary.Write(buf, binary.LittleEndian, c.argument)
			if err != nil {
				return nil, err
			}
		}
	}

	return buf.Bytes(), nil
}
