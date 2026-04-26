package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open("messages.txt")

	if err != nil {
		log.Fatal(err)
	}

	str := ""

	for {
		data := make([]byte, 8)

		n, err := f.Read(data)

		if err == io.EOF {
			break
		}

		data = data[:n]

		if i := bytes.IndexByte(data, '\n'); i != -1 {
			str += string(data[:i])
			data = data[i+1:]
			fmt.Printf("read: %s\n", str)
			str = ""
		}

		str += string(data)
	}

	if len(str) != 0 {
		fmt.Printf("read: %s\n", str)
	}
}
