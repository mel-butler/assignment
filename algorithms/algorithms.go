package algorithms

import (
	"runtime"
	"sort"
	"sync"
)

// ye olde mergesort with concurrency
func ConcurrentMergeSort(nums []int) []int {
	// found out that using a power of two for the number of threads
	// is good for workload distribution for your CPU cores
	// because mergesort splits arrays in half over and over,
	// this will let you split work evenly between cores
	threads := previousPowerOfTwo(runtime.NumCPU())
	// segmentation of the array, based on the array length
	//and number of threads
	segmentSize := (len(nums) + threads - 1) / threads

	// baby's first channel. the buffer size is the thread
	ch := make(chan []int, threads)

	// waitgroup to sync when the goroutines finish
	var wg sync.WaitGroup
	// this loop uses goroutines to sort the segments concurrently
	for i := 0; i < threads; i++ {
		start := i * segmentSize
		end := (i + 1) * segmentSize
		if end > len(nums) { // last segment for elements that don't fit.
			end = len(nums) // this is to avoid out of bounds errors
		}

		wg.Add(1)                 // wait group counter, to track the goroutines
		go func(start, end int) { // groutine for the segments
			defer wg.Done()
			part := nums[start:end]
			sort.Ints(part) // actual sorting here
			ch <- part      // sorted is sent to the channel, it'll be collected later
		}(start, end)
	}

	go func() { // goroutine to close the channel when the other go stuff is done
		wg.Wait()
		close(ch)
	}()

	var segments [][]int // slice for sorted segments from the channel
	for part := range ch {
		segments = append(segments, part)
	}

	return mergeSortedSegments(segments) // we merging
}

// bucket sort, make it concurrent
func ConcurrentBucketSort(nums []int, numCores int) []int {

	// splitting array
	half := len(nums) / 2
	ch := make(chan []int, 2) // channel for the two halves

	go func() {
		// sorting first half
		sortedHalf := BucketSort(nums[:half])
		ch <- sortedHalf
	}()

	go func() {
		// sorting second half
		sortedHalf := BucketSort(nums[half:])
		ch <- sortedHalf
	}()

	// receiving sorted halves from channels
	sortedFirstHalf := <-ch
	sortedSecondHalf := <-ch

	// merging sorted halves
	merged := merge(sortedFirstHalf, sortedSecondHalf)

	// a new channel created with a buffer size of
	// the number of cores.
	chMerge := make(chan []int, numCores)

	// loop to sort the segments concurrently.
	// a goroutine is launched for each core
	// each goroutine sorts a segment using insertion sort
	// partially sorted segments are sent to the channel
	for i := 0; i < numCores; i++ {
		go func(i int) {
			start := i * len(merged) / numCores
			end := (i + 1) * len(merged) / numCores
			chMerge <- insertionSort(merged[start:end])
		}(i)
	}

	// now we just merge them into sorted.
	// it iterates thru the cores, receiving
	// segments from the channel
	var sorted []int
	for i := 0; i < numCores; i++ {
		sorted = merge(sorted, <-chMerge)
	}

	return sorted
}

func ConcurrentQuickSort(nums []int, numCores int) []int {

	//we need a pivot that can balance the partitions
	// we find a pivot using the median of the first, middle and
	// last elements of the array
	pivotIdx := medianOfThree(nums)
	pivot := nums[pivotIdx]

	// elements less than the pivot are swapped to the left
	nums[pivotIdx], nums[len(nums)-1] = nums[len(nums)-1], nums[pivotIdx]

	// left pointer initialized to 0.
	// it keeps track of the elements less than the pivot
	left := 0
	for i := 0; i < len(nums)-1; i++ {
		if nums[i] < pivot {
			// if element is less than pivot, it is swapped to the left
			nums[left], nums[i] = nums[i], nums[left]
			left++
		}
	}

	// swaps the element left is pointing at with the pivot.
	// left is pointing at the first element greater than the pivot
	nums[left], nums[len(nums)-1] = nums[len(nums)-1], nums[left]

	// channels for the left and right partitions
	chLeft := make(chan struct{})
	chRight := make(chan struct{})

	// waitgroup to sync when the goroutines finish
	var wg sync.WaitGroup

	// this continues recirsively sorting the partitions
	// until the left partition is sorted
	// chleft is closed on completion
	if left > 1 {
		// Limit number of concurrent goroutines
		if numCores > 0 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				ConcurrentQuickSort(nums[:left], numCores-1)
				close(chLeft)
			}()
		} else {
			close(chLeft)
		}
	} else {
		close(chLeft)
	}

	// same but for the right partition
	if len(nums)-left-1 > 1 {
		if numCores > 0 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				ConcurrentQuickSort(nums[left+1:], numCores-1)
				close(chRight)
			}()
		} else {
			close(chRight)
		}
	} else {
		close(chRight)
	}

	// Wait for goroutines to finish
	wg.Wait()

	// sorted array!
	return nums
}

