package main

import (
	"assignment/algorithms"
	"assignment/numberGenerator"
	"fmt"
	"log"
	"runtime"
	"time"
)

func main() {
	// da file herself
	numberGenerator.CreateFile(-99999, 99999, 100000)

	// we reading the numbers from in.csv
	numbers, err := numberGenerator.ReadFile("in.csv")
	if err != nil {
		log.Fatalf("Error reading numbers: %v", err)
	}

	// starting timer
	startTime := time.Now()

	// algorithm time
	sortedNums := sortingAlgorithm("conquick", numbers)

	// ending timer
	elapsedTime := time.Since(startTime)

	// Write the sorted numbers to the output CSV file. Change the studentID to yours.
	if err := numberGenerator.WriteFile("out20259344.csv", sortedNums); err != nil {
		log.Fatalf("Error writing numbers: %v", err)
	}

	// Print the number of sorted elements and the time taken.
	fmt.Printf("Sorted %d numbers in %s.\n", len(numbers), elapsedTime)

}

func sortingAlgorithm(algorithm string, nums []int) []int {

	// a lovely if statement.
	var sortedNums []int
	if algorithm == "bucket" {
		sortedNums = algorithms.BucketSort(nums)
	} else if algorithm == "conbucket" {
		sortedNums = algorithms.ConcurrentBucketSort(nums, runtime.NumCPU())
	} else if algorithm == "conquick" {
		sortedNums = algorithms.ConcurrentQuickSort(nums, runtime.NumCPU())
	} else if algorithm == "conmerge" {
		sortedNums = algorithms.ConcurrentMergeSort(nums)
	} else {
		// I have mispelled a lot of things
		log.Fatal("did you spell it right?")
	}
	return sortedNums
}
