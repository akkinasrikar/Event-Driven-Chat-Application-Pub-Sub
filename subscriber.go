package pubsub

import "sync"

type Subscribers map[string][]*Subscriber

type Subscriber struct {
	Name      string
	Message   chan *Message
	Topics    map[string]bool
	Lock      sync.RWMutex
	Destroyed bool
}

func NewSubscriber(name string) *Subscriber {
	return &Subscriber{
		Name:    name,
		Message: make(chan *Message),
		Topics:  make(map[string]bool),
		Lock:    sync.RWMutex{},
	}
}

func (s *Subscriber) AddTopic(topic string) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.Topics[topic] = true
}

func (s *Subscriber) RemoveTopic(topic string) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	delete(s.Topics, topic)
}

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
func (s *Subscriber) GetMessage() *Message {
	return <-s.Message
}

// send message to subscriber
func (s *Subscriber) Signal(msg *Message) *Subscriber {
	s.Lock.RLock()
	defer s.Lock.RUnlock()
	if _, ok := s.Topics[msg.Topic]; ok {
		s.Message <- msg
	}
	return s
}

func (s *Subscriber) Destroy() {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	if !s.Destroyed {
		close(s.Message)
		s.Destroyed = true
	}
}

