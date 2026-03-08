package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

func main() {
	fmt.Println("Mini messeage queue : ")
	myBroker := &Broker{
		Subscribers: make(map[string][]net.Conn),
	}
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Couldn't connect sry\n")
		}

		go handleConnection(conn, myBroker)
	}
}

func handleConnection(conn net.Conn, b *Broker) {
	defer conn.Close()
	var myTopics []string

	defer func() {
		for _, topic := range myTopics {
			b.RemoveSubscriber(topic, conn)
		}
	}()

	buff := make([]byte, 1024)

	for {
		n, err := conn.Read(buff)
		if err != nil {
			fmt.Printf("client disconnected !\n")
			break
		}
		var msg Message
		err = json.Unmarshal(buff[:n], &msg)
		if err != nil {
			fmt.Printf("Invalid JSON received!\n")
			continue
		}
		// fmt.Printf("Received: %s\n", string(buff[:n])) // buff[:n] -> to not print toooo many 0s
		fmt.Printf("Parsed successfully -> Command: %s | Topic: %s | Payload: %s\n", msg.Command, msg.Topic, msg.Payload)

		switch msg.Command {
		case "SUB":
			b.Subscribe(msg.Topic, conn)
			myTopics = append(myTopics, msg.Topic)
		case "PUB":
			b.Publish(msg.Topic, msg.Payload)
		default:
			fmt.Println("Unknown command received.")
		}
	}

}
