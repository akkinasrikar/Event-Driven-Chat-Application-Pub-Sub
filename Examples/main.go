package main

import (
	"fmt"
	"math/rand"
	"time"

	pubsub "akkina.com/pub-sub/app"
)

func main() {
	
	// Create a new message broker 
	broker := pubsub.NewMessageBroker()

	// Create a new subscriber and register it into our main broker
	subscriber := broker.Attach("Srikar")

	// Subscribe a subscriber to a topic
	broker.Subscribe(subscriber, "DSA")

	// Get the message channel from the subscriber
	ch := subscriber.GetMessage()


	// Go routines to send and receive messages
	go send(broker)
	go receive(subscriber.Name, ch)

	// Wait for user input to exit
	fmt.Scanln()
	fmt.Println("done")
}

func send(broker *pubsub.MessageBroker) {
	fmt.Println("Producer, sending messages...")
	for {
		msg := &pubsub.Message{
			Topic:   "DSA",
			Payload: fmt.Sprintf("%d", rand.Intn(100)),
		}
		broker.Broadcast(msg.Payload, msg.GetTopic())
		fmt.Printf("Producer, sent %v\n", msg.GetPayload().(string))
		time.Sleep(time.Second)
	}
}

func receive(name string, ch1 <-chan *pubsub.Message) {
	fmt.Printf("Subscriber %v, receiving...\n", name)
	for {
		msg, ok := <-ch1
		if !ok {
			fmt.Println("channel closed")
			continue
		}
		fmt.Printf("Subscriber %v, received %v\n", name, msg.GetPayload().(string))
		time.Sleep(time.Second)
	}
}
