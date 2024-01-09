package commands

import (
	"github.com/rdnt/rdnt.dev/event"
	"github.com/rdnt/rdnt.dev/project"
)

type EventStore interface {
	Publish(event event.Event) error
	Subscribe(h func(e event.Event)) (dispose func(), err error)
	Events() (events []event.Event, err error)
}

type EventBus interface {
	Publish(event event.Event) error
}

type Commands struct {
	projects *project.Repo
	store    EventStore
	bus      EventBus
	dispose  func()
}

func New(store EventStore, bus EventBus,
	projects *project.Repo) *Commands {
	return &Commands{
		store:    store,
		bus:      bus,
		projects: projects,
	}
}

func (s *Commands) Start() error {
	events, err := s.store.Events()
	if err != nil {
		return err
	}

	s.processEvents(events...)

	dispose, err := s.store.Subscribe(func(e event.Event) {
		s.processEvents(e)
	})
	if err != nil {
		return err
	}

	s.dispose = dispose

	return nil
}

func (s *Commands) processEvents(events ...event.Event) {
	for _, e := range events {
		switch e.AggregateType() {
		case event.Project:
			s.projects.ProcessEvents(e)
		}
	}
}

func (s *Commands) publish(e event.Event) error {
	err := s.store.Publish(e)
	if err != nil {
		return err
	}

	return s.bus.Publish(e)
}
