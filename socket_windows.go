package pine

import (
	"errors"
	"net"
	"strconv"
)

type socket struct {
	port       int
	connection net.Conn
}

func newSocket(port int) *socket {
	return &socket{port: port}
}

func (s *socket) connect() error {
	connection, err := net.Dial("tcp", "127.0.0.1:" + strconv.Itoa(s.port))
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
