package main

import "net"
import "fmt"

const SERVER_IP = "10.100.23.242"


func main() {

  addr := fmt.Sprintf("%s:%d", SERVER_IP, 33546)
  sock, _ := net.Dial("tcp", addr)
  msg:=fmt.Sprintf("%s\000", "Hola")
  sock.Write([]byte(msg))

  buffer:= make([]byte, 1024)
  sock.Read(buffer)
  fmt.Printf("Message: %s \n" , buffer)

  sock.Read(buffer)
	fmt.Printf("Message: %s \n" , buffer)
  }
