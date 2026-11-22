//go:build darwin
// +build darwin

package pine

import (
	"errors"
	"net"
	"os"
)

type socket struct {
	slot       int
	connection net.Conn
}

func newSocket(slot int) *socket {
	return &socket{slot: slot}
}

func (s *socket) connect() error {
	socket_name := os.Getenv("TMPDIR")
	if socket_name == "" {
		socket_name = "/tmp"
	}
	socket_name += "/pcsx2.sock"

	connection, err := net.Dial("unix", socket_name)
	if err != nil {
		return err
	}
	s.connection = connection

	return nil
}

func (s *socket) readBytes(recvSize int) ([]byte, error) {
	if s.connection == nil {
		return nil, errors.New("not connected")
	}

	buffer := make([]byte, recvSize)
	recvLen, err := s.connection.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer[:recvLen], nil
}

func (s *socket) writeBytes(bytes []byte) (int, error) {
	if s.connection == nil {
		return -1, errors.New("not connected")
	}

	return s.connection.Write(bytes)
}

func (s *socket) disconnect() error {
	if s.connection != nil {
		err := s.connection.Close()
		s.connection = nil
		return err
	}

	return nil
}
