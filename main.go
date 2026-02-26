package main

import (
	"fmt"
	"net"
)

func main() {
 	fmt.Println("Mini messeage queue : ");
	ln, err := net.Listen("tcp", ":8080");
	if err != nil {
	// handle error
	}

	for {
	   conn, err := ln.Accept()
	   if err != nil {
		 	 fmt.Println("Couldn't connect sry")
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
			fmt.Println("client disconnected !")
			break
		}

		fmt.Printf("Received: %s\n", string(buff[:n])) // buff[:n] -> to not print toooo many 0s
	}
	
}
