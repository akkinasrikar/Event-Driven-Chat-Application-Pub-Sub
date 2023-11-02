package pubsub

type Publish struct {
	Sender   string `json:"sender"`
	Reciever string `json:"reciever"`
	Message  string `json:"message"`
}

type Join struct {
	UserName  string `json:"username"`
	GroupName string `json:"groupname"`
}

type Group struct {
	Name string `json:"name"`
}

type Broadcast struct {
	Sender  string `json:"sender"`
	Topic   string `json:"topic"`
	Message string `json:"message"`
}
