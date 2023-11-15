package main

import (
	"assignment/bubbleSort"
	"assignment/numberGenerator"
	"assignment/quickSort"
	"fmt"
	"log"
	"time"
)

func main() {
	numberGenerator.CreateFile(-9999, 9999, 100)

	// starting timer
	startTime := time.Now()

	sortingAlgorithm("bubble", "in.csv")

	// ending timer
	endTime := time.Since(startTime)
	fmt.Println("Time taken to sort:", endTime)

}

func sortingAlgorithm(algorithm string, fileName string) []int {

	// Read the CSV file
	nums, err := numberGenerator.ReadFile(fileName)
	if err != nil {
		log.Fatal("error reading file:", err)
	}

	var sortedNums []int
	if algorithm == "bubble" {
		fmt.Println("sorting algorithm used: bubble sort")
		sortedNums = bubbleSort.BubbleSort(nums)
	} else if algorithm == "quick" {
		fmt.Println("sorting algorithm used: quick sort")
		sortedNums = quickSort.QuickSort(nums)
	} else {
		log.Fatal("did you spell it right?")
	}

	// Write the sorted numbers back to out.csv
	if err := numberGenerator.WriteFile("out.csv", sortedNums); err != nil {
		log.Fatal("error writing to file:", err)
	}

	fmt.Println("Numbers sorted.")
	return sortedNums
}
