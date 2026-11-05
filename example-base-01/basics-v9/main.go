package main

import "fmt"

func main() {
	numbers := []int{1, 2, 3, 4, 5}

	fmt.Println("Numbers :", numbers)
	fmt.Printf("Numbers has data type: %T\n", numbers)
	fmt.Println("Length: ", len(numbers))

	numbers = append(numbers, 11, 22, 33, 44, 55, 66, 77)
	fmt.Println("Numbers :", numbers)
	fmt.Println("Length: ", len(numbers))
	fmt.Println("Capacity :", cap(numbers))

	numbers2 := make([]int, 3, 5)

	fmt.Println("Numbers :", numbers2)
	fmt.Printf("Numbers has data type: %T\n", numbers2)
	fmt.Println("Length: ", len(numbers2))
	fmt.Println("Capacity :", cap(numbers2))


	var stocks= make([]string, 0)
	fmt.Println("---------------------------------")
	fmt.Println("Stocks slice :", stocks)
	fmt.Printf("Numbers has data type: %T\n", stocks)
	fmt.Println("Length: ", len(stocks))
	fmt.Println("Capacity :", cap(stocks))

}
