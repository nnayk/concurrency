package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)


func sorter(nums []int, id int) {
	fmt.Printf("goroutine %v: going to sort array %v",label,nums)
	fmt.Printf("goroutine %v: finished sorting array %v",label,nums)
}

func main() {
	fmt.Println("Enter a list of numbers (should be a multiple of 4) separated by spaces:")

    // Create a new scanner
	scanner := bufio.NewScanner(os.Stdin)

	// Read the input from the user
	scanner.Scan()
	input := scanner.Text()

	// Split the input into a slice of strings
	numStrings := strings.Fields(input)
	
	// Assert the number of numbers is a multiple of 4
	if (len(numStrings))%4 != 0 {
		fmt.Printf("Error: The number of numbers entered must be a multiple of 4 (received %v numbers).",len(numStrings))
		return
	}

	// Convert the slice of strings to a slice of integers
	var numbers []int
	for _, numStr := range numStrings {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			fmt.Printf("Invalid number: %s\n", numStr)
			return
		}
		numbers = append(numbers, num)
	}
	// invoke goroutines each on 1/4 of the array
	// perform merge step of merge sort to update array in place s.t. it's sorted
    fmt.Printf("Numbers entered: %v\n", numbers)
}