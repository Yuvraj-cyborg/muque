package main

import (
	"net"
	"sync"
)

type Message struct {
	Command string `json:"command"`
	Topic   string `json:"topic"`
    Payload string `json:"payload"`
}

type Broker struct {
	// Topic -> Array of Sockets (what I call the hash table of pub/sub)
    Subscribers map[string][]net.Conn
	Lock sync.RWMutex
}

func (b *Broker) Subscribe(topic string, conn net.Conn) {
	b.Lock.Lock()
	b.Subscribers[topic] = append(b.Subscribers[topic], conn)
	b.Lock.Unlock()
	fmt.Printf("New subscriber added to topic: %s\n", topic)
}
