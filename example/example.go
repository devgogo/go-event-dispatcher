package main

import (
	"fmt"
	"github.com/wenmingtang/go-event-dispatcher/event"
)

const (
	UserCreated = "user.created"
)

type User struct {
	ID    int
	Login string
	Email string
}

type UserEvent struct {
	User User
}

type UserSubscriber struct {
}

func (s *UserSubscriber) SendActiveEmail(event event.Event, eventName string) {
	if evt, ok := event.(UserEvent); ok {
		fmt.Printf("send welcome email to %s\n", evt.User.Email)
	}
}

func (s *UserSubscriber) DoSomething(event event.Event, eventName string) {
	fmt.Printf("to do something\n")
}

func (s *UserSubscriber) SubscribedEvent() map[string][]event.Listener {
	return map[string][]event.Listener{
		UserCreated: {s.DoSomething, s.SendActiveEmail},
	}
}

func main() {
	d := event.NewDispatcher()

	user := User{ID: 1, Login: "twm", Email: "test@example.com"}
	evt := UserEvent{user}

	s := UserSubscriber{}

	d.AddSubscriber(&s)

	d.AddListener(UserCreated, func(event event.Event, eventName string) {
		fmt.Printf("event: %v\tname: %v\t\n", event, eventName)
	}, event.PriorityHigh)

	d.Dispatch(UserCreated, evt)
}
