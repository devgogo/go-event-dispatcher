package event

import (
	"sort"
	"sync"
)

const (
	PriorityDefault = 0
	PriorityLow     = -1
	PriorityHigh    = 100
)

type dispatcher struct {
	mux       sync.Mutex
	listeners map[string][]*listenerWrapper
}

func NewDispatcher() *dispatcher {
	return &dispatcher{
		listeners: make(map[string][]*listenerWrapper),
	}
}

type listenerWrapper struct {
	listener Listener
	priority int
}

type listenerSorter struct {
	listeners []*listenerWrapper
}

func (s *listenerSorter) Swap(i, j int) {
	s.listeners[i], s.listeners[j] = s.listeners[j], s.listeners[i]
}

func (s *listenerSorter) Less(i, j int) bool {
	return s.listeners[i].priority > s.listeners[j].priority
}

func (s *listenerSorter) Len() int {
	return len(s.listeners)
}

func (d *dispatcher) Dispatch(eventName string, event Event) {
	for _, listener := range d.SortedListeners(eventName) {
		listener.listener(event, eventName)
	}
}

func (d *dispatcher) AddListener(eventName string, l Listener, priority int) {
	d.mux.Lock()
	defer d.mux.Unlock()
	d.listeners[eventName] = append(d.listeners[eventName], &listenerWrapper{
		listener: l,
		priority: priority,
	})
}

func (d *dispatcher) RemoveListener(eventName string, listener Listener) {
	d.mux.Lock()
	defer d.mux.Unlock()
	delete(d.listeners, eventName)
}

func (d *dispatcher) SortedListeners(eventName string) []*listenerWrapper {
	d.mux.Lock()
	defer d.mux.Unlock()
	s := listenerSorter{d.listeners[eventName]}
	sort.Sort(&s)
	return d.listeners[eventName]
}

func (d *dispatcher) HasListeners(eventName string) bool {
	d.mux.Lock()
	defer d.mux.Unlock()
	return len(d.listeners[eventName]) > 0
}

func (d *dispatcher) AddSubscriber(subscriber Subscriber) {
	for eventName, listeners := range subscriber.SubscribedEvent() {
		for _, listener := range listeners {
			d.AddListener(eventName, listener, PriorityDefault)
		}
	}
}

func (d *dispatcher) RemoveSubscriber(subscriber Subscriber) {
	for eventName, listeners := range subscriber.SubscribedEvent() {
		for _, listener := range listeners {
			d.RemoveListener(eventName, listener)
		}
	}
}
