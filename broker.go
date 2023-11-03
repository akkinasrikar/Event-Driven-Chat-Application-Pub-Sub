package pubsub

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// MessageBroker is the main broker that will handle the subscribers and topic
type MessageBroker struct {
	subscribers Subscribers
	lock        sync.RWMutex
	topics      map[string]Subscribers
	topicLimit  map[string]int
	GroupAdmin  map[string][]string
	History     map[string][]MessagesHistory
}

type MessagesHistory struct {
	timeStamp string
	message   string
	sender    string
}

// NewMessageBroker creates a new message broker
func NewMessageBroker() *MessageBroker {
	return &MessageBroker{
		subscribers: Subscribers{},
		lock:        sync.RWMutex{},
		topics:      map[string]Subscribers{},
		topicLimit:  map[string]int{},
		GroupAdmin:  map[string][]string{},
		History:     map[string][]MessagesHistory{},
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

func (b *MessageBroker) SubscribeToGroup(subscriber *Subscriber, topic string, admin string) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	if _, ok := b.topics[topic]; !ok {
		b.topics[topic] = Subscribers{}
	}
	if len(b.topics[topic]) >= b.topicLimit[topic] {
		fmt.Printf("%v has reached the limit of %v subscribers", topic, b.topicLimit[topic])
		return errors.New("topic has reached the limit of subscribers")
	}
	if admin != "" {
		for _, adminUser := range b.GroupAdmin[topic] {
			if admin == adminUser {
				b.topics[topic][subscriber.Name] = subscriber
				subscriber.AddTopic(topic)
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
func (b *MessageBroker) Broadcast(payload string, sender string, topic string) error {
	messageObj := MessagesHistory{
		timeStamp: time.Now().Format("2006-01-02 15:04:05"),
		message:   payload,
		sender:    sender,
	}
	b.lock.RLock()
	if _, ok := b.topics[topic][sender]; !ok {
		fmt.Println("sender is not subscribed to the topic")
		return errors.New("sender is not subscribed to the topic")
	}
	b.History[topic] = append(b.History[topic], messageObj)
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
	return nil
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
	if _, ok := b.topics[GroupName]; !ok {
		b.topics[GroupName] = Subscribers{}
		b.topicLimit[GroupName] = Limit
		b.GroupAdmin[GroupName] = []string{Admin}
	}
	b.lock.Unlock()
	subscriber := b.Attach(Admin)
	b.SubscribeToGroup(subscriber, GroupName, Admin)
}


func (b *MessageBroker) GetHistory(topic string) []MessagesHistory {
	return b.History[topic]
}
