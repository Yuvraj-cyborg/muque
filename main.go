package main

import (
	"fmt"
	"net"
	"encoding/json"
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
	myBroker := &Broker{
		Subscribers: make(map[string][]net.Conn),
	}

	go handleConnection(conn,myBroker)
  }
}

func handleConnection(conn net.Conn,b *Broker) {
	defer conn.Close()
	buff := make([]byte,1024)

	for {
		n, err := conn.Read(buff);
		if (err != nil) {
			fmt.Printf("client disconnected !\n")
			break
		}
		var msg Message
		err = json.Unmarshal(buff[:n],&msg)
		if err != nil {
			fmt.Printf("Invalid JSON received!\n")
            continue
		}
		// fmt.Printf("Received: %s\n", string(buff[:n])) // buff[:n] -> to not print toooo many 0s
		fmt.Printf("Parsed successfully -> Command: %s | Topic: %s | Payload: %s\n", msg.Command, msg.Topic, msg.Payload)

		if msg.Command == "SUB" {
			b.Subscribe(msg.Topic, conn)
		}else if msg.Command == "PUB" {
			b.Publish(msg.Topic, msg.Payload)
		} else {
			fmt.Println("Unknown command received.")
		}
	}
	
}
