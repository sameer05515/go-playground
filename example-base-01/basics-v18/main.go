package main

import (
	"fmt"
	"time"
)

func main() {
	currentTime := time.Now()
	fmt.Println("Current time:", currentTime)
	fmt.Printf("Type of currentTime: %T\n", currentTime)

	fmt.Println("Formatted time: ", currentTime.Format("2006/01/02, Monday"))
	fmt.Println("Formatted time: ", currentTime.Format("02/01/2006, Monday"))

	layoutStr := "02/01/2006"
	dateStr := "25/11/2026"
	formattedTime, _ := time.Parse(layoutStr, dateStr)
	fmt.Println("Formatted time:", formattedTime)

	// Add 1 more day in currentTime
	newDate := currentTime.Add(24 * time.Hour)
	fmt.Println("newDate time", newDate)
	formattedNewTime := newDate.Format("2006/01/02, Monday")
	fmt.Println("formattedNewTime time", formattedNewTime)
}
