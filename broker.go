package main

import (
	"fmt"
	"net"
	"sync"
)

type Message struct {
	Command string `json:"command"`
	Topic   string `json:"topic"`
	Payload string `json:"payload"`
}

// Subscriber wraps a connection with its own mutex to serialize writes.
type Subscriber struct {
	Conn net.Conn
	Mu   sync.Mutex
}

// WriteMsg sends a message to this subscriber, serialized via its mutex.
// Returns an error if the write fails (dead connection).
func (s *Subscriber) WriteMsg(data []byte) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	_, err := s.Conn.Write(data)
	return err
}

type Broker struct {
	// Topic -> Array of Subscribers (what I call the hash table of pub/sub)
	Subscribers map[string][]*Subscriber
	Lock        sync.RWMutex
}

func (b *Broker) Subscribe(topic string, sub *Subscriber) {
	b.Lock.Lock()
	b.Subscribers[topic] = append(b.Subscribers[topic], sub)
	b.Lock.Unlock()
	fmt.Printf("New subscriber added to topic: %s\n", topic)
}

func (b *Broker) Publish(topic string, payload string) {
	b.Lock.RLock()
	subscribers := b.Subscribers[topic]
	outgoingMsg := []byte(fmt.Sprintf("BROKER BROADCAST [%s]: %s\n", topic, payload))

	var dead []*Subscriber
	for _, sub := range subscribers {
		if err := sub.WriteMsg(outgoingMsg); err != nil {
			dead = append(dead, sub)
		}
	}
	b.Lock.RUnlock()

	// Clean up dead connections outside the read lock
	for _, sub := range dead {
		b.RemoveSubscriber(topic, sub)
	}

	fmt.Printf("Broadcasted to %d subscribers on topic: %s\n", len(subscribers), topic)
}

func (b *Broker) RemoveSubscriber(topic string, target *Subscriber) {
	b.Lock.Lock()
	defer b.Lock.Unlock()

	subscribers := b.Subscribers[topic]

	for i, sub := range subscribers {
		if sub == target {
			// Append everything BEFORE the index, with everything AFTER the index.
			b.Subscribers[topic] = append(subscribers[:i], subscribers[i+1:]...)

			fmt.Printf("Cleaned up dead socket from topic: %s\n", topic)
			break
		}
	}
}
