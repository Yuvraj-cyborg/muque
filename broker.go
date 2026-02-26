package main

import (
	"fmt"
	"log"
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
