package pine

import (
	"bytes"
	"errors"
)

// Safe call to Buffer.Next(n) with check if there is enough bytes to read.
func nextBytes(b *bytes.Buffer, n int) ([]byte, error) {
	if n == 0 {
		return nil, nil
	}
	if b.Len() < n {
		return nil, errors.New("not enough bytes in buffer to read")
	}
	return b.Next(n), nil
}

// Read next bytes in Buffer as string. Size of string is defined by first 4 bytes.
func nextString(b *bytes.Buffer) ([]byte, error) {
	if b.Len() < 4 {
		return nil, errors.New("not enough bytes in buffer to read")
	}
	size, err := toUint32(b.Next(4))
	if err != nil {
		return nil, err
	}
	return nextBytes(b, int(size))
}
