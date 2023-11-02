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
		broker.Subscribe(subscriber, sub.SubscribedTo)
		ch := subscriber.GetMessage()
		go receive(subscriber.Name, ch)
		c.JSON(200, gin.H{
			"message": fmt.Sprintf("%v subscribed to %v", sub.Subscriber, sub.SubscribedTo),
		})
	})

	router.POST("/send/", func(c *gin.Context) {
		var pub pubsub.Publish
		c.BindJSON(&pub)
		broker.Send(pub.Message, pub.Sender, pub.Reciever)
		c.JSON(200, gin.H{
			"message": fmt.Sprintf("%v sent a message to %v", pub.Sender, pub.Reciever),
		})
	})

	// create a Group
	router.POST("/group/", func(c *gin.Context) {
		var group pubsub.Group
		c.BindJSON(&group)
		broker.CreateTopic(group.Name)
		c.JSON(200, gin.H{
			"message": fmt.Sprintf("Group %v created", group.Name),
		})
	})

	// Join a group
	router.POST("/join/", func(c *gin.Context) {
		var sub pubsub.Join
		c.BindJSON(&sub)
		subscriber := broker.Attach(sub.UserName)
		broker.Subscribe(subscriber, sub.GroupName)
		ch := subscriber.GetMessage()
		go receive(subscriber.Name, ch)
		c.JSON(200, gin.H{
			"message": fmt.Sprintf("%v joined %v", sub.UserName, sub.GroupName),
		})
	})

	// publish to a group
	router.POST("/publish/topic/", func(c *gin.Context) {
		var pub pubsub.Broadcast
		c.BindJSON(&pub)
		fmt.Printf("%v sending message to %v", pub.Sender, pub.Topic)
		broker.Broadcast(pub.Message, pub.Sender, pub.Topic)
		c.JSON(200, gin.H{
			"message": fmt.Sprintf("%v published to %v", pub.Sender, pub.Topic),
		})
	})

	router.Run(":8080")
}

func receive(name string, ch <-chan *pubsub.Message) {
	for {
		_, ok := <-ch
		if !ok {
			fmt.Println("channel closed")
			continue
		}
		time.Sleep(time.Second)
	}
}
