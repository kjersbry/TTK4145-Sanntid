package main

import (
	"./network/bcast"
	"./network/localip"
	//"./network/peers"
	//"flag"
	"fmt"
	"os"
	"time"
)
//import

//import

/*Move enum and struct to other module!!!*/
type Direction int
const (
	DOWN Direction = 0
	UP Direction = 1
	CAB Direction = 2
)

type Order struct { //default init?
	Order_ID int
	Floor int
	Direction Direction
	Is_served bool
	Cab_num int
}
func find_ip()string{

	localIP, err := localip.LocalIP()
	if err != nil {
	fmt.Println(err)
	localIP = "DISCONNECTED"
	}
id := fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
return id
}



func main() {
	orderTx := make(chan Order)
	orderRx := make(chan Order)
	helloOrder:=Order{666,3,UP,false,2}

	go bcast.Transmitter(16569, orderTx)
	go bcast.Receiver(16569, orderRx)
	time.Sleep(1 * time.Second)
	orderTx <- helloOrder
	fmt.Println("Started")
	for {
		select {

		case a := <-orderRx:
			fmt.Printf("Received: %#v\n", a)
		}
	}
}
