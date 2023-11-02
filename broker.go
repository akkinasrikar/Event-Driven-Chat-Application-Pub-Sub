package pubsub

import (
	"fmt"
	"sync"
)

// MessageBroker is the main broker that will handle the subscribers and topics
type MessageBroker struct {
	subscribers Subscribers
	lock        sync.RWMutex
	topics      map[string]Subscribers
}

// NewMessageBroker creates a new message broker
func NewMessageBroker() *MessageBroker {
	return &MessageBroker{
		subscribers: Subscribers{},
		lock:        sync.RWMutex{},
		topics:      map[string]Subscribers{},
	}
}

// Attach Create a new subscriber and register it into our main broker
func (b *MessageBroker) Attach(name string) *Subscriber {
	subscriber := NewSubscriber(name)
	b.lock.Lock()
	b.subscribers[subscriber.Name] = subscriber
	b.lock.Unlock()
	return subscriber
}

// subscribe a subscriber to a topic
func (b *MessageBroker) Subscribe(subscriber *Subscriber, topic string) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if _, ok := b.topics[topic]; !ok {
		b.topics[topic] = Subscribers{}
	}
	b.topics[topic][subscriber.Name] = subscriber
	subscriber.AddTopic(topic)
}

// unsubscribe a subscriber from a topic
func (b *MessageBroker) Unsubscribe(subscriber *Subscriber, topic string) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if _, ok := b.topics[topic]; ok {
		delete(b.topics[topic], subscriber.Name)
		subscriber.RemoveTopic(topic)
	}
}

// detach a subscriber from the broker
func (b *MessageBroker) Detach(subscriber *Subscriber) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if _, ok := b.subscribers[subscriber.Name]; ok {
		delete(b.subscribers, subscriber.Name)
	}
	for _, topic := range subscriber.GetTopics() {
		b.Unsubscribe(subscriber, topic)
	}
}

// get the number of subscribers for a topic
func (b *MessageBroker) Subscribers(topic string) int {
	b.lock.RLock()
	defer b.lock.RUnlock()
	return len(b.topics[topic])
}

// broadcast a message to all subscribers of a topic
func (b *MessageBroker) Broadcast(payload interface{}, topics ...string) {
	for _, topic := range topics {
		if b.Subscribers(topic) == 0 {
			continue
		}
		b.lock.RLock()
		for _, s := range b.topics[topic] {
			m := &Message{
				Topic:   topic,
				Payload: payload,
			}
			go (func(s *Subscriber) {
				s.Signal(m)
			})(s)
		}
		b.lock.RUnlock()
	}
}

// send a message to a specific subscriber
func (b *MessageBroker) Send(payload string, sender string, reciever string) {
	b.lock.RLock()
	defer b.lock.RUnlock()
	m := &Message{
		Topic:   reciever,
		Payload: payload,
	}
	// check if the reciever is subscribed to the sender
	if _, ok := b.topics[reciever][sender]; !ok {
		fmt.Println("reciever is not subscribed to sender")
		return
	}
	go (func(s *Subscriber) {
		s.Signal(m)
	})(b.subscribers[reciever])
}

func (b *MessageBroker) CreateTopic(name string) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if _, ok := b.topics[name]; !ok {
		b.topics[name] = Subscribers{}
	}
}
