package main

import (
	"fmt"
	"os"
	"reflect"
	"sync"
	"time"
)

// Type defs
type Chopstick struct {
	taken sync.Mutex
}

type Philosopher struct {
	id int
	left *Chopstick
	right *Chopstick
	eating bool // false by default
	eat_count int 
}

// Constants
const DEBUG = 1
const NUM_PHILOSOPHERS = 5
const NUM_CHOPSTICKS = NUM_PHILOSOPHERS
const DINE_COUNT = 3
const CONCURRENT_DINE = 2
const ACCEPT = 1
const REJECT = 0

func main() {
	// init chopsticks, philosophers, other data
	var wg sync.WaitGroup
	wg.Add(NUM_PHILOSOPHERS)
	chopsticks := make([]*Chopstick,NUM_CHOPSTICKS) 
	for i:=0;i<NUM_CHOPSTICKS;i++{
		chopsticks[i] = new(Chopstick)
	}
	// launch host goroutine
	phils := make([]*Philosopher,NUM_PHILOSOPHERS)
	channels := make([]chan int,NUM_PHILOSOPHERS)
	for i:=0;i<NUM_PHILOSOPHERS;i++{
		phils[i] = new(Philosopher)
		phils[i].left = chopsticks[i]
		phils[i].right = chopsticks[(i+1)%NUM_CHOPSTICKS]
		phils[i].id = i
		channels[i] = make(chan int)
		// launch philosopher
		go eat(&wg,phils[i],channels[i])
	}
	go host(phils,channels)
	wg.Wait()
}

func host(phils []*Philosopher,channels []chan int) {
	remaining := NUM_PHILOSOPHERS
	for {
		cases := make([]reflect.SelectCase, len(channels))
		for i, ch := range channels {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
		}
		chosen, id_val, _ := reflect.Select(cases)
		chosen = int(chosen)
		id := int(id_val.Int())
		if chosen != id {
			// fmt.Printf("Discrepancy between channel id (%v) and philosopher id (%v)\n",chosen,id)
		}
		phil := phils[id]
		// fmt.Printf("Received eat request from philosopher %v\n",phil.id)
		if(phil.eating) {
			// fmt.Printf("Error: Philosopher %v requested to eat but he is already eating!\n",phil.id)
			os.Exit(1)
		} else if(phil.eat_count==DINE_COUNT) {
			// fmt.Printf("Error: Philosopher %v requested to eat but he has already eaten %v times!\n",phil.id,DINE_COUNT)
			os.Exit(1)
		} else {
			var eating_phils []*Philosopher
			for _,p := range eating_phils {
				if(p.eating) {
					eating_phils = append(eating_phils, p)
				}
			}
			if(len(eating_phils) > 2) {
				// fmt.Printf("Uh oh, >2 philosophers are eating (eating philosophers: %v)",eating_phils)
				os.Exit(1)
			} else if(len(eating_phils)==2) {
				// fmt.Printf("Sorry, %v philosophers (%v and %v) are already eating!\n",CONCURRENT_DINE,eating_phils[0],eating_phils[1])
			} else {
				// if a neighbor is eating then reject the request
				left_index := id-1
				if id==0 {
					left_index=4
				}
				if(phils[(id+1)%NUM_PHILOSOPHERS].eating || phils[left_index].eating) {
					// fmt.Printf("Sorry, a neigbor is already eating!\n")
					channels[id] <- REJECT
				} else {
					if phil.eat_count+1 == DINE_COUNT {
						remaining -= 1
					}
					channels[id] <- ACCEPT
					if(remaining==0) {
						fmt.Printf("Everyone finished eating!\n")
						break
					}
				}
				
			}
		}
	}
}

func eat(wg *sync.WaitGroup,phil *Philosopher,ch chan int) {
	defer wg.Done()
	for {
		// fmt.Printf("Philosopher %v sent a request\n",phil.id)
		ch <- phil.id
		res := <- ch
		if res == ACCEPT {
			// fmt.Printf("Philosopher %v is approved to eat!\n",phil.id)
			// lock chopsticks
			phil.left.taken.Lock()
			phil.right.taken.Lock()
			fmt.Printf("starting to eat %v\n",phil.id)
			phil.eating = true
			phil.eat_count += 1
			time.Sleep(2*time.Second)
			fmt.Printf("finishing eating %v\n",phil.id)
			fmt.Printf("eat count for philosopher %v = %v\n",phil.id,phil.eat_count)
			// unlock chopsticks
			phil.left.taken.Unlock()
			phil.right.taken.Unlock()
			phil.eating = false
			if(phil.eat_count==DINE_COUNT) {
				break
			}
		} else {
			fmt.Printf("Philosopher %v is rejected to eat!\n",phil.id)
		}	
	}
}