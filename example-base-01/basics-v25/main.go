package main

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	IsAdult bool   `json:"is_adult"`
}

func main() {
	fmt.Println("We are learning JSON")
	person := Person{Name: "Prem", Age: 40, IsAdult: true}
	fmt.Println("Person data: ", person)

	//convert person into JSON Encoding (Marshalling)
	jsondata, err := json.Marshal(person)
	if err!=nil {
		fmt.Println("Error marshalling", err)
		return
	}

	fmt.Println("Person data", string(jsondata))

	// Decoding (Unmarshalling)

	var personData Person
	errors := json.Unmarshal(jsondata, &personData)
	if errors!=nil{
		fmt.Println("Error unmarshalling", errors)
		return
	}

	fmt.Println("Person data after unmarshalling", personData)



}
