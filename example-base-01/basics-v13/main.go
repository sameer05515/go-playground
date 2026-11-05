package main

import "fmt"

func main() {
	studentGrades := make(map[string]int)

	studentGrades["Prem"] = 34
	studentGrades["Naren"] = 99
	studentGrades["Bob"] = 56
	studentGrades["Charlie"] = 67

	fmt.Println("Marks of Prem:", studentGrades["Prem"])
	studentGrades["Prem"] = 100
	fmt.Println("Marks of Prem:", studentGrades["Prem"])

	delete(studentGrades, "Bob")
	fmt.Println("Marks of Bob:", studentGrades["Bob"])

	//check if a given key exists
	grade, exists := studentGrades["Bob"]
	fmt.Println("Grade of Bob:", grade)
	fmt.Println("Bob exists:", exists)

	for key, value := range studentGrades {
		fmt.Printf("Key == %s, value == %d\n", key, value)
	}

	// Another way to declare and intialize map
	country := map[string]string{
		"India": "New Delhi",
		"Nepal": "Kathmandu",
	}

	for key, value := range country {
		fmt.Printf("Key == %s, value == %s\n", key, value)
	}

}
