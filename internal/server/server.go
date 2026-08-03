package server

import (
	"fmt"
	"io"
	"net"
)

type Server struct {
	closed bool
}

func runConnection(s *Server, conn io.ReadWriteCloser) {
  out := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello World!")
  conn.Write(out)
  conn.Close()
}

func runServer(s *Server, listener net.Listener) {
	for {
		conn, err := listener.Accept()

		if s.closed {
			return
		}

		if err != nil {
			return
		}

		go runConnection(s, conn)
	}
}

func Serve(port uint16) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return nil, err
	}

	server := &Server{closed: false}

	go runServer(server, listener)

	return server, nil
}

func (s *Server) Close() error {
  s.closed = true
	return nil
}
