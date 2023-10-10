package pubsub

// Message is the struct that is passed to the subscribers
type Message struct {
	Topic   string
	Payload interface{}
}

// GetTopic returns the topic of the message
func (m *Message) GetTopic() string {
	return m.Topic
}

// GetPayload returns the payload of the message
func (m *Message) GetPayload() interface{} {
	return m.Payload
}
