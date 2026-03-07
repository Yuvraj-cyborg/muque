package main

import (
	"net"
	"sync"
	"fmt"
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

func (b *Broker) Publish(topic string, payload string) {
	b.Lock.RLock()
	defer b.Lock.RUnlock()
	subscribers := b.Subscribers[topic]
	for _, conn := range subscribers {
        // Format the string and cast it to a byte array so the socket can send it
        outgoingMsg := fmt.Sprintf("BROKER BROADCAST [%s]: %s\n", topic, payload)
        conn.Write([]byte(outgoingMsg))
    }

	fmt.Printf("Broadcasted to %d subscribers on topic: %s\n", len(subscribers), topic)
}
