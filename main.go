package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	string_list := make(chan string)

	go func() {
		defer f.Close()
		defer close(string_list)
		buffer := make([]byte, 8)
		var sb strings.Builder
		for {
			length, err := f.Read(buffer)

			if err != nil {
				if err == io.EOF {
					if sb.Len() != 0 {
						string_list <- sb.String()
					}
					break
				}
				panic("Issue reading bytes")
			}

			for ind := range length {
				val := buffer[ind]
				sb.WriteByte(val)
				if val == '\n' {
					string_list <- sb.String()
					sb.Reset()
				}
			}
		}
	}()

	return string_list

}

func main() {
	file, err := os.Open("messages.txt")

	if err != nil {
		panic("Error opening messages")
	}

	channel := getLinesChannel(file)

	for str := range channel {
		fmt.Printf("read: %s", str)
	}
}