// just an insertion sort. nothing special
func insertionSort(nums []int) []int {
	for i := 0; i < len(nums); i++ {
		temp := nums[i]
		j := i - 1
		for ; j >= 0 && nums[j] > temp; j-- {
			nums[j+1] = nums[j]
		}
		nums[j+1] = temp
	}
	return nums
}

// a simple bucket sort using insertion sort
func BucketSort(nums []int) []int {
	bucketSize := 10
	var max, min int
	for _, n := range nums {
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
	}
	nBuckets := int(max-min)/bucketSize + 1
	buckets := make([][]int, nBuckets)
	for i := 0; i < nBuckets; i++ {
		buckets[i] = make([]int, 0)
	}

	for _, n := range nums {
		idx := int(n-min) / bucketSize
		buckets[idx] = append(buckets[idx], n)
	}

	sorted := make([]int, 0)
	for _, bucket := range buckets {
		if len(bucket) > 0 {
			sorted = append(sorted, insertionSort(bucket)...)
		}
	}
	return sorted
}

/*func merge(left, right []int) []int {
	return append(left, right...)
}
*/

// found out it was a lot faster to create my own merge function.
func merge(arr1, arr2 []int) []int {
	// checks if they're empty
	if len(arr1) == 0 {
		return arr2
	}
	if len(arr2) == 0 {
		return arr1
	}

	// initializing slice first to make things speedy
	merged := make([]int, len(arr1)+len(arr2))
	i, j, k := 0, 0, 0

	// compares the elemments and merges them, incrementing
	for i < len(arr1) && j < len(arr2) {
		if arr1[i] <= arr2[j] {
			merged[k] = arr1[i]
			i++
		} else {
			merged[k] = arr2[j]
			j++
		}
		k++
	}

	// this handles any remaining elements
	for i < len(arr1) {
		merged[k] = arr1[i]
		i++
		k++
	}

	for j < len(arr2) {
		merged[k] = arr2[j]
		j++
		k++
	}

	return merged[:k] // truncate the merged slice to actual size
}

// merge sorted segments into a single sorted array.
// it takes a slice of ints, and returns a single Sorted slice of ints
func mergeSortedSegments(segments [][]int) []int {
	var merged []int // empty slice that will hold the merged segments

	// iterate through each segment and merge them
	for _, segment := range segments {
		merged = merge(merged, segment)
	}

	return merged
}

// finds the biggest power of two that is less than or equal to x
func previousPowerOfTwo(x int) int {
	// this is to handle big values,,, and stops overflow
	num := uint64(x)
	// if it is less than 3, it is already a power of two, or less
	// then it converts it back to an int
	if num < 3 {
		return int(num)
	}
	// using OR operations to set all the bits to
	// the right of the highest bit set in num
	// shifts one bit to the right & OR with the num
	num |= num >> 1
	// shifts two bits to the right. This continues doubling the
	// number of bits.
	num |= num >> 2
	num |= num >> 4
	num |= num >> 8
	num |= num >> 16
	num |= num >> 32
	// subtracts num >> 1 from num to get the closest
	// power of two, that is less than or equal to x
	// It gets rid of the bits on the right of the highest set
	// bit, and leaves the high bits to the left.
	// then it gets converted to an int!!
	return int(num - (num >> 1))
}

//	evaluates three elements from the given slice and determines the
//
// median among them by comparing their values and positions in the array.
func medianOfThree(nums []int) int {
	first, middle, last := nums[0], nums[len(nums)/2], nums[len(nums)-1] // gets three elements from a slice
	if (first <= middle && middle <= last) || (last <= middle && middle <= first) {
		return len(nums) / 2 // middle element is median
	} else if (middle <= first && first <= last) || (last <= first && first <= middle) {
		return 0 // first element is median
	}
	return len(nums) - 1 // last element is median
}
