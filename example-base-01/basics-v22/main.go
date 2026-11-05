package main

import (
	"fmt"
	"os"
)

func main() {
	//content, err := ioutil.ReadFile("example.txt")
	content, err := os.ReadFile("example.txt") 

	if err != nil {
		fmt.Println("Error while reading file", err)
		return
	}

	fmt.Println(string(content))
}
