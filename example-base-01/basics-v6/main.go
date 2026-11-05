package main

import "fmt"

func simpleFunction() {
	fmt.Println("I am simple function")
}

func add(a int, b int) (result int) {
	result = a + b

	return
}

func main() {
	simpleFunction()

	fmt.Println("a + b =", 3, "+", 4, "=", add(3, 4))
}
