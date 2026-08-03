package server

import (
	"bytes"
	"fmt"
	"io"
	"net"

	"httpfromtcp.kritishdhaubanjar.com.np/internal/request"
	"httpfromtcp.kritishdhaubanjar.com.np/internal/response"
)

type Server struct {
	closed  bool
	handler Handler
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func runConnection(s *Server, conn io.ReadWriteCloser) {
	defer conn.Close()

	headers := response.GetDefaultHeaders(0)
	r, err := request.RequestFromReader(conn)

	if err != nil {
		response.WriteStatusLine(conn, response.StatusBadRequest)
		response.WriteHeaders(conn, headers)
		return
	}

	writer := bytes.NewBuffer([]byte{})
	handlerError := s.handler(writer, r)

	var body []byte = nil
	var status response.StatusCode = response.StatusOK

	if handlerError != nil {
		status = handlerError.StatusCode
		body = []byte(handlerError.Message)
	} else {
		body = writer.Bytes()
	}

	/**/
	if r.RequestLine.RequestTarget == "/chunked-encoding" {
		headers.Delete("Content-Length")
		headers.Set("Transfer-Encoding", "chunked")

		response.WriteStatusLine(conn, status)
		response.WriteHeaders(conn, headers)

		var image = []string{
			"89504E470D0A1A0A0000000D4948445200000028000000280803000000BB20485F",
			"0000000373424954080808DBE14FE000000021504C5445A0C3FF4374E0A3C6FF",
			"3C70DF346BDE86AAF4789EEF6B93EB98BBFC4F7CE35D87E7D3105360000000AC",
			"49444154388DEDD2CB0E84200C0550E8CBB6FFFFC1D331C4805A6492597A176E",
			"3CB950A094373F05F63C33152332D1070A425823483295E16ACB54021F2E244F",
			"A0D52E96572AF590345F197B889C56AEC2E54658DD6301EFA14F0EB2AF9C140E",
			"274EF9CC6D1E8C89E23377210BBB99735978686BEFB1B90719FF95658B086B49",
			"31003B12B610651B05AE385EE1FD33071FEEAF59BB9C3A385EDD3767C8377D7B",
			"E7697590A410B717FE157E008FEF0524CD3DCFBA0000000049454E44AE426082",
		}

		for _, data := range image {
			conn.Write([]byte(fmt.Sprintf("%x\r\n", len(data))))
			conn.Write([]byte(data))
			conn.Write([]byte("\r\n"))
		}

		conn.Write([]byte("0\r\n\r\n"))

		return
	}

	/**/

	headers.Replace("Content-length", fmt.Sprintf("%d", len(body)))

	response.WriteStatusLine(conn, status)
	response.WriteHeaders(conn, headers)
	conn.Write(body)
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

func Serve(port uint16, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return nil, err
	}

	server := &Server{closed: false, handler: handler}

	go runServer(server, listener)

	return server, nil
}

func (s *Server) Close() error {
	s.closed = true
	return nil
}
