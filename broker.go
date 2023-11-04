package pubsub

import (
	"errors"
	"fmt"
	"sync"
)

// MessageBroker is the main broker that will handle the subscribers and topics
type MessageBroker struct {
	subscribers Subscribers
	lock        sync.RWMutex
	topics     map[string]Subscribers
	topicLimit map[string]int
	GroupAdmin map[string][]string

}

// NewMessageBroker creates a new message broker
func NewMessageBroker() *MessageBroker {
	return &MessageBroker{
		subscribers: Subscribers{},
		lock:        sync.RWMutex{},
		topics:      map[string]Subscribers{},
		topicLimit:  map[string]int{},
		GroupAdmin:  map[string][]string{},
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
}

func (b *MessageBroker) SubscribeToGroup(subscriber *Subscriber, sub Join) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	if _, ok := b.topics[sub.GroupName]; !ok {
		b.topics[sub.GroupName] = Subscribers{}
	}
	if len(b.topics[sub.GroupName]) >= b.topicLimit[sub.GroupName] {
		fmt.Printf("%v has reached the limit of %v subscribers", sub.GroupName, b.topicLimit[sub.GroupName])
		return errors.New("topic has reached the limit of subscribers")
	}
	if sub.Admin != "" {
		for _, adminUser := range b.GroupAdmin[sub.GroupName] {
			if sub.Admin == adminUser {
				b.topics[sub.GroupName][subscriber.Name] = subscriber
				subscriber.AddTopic(sub.GroupName)
				if sub.MakeAdmin {
					b.GroupAdmin[sub.GroupName] = append(b.GroupAdmin[sub.GroupName], subscriber.Name)
				}
				return nil
			}
		}
	}
	fmt.Println("admin is not owner of the group")
	return errors.New("admin is not owner of the group")
}

// leave a group
func (b *MessageBroker) LeaveGroup(subscriber *Subscriber, topic string, admin string) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	if admin != "" {
		for _, adminUser := range b.GroupAdmin[topic] {
			if admin == adminUser {
				if _, ok := b.topics[topic][subscriber.Name]; !ok {
					fmt.Println("subscriber is not subscribed to the topic")
					return errors.New("subscriber is not subscribed to the topic")
				}
				delete(b.topics[topic], subscriber.Name)
				return nil
			}
		}
	}
	fmt.Println("admin is not owner of the group")
	return errors.New("admin is not owner of the group")
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
func (b *MessageBroker) Send(payload string, sender string, reciever string) error {
	b.lock.RLock()
	defer b.lock.RUnlock()
	m := &Message{
		Topic:   reciever,
		Payload: payload,
	}
	// check if the reciever is subscribed to the sender
	if _, ok := b.topics[reciever][sender]; !ok {
		fmt.Println("reciever is not subscribed to sender")
		return errors.New("reciever is not subscribed to sender")
	}
	go (func(s *Subscriber) {
		s.Signal(m)
	})(b.subscribers[reciever])
	return nil
}

func (b *MessageBroker) CreateTopic(GroupName string, Limit int, Admin string) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if _, ok := b.topics[GroupName]; !ok {
		b.topics[GroupName] = Subscribers{}
		b.topicLimit[GroupName] = Limit
		b.GroupAdmin[GroupName] = []string{Admin}

	}
}
