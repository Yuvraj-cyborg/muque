package main

import (
	"fmt"
	"net"
)

func main() {
 	fmt.Println("Mini messeage queue : ");
	ln, err := net.Listen("tcp", ":8080");
	if err != nil {
		fmt.Printf("Couldn't listen sry !\n")
	}

	for {
	   conn, err := ln.Accept()
	   if err != nil {
		fmt.Printf("Couldn't connect sry\n")
	   }
	go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buff := make([]byte,1024)

	for {
		n, err := conn.Read(buff);
		if (err != nil) {
			fmt.Printf("client disconnected !\n")
			break
		}

		fmt.Printf("Received: %s\n", string(buff[:n])) // buff[:n] -> to not print toooo many 0s
	}
	
}
