package main

import "fmt"

func main() {
	age := 25
	name := "Premendra"
	height := 5.823478

	fmt.Println("age: ",age,"height: ",height,"name: ",name)

	fmt.Println("Hello world")

	fmt.Printf("Age is %d \t Height is %.4f \n", age, height)

	fmt.Printf("Type of age is %T , height is %T, name is %T", age, height, name)


}