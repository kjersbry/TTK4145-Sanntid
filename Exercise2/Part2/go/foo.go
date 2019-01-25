// Use `go run foo.go` to run your program

package main

import (
    . "fmt"
    "runtime"
)

// Control signals
const (
	GetNumber = iota
	Exit
)

func number_server(add_number <-chan int, control <-chan int, number chan<- int) {
	var i = 0
	for {
		select {
		//receive different messages and handle them correctly
		case receive:= <- add_number: 
			i += receive

		case receive:= <- control:
			if receive == GetNumber {
				number <- i //send number
			} else if receive == Exit {
				return;
			}
		}
	}
}

func incrementing(add_number chan<-int, finished chan<- bool) {
	for j := 0; j<1000000; j++ {
		add_number <- 1
	}
	//signal that the goroutine is finished
	finished <- true
}

func decrementing(add_number chan<- int, finished chan<- bool) {
	for j := 0; j<1000002; j++ {
		add_number <- -1
	}
	//signal that the goroutine is finished
	finished <- true
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Construct the required channels
	// Think about wether the receptions of the number should be unbuffered, or buffered with a fixed queue size.
	add_number := make(chan int)
	finished := make(chan bool) 
	control := make(chan int)
	number := make(chan int)

	// Spawn the required goroutines
	go incrementing(add_number, finished)
	go decrementing(add_number, finished)
	go number_server(add_number, control, number)

	// Block on finished from both "worker" goroutines
	var fin = 0
	for fin < 2 {
		select {
		case <- finished:
			fin++
		}
	}
	
	control<-GetNumber
	Println("The magic number is:", <- number)
	control<-Exit
}
