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

func (s *socket) connect() error {
	connection, err := net.Dial("tcp", "127.0.0.1:"+strconv.Itoa(s.port))
	if err != nil {
		return err
	}
	s.connection = connection

	return nil
}

func (s *socket) readBytes(recvSize int) ([]byte, error) {
	if s.connection == nil {
		// TODO better message
		return nil, errors.New("please call Connect first")
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
		// TODO better message
		return -1, errors.New("please call Connect first")
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
