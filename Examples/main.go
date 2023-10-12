package main

import (
	"fmt"
	"time"

	pubsub "akkina.com/pub-sub/app"
	"github.com/gin-gonic/gin"
)

func main() {
	// Create a new message broker
	broker := pubsub.NewMessageBroker()

	router := gin.Default()
	router.POST("/subscribe/", func(c *gin.Context) {
		var sub pubsub.Subscribe
		c.BindJSON(&sub)
		subscriber := broker.Attach(sub.Subscriber)
		broker.Subscribe(subscriber, sub.Topic)
		ch := subscriber.GetMessage()
		go receive(subscriber.Name, ch)
		c.JSON(200, gin.H{
			"message": fmt.Sprintf("%v subscribed to %v", sub.Subscriber, sub.Topic),
		})
	})
	router.POST("/publish/", func(c *gin.Context) {
		var pub pubsub.Publish
		c.BindJSON(&pub)

		broker.Broadcast(pub.Message, pub.Sender)
		c.JSON(200, gin.H{
			"message": fmt.Sprintf("%v published to %v", pub.Sender, pub.Reciever),
		})
	})

	router.Run(":8080")
}

func receive(name string, ch <-chan *pubsub.Message) {
	for {
		msg, ok := <-ch
		if !ok {
			fmt.Println("channel closed")
			continue
		}
		fmt.Printf("Subscriber %v, received %v\n", name, msg.GetPayload().(string))
		time.Sleep(time.Second)
	}
}
