package main

import "net"
import "fmt"
import "os"



func main() {

	listen, err := net.ListenPacket("udp", ":30000")

	if err!= nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
	buffer:= make([]byte, 1024)
	listen.ReadFrom(buffer)

	fmt.Printf("IP?: %s \n" , buffer)
	defer listen.Close()
}
