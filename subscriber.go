package pubsub

import (
	"fmt"
	"sync"
)

// Subscribers is a map of subscriber name to subscriber
type Subscribers map[string]*Subscriber

// Subscriber is a struct that holds the subscriber name, message channel and topics
type Subscriber struct {
	Name      string
	Message   chan *Message
	Topics    map[string]bool
	Lock      sync.RWMutex
	Destroyed bool
}

// NewSubscriber creates a new subscriber
func NewSubscriber(name string) *Subscriber {
	return &Subscriber{
		Name:    name,
		Message: make(chan *Message),
		Topics:  map[string]bool{},
		Lock:    sync.RWMutex{},
	}
}

// add topic to subscriber
func (s *Subscriber) AddTopic(topic string) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.Topics[topic] = true
}

// remove topic from subscriber
func (s *Subscriber) RemoveTopic(topic string) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	delete(s.Topics, topic)
}

// get topics from subscriber
func (s *Subscriber) GetTopics() []string {
	s.Lock.RLock()
	defer s.Lock.RUnlock()
	topics := make([]string, 0, len(s.Topics))
	for topic := range s.Topics {
		topics = append(topics, topic)
	}
	return topics
}

// get message from subscriber
func (s *Subscriber) GetMessage() <-chan *Message {
	return s.Message
}

// signal message to subscriber
func (s *Subscriber) Signal(msg *Message) *Subscriber {
	fmt.Printf("Delivered message to %v : %v \n", s.Name, msg.Payload)
	s.Lock.RLock()
	if !s.Destroyed {
		s.Message <- msg
	}
	s.Lock.RUnlock()
	return s
}

// destroy subscriber
func (s *Subscriber) Destroy() {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	if !s.Destroyed {
		close(s.Message)
		s.Destroyed = true
	}
}
