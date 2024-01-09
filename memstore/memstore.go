package memstore

import (
	"sync"

	"github.com/sanity-io/litter"

	"github.com/rdnt/rdnt.dev/broker"
	"github.com/rdnt/rdnt.dev/event"
	"github.com/rdnt/rdnt.dev/uuid"
)

type Store struct {
	mux    sync.Mutex
	events []event.Event
	broker *broker.Broker[string, event.Event]
}

func (s *Store) Subscribe(handler func(e event.Event)) (dispose func(), err error) {
	return s.broker.Subscribe("", func(e event.Event) {
		handler(e)
	}), nil
}

func (s *Store) Events() ([]event.Event, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.events, nil
}

func (s *Store) Publish(e event.Event) error {
	s.mux.Lock()
	s.events = append(s.events, e)
	s.broker.Publish("", e)

	s.mux.Unlock()

	return nil
}

func (s *Store) String() string {
	s.mux.Lock()
	storedEvents := s.events
	s.mux.Unlock()

	type storedEvent struct {
		Type          event.Type
		AggregateType event.AggregateType
		AggregateId   uuid.UUID

		Event event.Event
	}

	var events []storedEvent

	for _, e := range storedEvents {
		events = append(events, storedEvent{
			Type:          e.Type(),
			AggregateType: e.AggregateType(),
			AggregateId:   e.AggregateId(),
			Event:         e,
		})
	}

	return litter.Sdump(events)
	//b, err := json.Marshal(events)
	//if err != nil {
	//	return "error"
	//}

	//return string(b)
}

func New() *Store {
	return &Store{
		events: make([]event.Event, 0),
		broker: broker.New[string, event.Event](),
	}
}
