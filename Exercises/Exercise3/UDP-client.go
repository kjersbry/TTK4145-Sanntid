package main

import "net"
import "fmt"
import "os"

//SERVER IP: 10.100.23.242
const SERVER_IP = "10.100.23.242"


func UDPclient_serverdial(ip string, port int){
	addr := fmt.Sprintf("%s:%d", ip, port)
	sock, _ := net.Dial("udp", addr)
	sock.Write([]byte("Helloooo"))
}


func UDPclient_serverlisten(port int){
	portt := fmt.Sprintf(":%d", port)
	listen, err := net.ListenPacket("udp", portt)

	if err!= nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}

	buffer:= make([]byte, 1024)
	listen.ReadFrom(buffer)

	fmt.Printf("IP: %s \n" , buffer)
	defer listen.Close()

}

func main() {

	UDPclient_serverdial(SERVER_IP, 20006)
	UDPclient_serverlisten(20006)

}
