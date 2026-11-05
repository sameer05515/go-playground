package main

import "fmt"

func main() {
	var names [5]string

	names[0] = "Prem"
	names[2] = "Kumar"

	fmt.Println("Names of persons:", names)

	fmt.Printf("Names of persons: %q\n", names)
	fmt.Println("Length of names array:", len(names))

	var numbers= [8]int{1,2,3,4,5}

	fmt.Println("Number array elements:", numbers)
	fmt.Println("Length of numbers array:", len(numbers))

	fmt.Println("Value of 2nd index in numbers array:", numbers[1])
}