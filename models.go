package pubsub

type Subscribe struct {
	Subscriber   string `json:"subscriber"`
	SubscribedTo string `json:"subscribed_to"`
}

type Unsubscribe struct {
	Subscriber  string `json:"subscriber"`
	Unsubscribe string `json:"unsubscribe"`
}

type Publish struct {
	Sender   string `json:"sender"`
	Reciever string `json:"reciever"`
	Message  string `json:"message"`
}

type Join struct {
	Admin     string `json:"admin"`
	UserName  string `json:"username"`
	GroupName string `json:"groupname"`
}

type Leave struct {
	Admin     string `json:"admin"`
	UserName  string `json:"username"`
	GroupName string `json:"groupname"`
}

type Group struct {
	Admin     string `json:"admin"`
	GroupName string `json:"groupname"`
	Limit     int    `json:"limit"`
}

type Broadcast struct {
	Sender  string `json:"sender"`
	Topic   string `json:"topic"`
	Message string `json:"message"`
}
