package main

import (
	"fmt"
	"strings"
)

func main() {
	data := "apple,orange,banana,chiku"
	parts := strings.Split(data, ",")
	fmt.Println(parts)
	for index, value := range strings.Split(data, ",") {
		fmt.Printf("Index == %d, value == %s\n", index, value)
	}

	str := "one two three four two two five"
	count := strings.Count(str, "two")
	fmt.Println("Count:", count)

	str = "   Hello,              Go!  "
	trimmed := strings.TrimSpace(str)
	fmt.Println("trimmed:", trimmed)

	for index, value := range strings.Split(str, ",") {
		fmt.Printf("Index == %d, value == %s\n", index, strings.TrimSpace(value))
	}

	str1 := "Premendra"
	str2 := "Gupta"
	result := strings.Join([]string{str1, "Kumar", str2}, " ")
	fmt.Println("result:", result)

}
