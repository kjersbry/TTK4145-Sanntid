package main

import (
	"encoding/gob"
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
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
/****************************************/
func UDP_send_order(ip string, port int, order Order){
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)

	err:= enc.Encode(Order{11,2,UP, false, 4})
	if err!= nil{
		log.Fatal("encode error: ", err)
	}

	addr := fmt.Sprintf("%s:%d", ip, port)
	sock, _ := net.Dial("udp", addr)
	msg, _:= bytes.ReadBytes(&buffer)
	sock.Write(msg)
}


func UDP_receive_order(port int){
	portt := fmt.Sprintf(":%d", port)
	listen, err := net.ListenPacket("udp", portt)

	if err!= nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}

	msg:= make([]byte, 1024)
	listen.ReadFrom(msg)

	var buffer bytes.Buffer
	buffer.Write(msg)

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

}
