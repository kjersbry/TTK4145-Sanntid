package main

import (
	"encoding/gob"
	"bytes"
	"fmt"
	"log"
)
//import "net"

//import "os"

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
/****************************************/
func UDP_send_order(ip string, port int, order Order){
	var packet bytes.Buffer
	enc := gob.NewEncoder(&packet)

	err:= enc.Encode(Order{11,2,UP, false, 4})
	if err!= nil{
		log.Fatal("encode error: ", err)
	}



	addr := fmt.Sprintf("%s:%d", ip, port)
	sock, _ := net.Dial("udp", addr)
	sock.Write(packet)
}


func UDP_receive_order(port int){
	portt := fmt.Sprintf(":%d", port)
	listen, err := net.ListenPacket("udp", portt)

	if err!= nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}

	buffer:= make([]byte, 1024)
	listen.ReadFrom(buffer)

	dec := gob.NewDecoder(&buffer)
	var ord Order
	err = dec.Decode(&ord)
	if err!= nil{
		log.Fatal("encode error: ", err)
	}

	fmt.Printf("ID: %d, Floor: %d, Cab: %d\n", ord.Order_ID, ord.Floor, ord.Cab_num)


	defer listen.Close()

}







func main() {
	var packet bytes.Buffer
	enc := gob.NewEncoder(&network)

	err:= enc.Encode(Order{11,2,UP, false, 4})
	if err!= nil{
		log.Fatal("encode error: ", err)
	}


	sock.Write(packet)


	dec := gob.NewDecoder(&network)
	var ord Order
	err = dec.Decode(&ord)
	if err!= nil{
		log.Fatal("encode error: ", err)
	}

	fmt.Printf("ID: %d, Floor: %d, Cab: %d\n", ord.Order_ID, ord.Floor, ord.Cab_num)

}
