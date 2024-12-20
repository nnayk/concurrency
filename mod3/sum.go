package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

const NUM_SPLITS = 4

func merge_sort(nums []int) ([]int,error) {
	size := len(nums)
	if size <= 1 {
		return nums, nil
	} else {
		mid := size/2
		left,_ := merge_sort(append([]int{},nums[0:mid]...))
		right,_ := merge_sort(append([]int{},nums[mid:size]...))
		return merge(left,right,nums)
	}
}

func merge(left []int,right []int,nums []int) ([] int, error) {
	// fmt.Printf("gonna merge %v and %v\n",left,right)
	left_ind := 0
	right_ind := 0
	left_size := len(left)
	right_size := len(right)
	if(left_size+right_size > len(nums)) {
		return nums, fmt.Errorf("Expected left size = %v and right size = %v to sum up to nums size = %v",left_size,right_size,len(nums))
	}
	index := 0
	for ; left_ind < left_size && right_ind < right_size; {
		if(left[left_ind] <= right[right_ind]) {
			nums[index] = left[left_ind]
			left_ind++
		} else {
			nums[index] = right[right_ind]
			right_ind++	
		}
		index++
	}
	for ; left_ind < left_size; {
		nums[index] = left[left_ind]
		index += 1
		left_ind += 1
	}

	for ; right_ind < right_size; {
		nums[index] = right[right_ind]
		index += 1
		right_ind += 1
	}
	return nums, nil
}

func sort_slice(ch chan []int, nums []int, id int) {
	fmt.Printf("goroutine %v: going to sort array %v\n",id,nums)
	merge_sort(nums)
	ch <- nums
	fmt.Printf("goroutine %v: finished sorting array %v\n",id,nums)
}

func main() {
	args := os.Args
	fmt.Println("Enter a list of numbers (should be a multiple of 4) separated by spaces:")
	ch := make(chan []int,NUM_SPLITS)
	var scanner *bufio.Scanner
    // Create a new scanner
	if len(args) == 2 {
		fmt.Printf("file = %v\n",args[1])
		file, err := os.Open(args[1])
        if err != nil {
            fmt.Println("Error opening file:", err)
            return
        }
        defer file.Close() // Ensure the file is closed after reading
		scanner = bufio.NewScanner(file)
	} else {
		scanner = bufio.NewScanner(os.Stdin)
	}

	// Read the input from the user
	scanner.Scan()
	input := scanner.Text()
	fmt.Printf("input = %v\n",input)

	// Split the input into a slice of strings
	numStrings := strings.Fields(input)
	total_size := len(numStrings)
	// Assert the number of numbers is a multiple of 4
	// if total_size%NUM_SPLITS != 0 {
	// 	fmt.Printf("Error: The number of numbers entered must be a multiple of 4 (received %v numbers).",len(numStrings))
	// 	return
	// }
	var numbers []int
	// Convert the slice of strings to a slice of integers
	for _, numStr := range numStrings {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			fmt.Printf("Invalid number: %s\n", numStr)
			return
		}
		numbers = append(numbers, num)
	}
	
	// Experiment 1: Sort the entire array
	numbers_copy := make([]int,len(numbers))
	copy(numbers_copy,numbers)
	start := time.Now()
	merge_sort(numbers_copy)
	end := time.Since(start)
	fmt.Printf("Sorting entire array took %v\n",end)
	// fmt.Printf("numbers_copy = %v\n",numbers_copy)

	// Experiment 2: Divide the sorting work
	// invoke goroutines each on 1/4 of the array
	// perform merge step of merge sort to update array in place s.t. it's sorted
	split_size := int(math.Ceil(float64(total_size)/float64(NUM_SPLITS)))
	if split_size == 0 {
		split_size = 1
	}
	start = time.Now()
	var low,high int
	for i := 0; i < NUM_SPLITS; i++ {
		low = i*(split_size)
		if i== NUM_SPLITS-1 {
			high = total_size
		} else {
			high = (i+1)*(split_size)
		}
		go sort_slice(ch,numbers[low:high],i)
	}

	subarrays := make([][]int,NUM_SPLITS)
	for i := 0; i < NUM_SPLITS; i++ {
		subarrays[i] = <- ch
		// fmt.Printf("Received sub array = %v\n",subarrays[i])
	}
	fmt.Printf("subarrays = %v\n",subarrays)
	// merge
	num_groups := NUM_SPLITS
	// group_size := split_size
	var new_subarrays [][]int
	for ; num_groups > 1 ; {
		num_groups /= 2
		// group_size = group_size * 2
		new_subarrays = make([][]int,num_groups)
		// fmt.Printf("num_groups = %v, group_size = %v\n",num_groups,group_size)
		for i := 0; i < num_groups; i++ {
			new_subarrays[i] = make([]int,len(subarrays[i*2])+len(subarrays[i*2+1]))
			fmt.Printf("invoking merge for %v and %v\n",subarrays[i*2],subarrays[i*2+1])
			_, err := merge(subarrays[i*2],subarrays[i*2+1],new_subarrays[i])
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			fmt.Printf("sorted subarray = %v\n",new_subarrays)
		}
		subarrays = new_subarrays
		fmt.Printf("new subarrays = %v\n",subarrays)
	}
	end = time.Since(start)
	fmt.Printf("Dividing sort work took %v\n",end)
	fmt.Printf("Sorted array = %v\n",new_subarrays[0])
}