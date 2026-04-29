//go:build aix || darwin || dragonfly || freebsd || (js && wasm) || linux || nacl || netbsd || openbsd || solaris
// +build aix darwin dragonfly freebsd js,wasm linux nacl netbsd openbsd solaris

package pine

import (
	"errors"
)

// TODO evalueate possibility of creating interface
// TODO implementation and tests

type socket struct {
}

func (s *socket) connect() error {
	return errors.New("not implemented yet")
}

func (s *socket) readBytes(recvSize int) ([]byte, error) {
	return nil, errors.New("not implemented yet")
}

func (s *socket) writeBytes(bytes []byte) (int, error) {
	return -1, errors.New("not implemented yet")
}

func (s *socket) disconnect() error {
	return errors.New("not implemented yet")
}
