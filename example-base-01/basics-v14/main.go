package main

import "fmt"

type Person struct {
	FirstName string
	LastName  string
	Age       int
}

type Contact struct{
	Email string
	Phone string
	Fax string
}

type Address struct{
	HouseNo int
	City string
	State string
}

type Employee struct{
	PersonDetails Person
	PersonContact Contact
	PersonAddress Address
}

func main() {
	var person1 Person
	fmt.Println("Person1 is: ", person1)

	person1.FirstName = "Premendra"
	person1.LastName = "Kumar"
	person1.Age = 40
	fmt.Println("Person1 is: ", person1)

	person2 := Person{
		FirstName: "Narendra",
		LastName:  "Karnik",
		Age:       37,
	}
	fmt.Println("Person2 is: ", person2)

	var person3 = new(Person)
	person3.FirstName = "Rashmi"
	person3.LastName = "Modi"
	person3.Age = 40
	fmt.Println("Person3 is: ", person3)

	fmt.Println("---------------------------------")

	// Method 4: Nested struct initialization
	employee1 := Employee{
		PersonDetails: Person{
			FirstName: "Arun",
			LastName:  "Kumar",
			Age:       35,
		},
		PersonContact: Contact{
			Email: "arun@example.com",
			Phone: "9876543210",
			Fax:   "123456",
		},
		PersonAddress: Address{
			HouseNo: 101,
			City:    "Faridabad",
			State:   "Haryana",
		},
	}

	fmt.Printf("Employee1: %+v\n", employee1)
	fmt.Printf("Employee Name: %s %s\n",
		employee1.PersonDetails.FirstName,
		employee1.PersonDetails.LastName)

}
