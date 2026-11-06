package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Create("example.txt")

	if err != nil {
		fmt.Println("Error while creating file", err)
		return
	}

	defer file.Close()

	content:= "First file created with Go"
	_, errors:= io.WriteString(file, content+"\n")
	if(errors!=nil){
		fmt.Println("Error while writing file", errors)
		return
	}

	fmt.Println("Successfully created file")
}
