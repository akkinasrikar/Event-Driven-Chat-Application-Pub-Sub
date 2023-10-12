package pubsub

type Publish struct {
	Sender   string `json:"sender"`
	Reciever string `json:"reciever"`
	Message  string `json:"message"`
}

type Subscribe struct {
	Subscriber string `json:"subscriber"`
	Topic      string `json:"topic"`
}
