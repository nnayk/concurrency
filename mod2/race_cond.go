/*
EXPLANATION:
In this code the "populate" goroutine -- as the name suggests -- is responsible
for populating the array with non-zero values. The "output" goroutine then
prints the array contents. The race condition here is that there is no guarantee
that the "populate" goroutine will finish before the "output" goroutine starts,
and thus the "output" goroutine may print some or all array contents before it
is even populated. Thus, although the output is expected to be
*/

package main

import (
	"fmt"
	"sync"
)
var wg sync.WaitGroup
var ARRAY_SIZE = 100

func populate(nums []int) {
	defer wg.Done()
	for i, _ := range nums {
		nums[i] = i+1
	}
}

func output(nums []int) {
	defer wg.Done()
	fmt.Println("inside output")
	for _, num := range nums {
		fmt.Printf("%v ", num)
	}
}
func main() {
	wg.Add(2)
	// create an array
	nums := make([]int, ARRAY_SIZE*10)
	go populate(nums)
	go output(nums)
	wg.Wait()
}