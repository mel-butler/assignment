package main

import (
	"assignment/algorithms"
	"fmt"
	"math/rand"
	"testing"
)

/*
func BenchmarkBucketSort(b *testing.B) {
	inputSize := []int{10, 100, 1000, 10000, 100000}
	for _, size := range inputSize {
		b.Run(fmt.Sprintf("input_size_%d", size), func(b *testing.B) {
			testList := make([]int, size)
			for i := 0; i < size; i++ {
				testList[i] = rand.Intn(size)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				algorithms.BucketSort(testList)
			}
		})
	}
}
*/

func BenchmarkConcurrentBucketSort(b *testing.B) {
	inputSize := []int{10, 100, 1000, 10000, 100000} // the input size
	cores := []int{1, 2, 4, 8}                       // the number of cores

	for _, size := range inputSize { // checks input size
		for _, numCores := range cores {
			b.Run(fmt.Sprintf("input_size_%d_cores_%d", size, numCores), func(b *testing.B) { // runs the benchmark
				testList := make([]int, size) // making the list to test on
				for i := 0; i < size; i++ {
					testList[i] = rand.Intn(size)
				}
				b.ResetTimer() // timer starts here
				for i := 0; i < b.N; i++ {
					algorithms.ConcurrentBucketSort(testList, numCores) // tests algorithm/ passes the array to test on and the number of cores
				}
			})
		}
	}
}

func BenchmarkConcurrentQuickSort(b *testing.B) {
	inputSize := []int{10, 100, 1000, 10000, 100000}
	cores := []int{1, 2, 4, 8} // Define the number of cores

	for _, size := range inputSize {
		for _, numCores := range cores {
			b.Run(fmt.Sprintf("input_size_%d_cores_%d", size, numCores), func(b *testing.B) {
				testList := make([]int, size)
				for i := 0; i < size; i++ {
					testList[i] = rand.Intn(size)
				}
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					algorithms.ConcurrentQuickSort(testList, numCores)
				}
			})
		}
	}
}
