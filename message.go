package pubsub

type Message struct {
	Topic   string
	Payload interface{}
}

func (m *Message) GetTopic() string {
	return m.Topic
}

func (m *Message) GetPayload() interface{} {
	return m.Payload
}
