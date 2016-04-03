package event
import (
	"errors"
	"fmt"
	"github.com/orivil/sorter"
)

var (
	ErrEventExist = errors.New("dispatcher.Dispatcher: event already exist")
	ErrEventNotExist = errors.New("dispatcher.Dispatcher: event not exist")
	ErrListenerExist = errors.New("dispatcher.Dispatcher: listener already exist")
	ErrListenerNotExist = errors.New("dispatcher.Dispatcher: listener not exist")
)

type Event struct {
	Name string
	Call func(listener interface{}, param ...interface{})
}

type Subscribe struct {
	Name string // event name
	Priority int
}

type Listener interface {
	GetSubscribe() (listenerName string, subscribes []Subscribe)
}

type Dispatcher struct {

	events map[string]*Event

	listeners map[string]Listener

	// {eventName: {listenerName: priority}}
	priorities map[string]map[string]int

	// {eventName: {listenerName}}
	sortedListeners map[string][]string
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		events: make(map[string]*Event, 10),
		listeners: make(map[string]Listener, 20),
		priorities: make(map[string]map[string]int, 10),
		sortedListeners: make(map[string][]string, 10),
	}
}

// AddEvent add an event to dispatcher
//
// error: if event exist return error 'ErrEventExist', else return nil
func (d *Dispatcher) AddEvent(e *Event) error {
	if _, ok := d.events[e.Name]; ok {
		return ErrEventExist
	} else {
		d.events[e.Name] = e
	}
	return nil
}

// AddEvent add events to dispatcher
//
// error: event may already exist
func (d *Dispatcher) AddEvents(es []*Event) error {
	for _, e := range es {
		err := d.AddEvent(e)
		if err == ErrEventExist {
			return fmt.Errorf("dispatcher.Dispatcher: add exist event %s", e.Name)
		}
	}
	return nil
}

// AddListener add event listeners
//
// error: if listener name exist return error 'ErrListenerExist'
func (d *Dispatcher) AddListener(ls ...Listener) error {
	for _, l := range ls {
		listenerName, subscribes := l.GetSubscribe()
		if _, ok := d.listeners[listenerName]; ok {
			return ErrListenerExist
		} else {
			// add instance
			d.listeners[listenerName] = l
		}

		// subscribe event
		for _, s := range subscribes {
			// initialize sorted listeners
			d.sortedListeners[s.Name] = nil

			// add priorities to priority container
			if d.priorities[s.Name] == nil {
				d.priorities[s.Name] = map[string]int{listenerName: s.Priority}
			} else {
				d.priorities[s.Name][listenerName] = s.Priority
			}
		}
	}
	return nil
}

// DelListener delete listener
//
// name: listener name which the 'GetSubscribe' function returned
//
// error: if listener not exist return error 'ErrListenerNotExist'
func (d *Dispatcher) DelListener(name string) error {
	listener := d.listeners[name]
	if listener == nil {
		return ErrListenerNotExist
	}

	// delete listener instance
	delete(d.listeners, name)

	// delete subscribes
	_, subscribes := listener.GetSubscribe()
	for _, s := range subscribes {
		delete(d.priorities[s.Name], name)
		// initialize sorted listeners
		d.sortedListeners[s.Name] = nil
	}
	return nil
}

// Trigger trigger an event with parameters
//
// error: if event not exist return error 'ErrEventNotExist'
func (d *Dispatcher) Trigger(event string, param ...interface{}) error {
	if e, ok := d.events[event]; ok {
		// read from cache
		listeners := d.sortedListeners[event]
		if listeners == nil {
			// reverse sort listeners
			sorter := sorter.NewPrioritySorter(d.priorities[event])
			listeners = sorter.SortReverse()
			d.sortedListeners[event] = listeners
		}
		for _, listener := range listeners {
			e.Call(d.listeners[listener], param...)
		}
	} else {
		return ErrEventNotExist
	}
	return nil
}
