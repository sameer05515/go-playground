package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("example.txt")

	if err != nil {
		fmt.Println("Error while creating file", err)
		return
	}

	defer file.Close()

	buffer := make([]byte, 1024)

	for {
		n, err := file.Read(buffer)

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("Error while reading file", err)
			return
		}

		fmt.Println(string(buffer[:n]))
	}
}
