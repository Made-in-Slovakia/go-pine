package pine

import (
	"bytes"
	"errors"
)

// Safe call to Buffer.Next(n) with check if there is enough bytes to read.
func nextBytes(buffer *bytes.Buffer, n int) ([]byte, error) {
	if n == 0 {
		return nil, nil
	}
	if buffer.Len() < n {
		return nil, errors.New("not enough bytes in buffer to read")
	}
	return buffer.Next(n), nil
}

// Read next bytes in Buffer as string. Size of string is defined by first 4 bytes.
func nextString(buffer *bytes.Buffer) ([]byte, error) {
	if buffer.Len() < 4 {
		return nil, errors.New("not enough bytes in buffer to read")
	}
	size, err := toUint32(buffer.Next(4))
	if err != nil {
		return nil, err
	}
	return nextBytes(buffer, int(size))
}
