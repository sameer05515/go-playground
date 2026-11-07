package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	fmt.Println("Learning web service")

	res, err := http.Get("https://jsonplaceholder.typicode.com/posts/1")

	if(err!= nil){
		fmt.Println("Error while fetching get response", err)
		return
	}

	defer res.Body.Close()

	fmt.Printf("Type of response: %T\n", res)

	data, err := ioutil.ReadAll(res.Body)

	if(err!=nil){
		fmt.Println("Error occured while reading respones", err)
		return
	}

	fmt.Println("Response:", string(data))
}