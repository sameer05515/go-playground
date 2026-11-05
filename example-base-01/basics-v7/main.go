package main

import "fmt"

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("Denominator should not be zero")
	}

	return a/b, nil
}

func main() {

	data, _:= divide(100, 40)
	fmt.Println("Division is ", data)


	data, err := divide(10,0)
	if(err != nil){
		fmt.Println(err)
	}
	fmt.Println("Division is ", data)

}