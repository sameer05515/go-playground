package main

import "fmt"

func main() {
	x := 21
	if x > 5 {
		fmt.Println("X is greater than 5")
	} else {
		fmt.Println("X is less than 5")
	}

	z := 10
	if z > 10 {
		fmt.Println("z is greater than 10")
	} else if z > 5 {
		fmt.Println("z is greater 5 but less than or equal to 10")
	} else {
		fmt.Println("z is less than 5")
	}

	y := 10
	if x > 5 && (y > 5 || z < 43) {
		fmt.Println("Hi how are you?")
	} else {
		fmt.Println("Lets learn go lang")
	}
}
