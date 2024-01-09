package query

import (
	"github.com/rdnt/rdnt.dev/event"
	"github.com/rdnt/rdnt.dev/project"
	"github.com/rdnt/rdnt.dev/uuid"
)

type EventBus interface {
	Subscribe(handler func(event event.Event)) (func(), error)
	Events() (events []event.Event, err error)
}

type ProjectView interface {
	ProcessEvents(events ...event.Event)
	Project(id uuid.UUID) (project.Project, error)
}

type Queries struct {
	store    EventBus
	projects ProjectView
	dispose  func()
}

func (q *Queries) Project(id uuid.UUID) (project.Project, error) {
	return q.projects.Project(id)
}

func New(events EventBus, projects ProjectView) *Queries {
	s := &Queries{
		store:    events,
		projects: projects,
	}

	//go func() {
	//	for {
	//		func() {
	//			store, err := s.store.Subscribe()
	//			if err != nil {
	//				log.Error(err)
	//				return
	//			}
	//
	//			for e := range store {
	//				log.Debug("[view] receive ", e)
	//				err = s.handleEvent(e)
	//				if err != nil {
	//					log.Error(err)
	//					continue
	//				}
	//			}
	//		}()
	//	}
	//}()

	return s
}

func (q *Queries) Start() error {
	events, err := q.store.Events()
	if err != nil {
		return err
	}

	q.processEvents(events...)

	dispose, err := q.store.Subscribe(func(e event.Event) {
		q.processEvents(e)
	})
	if err != nil {
		return err
	}

	q.dispose = dispose

	return nil
}

func (q *Queries) processEvents(events ...event.Event) {
	for _, e := range events {
		switch e.AggregateType() {
		case event.Project:
			q.projects.ProcessEvents(e)
		}
	}
}
