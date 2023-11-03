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

	// unsubscribe a subscriber from a topic
	router.POST("/unsubscribe/", func(c *gin.Context) {
		var sub pubsub.Unsubscribe
		c.BindJSON(&sub)
		subscriber := broker.Attach(sub.Subscriber)
		broker.Unsubscribe(subscriber, sub.Unsubscribe)
		c.JSON(200, gin.H{
			"message": fmt.Sprintf("%v unsubscribed from %v", sub.Subscriber, sub.Unsubscribe),
		})
	})

	router.POST("/send/", func(c *gin.Context) {
		var pub pubsub.Publish
		c.BindJSON(&pub)
		err := broker.Send(pub.Message, pub.Sender, pub.Reciever)
		if err != nil {
			c.JSON(200, gin.H{
				"message": fmt.Sprintf("%v", err),
			})
			return
		}
		c.JSON(200, gin.H{
			"message": fmt.Sprintf("%v sent a message to %v", pub.Sender, pub.Reciever),
		})
	})

	// create a Group
	router.POST("/group/", func(c *gin.Context) {
		var group pubsub.Group
		c.BindJSON(&group)
		broker.CreateTopic(group.GroupName, group.Limit, group.Admin)
		c.JSON(200, gin.H{
			"message": fmt.Sprintf("%v created a group %v", group.Admin, group.GroupName),
		})
	})

	// Join a group
	router.POST("/join/", func(c *gin.Context) {
		var sub pubsub.Join
		c.BindJSON(&sub)
		subscriber := broker.Attach(sub.UserName)
		err := broker.SubscribeToGroup(subscriber, sub.GroupName, sub.Admin)
		if err != nil {
			c.JSON(200, gin.H{
				"message": fmt.Sprintf("%v", err),
			})
			return
		}
		ch := subscriber.GetMessage()
		go receive(subscriber.Name, ch)
		c.JSON(200, gin.H{
			"message": fmt.Sprintf("%v joined %v", sub.UserName, sub.GroupName),
		})
	})

	// Leave a group
	router.POST("/leave/", func(c *gin.Context) {
		var sub pubsub.Leave
		c.BindJSON(&sub)
		subscriber := broker.Attach(sub.UserName)
		err := broker.LeaveGroup(subscriber, sub.GroupName, sub.Admin)
		if err != nil {
			c.JSON(200, gin.H{
				"message": fmt.Sprintf("%v", err),
			})
			return
		}
		c.JSON(200, gin.H{
			"message": fmt.Sprintf("%v left %v", sub.UserName, sub.GroupName),
		})
	})

	// publish to a group
	router.POST("/publish/topic/", func(c *gin.Context) {
		var pub pubsub.Broadcast
		c.BindJSON(&pub)
		fmt.Printf("%v sending message to %v\n", pub.Sender, pub.Topic)
		err := broker.Broadcast(pub.Message, pub.Sender, pub.Topic)
		if err != nil {
			c.JSON(200, gin.H{
				"message": fmt.Sprintf("%v", err),
			})
			return
		}
		c.JSON(200, gin.H{
			"message": fmt.Sprintf("%v published to %v", pub.Sender, pub.Topic),
		})
	})

	// History of a Specific Group
	router.GET("/history/:group", func(c *gin.Context) {
		group := c.Param("group")
		history := broker.History[group]
		c.JSON(200, gin.H{
			"history": fmt.Sprintf("%v", history),
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
