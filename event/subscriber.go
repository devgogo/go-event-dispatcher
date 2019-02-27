package event

type Subscriber interface {
	SubscribedEvent() map[string][]Listener
}
