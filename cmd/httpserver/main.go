package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"httpfromtcp.kritishdhaubanjar.com.np/internal/request"
	"httpfromtcp.kritishdhaubanjar.com.np/internal/response"
	"httpfromtcp.kritishdhaubanjar.com.np/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, func(w io.Writer, req *request.Request) *server.HandlerError {
		if req.RequestLine.RequestTarget == "/yourproblem" {
			return &server.HandlerError{
				StatusCode: response.StatusBadRequest,
				Message:    "Whoopsie, your problem\n",
			}
		} else if req.RequestLine.RequestTarget == "/myproblem" {
			return &server.HandlerError{
				StatusCode: response.StatusInternalServerError,
				Message:    "Whoopsie, my bad\n",
			}
		} else {
			w.Write([]byte("All good\n"))
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	defer server.Close()

	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
