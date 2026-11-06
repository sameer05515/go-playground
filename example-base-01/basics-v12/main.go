package main

import "fmt"

func main() {
	for i := 1; i < 5; i++ {
		fmt.Println("Number is:", i)
	}

	counter := 0
	for {
		fmt.Println("Infinite loop")
		counter++
		if counter == 4 {
			break
		}
	}

	numbers := []int{11, 22, 33, 44, 55, 66}
	for index, value := range numbers {
		fmt.Printf("Index == %d, value == %d\n", index, value)
	}

	data := "Hello, world!"
	fmt.Println("-----------------------")
	for index, char := range data {
		fmt.Printf("Index == %d, value == %c\n", index, char)
	}

}
