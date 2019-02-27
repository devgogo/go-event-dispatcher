package event

import (
	"reflect"
	"sort"
	"sync"
)

const (
	PriorityDefault = 0
	PriorityLow     = -1
	PriorityHigh    = 100
)

type dispatcher struct {
	mux       sync.RWMutex
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

func (d *dispatcher) Dispatch(eventName string, e Event) {
	for _, listener := range d.SortedListeners(eventName) {
		listener.listener(e, eventName)
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

func (d *dispatcher) RemoveListener(eventName string, l Listener) {
	d.mux.Lock()
	defer d.mux.Unlock()
	index := -1
	listeners := d.listeners[eventName]
	for i, l2 := range listeners {
		if reflect.ValueOf(l2.listener).Pointer() == reflect.ValueOf(l).Pointer() {
			index = i
			break
		}
	}
	if index == -1 {
		return
	}
	listeners[index] = listeners[len(listeners)-1]
	d.listeners[eventName] = listeners[:len(listeners)-1]
}

func (d *dispatcher) SortedListeners(eventName string) []*listenerWrapper {
	d.mux.Lock()
	defer d.mux.Unlock()
	s := listenerSorter{d.listeners[eventName]}
	sort.Sort(&s)
	return d.listeners[eventName]
}

func (d *dispatcher) HasListeners(eventName string) bool {
	d.mux.RLock()
	defer d.mux.RUnlock()
	return len(d.listeners[eventName]) > 0
}

func (d *dispatcher) AddSubscriber(s Subscriber) {
	for eventName, listeners := range s.SubscribedEvent() {
		for _, listener := range listeners {
			d.AddListener(eventName, listener, PriorityDefault)
		}
	}
}

func (d *dispatcher) RemoveSubscriber(s Subscriber) {
	for eventName, listeners := range s.SubscribedEvent() {
		for _, listener := range listeners {
			d.RemoveListener(eventName, listener)
		}
	}
}
