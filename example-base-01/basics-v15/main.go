package main

import "fmt"

func modifyValueByReference(a *int) {
	*a = *a * 20
}

func main() {
	num := 2

	ptr := &num

	fmt.Println("Num has value:", num)
	fmt.Println("Pointer contains:", ptr)
	fmt.Println("Data contains through pointer:", *ptr)

	var pointer *int
	if pointer == nil {
		fmt.Println("Pointer is not assigned")
	}

	value := 10
	modifyValueByReference(&value)
	fmt.Println("Value contains:", value)
}
