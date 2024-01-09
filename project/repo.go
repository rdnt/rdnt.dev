package project

import (
	"errors"
	"sync"

	"github.com/rdnt/rdnt.dev/event"
	"github.com/rdnt/rdnt.dev/uuid"
)

type EventStore interface {
	Events() ([]event.Event, error)
	Subscribe(h func(e event.Event)) (dispose func(), err error)
}

type Repo struct {
	mux      sync.Mutex
	projects map[uuid.UUID]*Aggregate
	events   EventStore
}

func NewRepo(events EventStore) *Repo {
	r := &Repo{
		events:   events,
		projects: map[uuid.UUID]*Aggregate{},
	}

	return r
}

func (r *Repo) ProcessEvents(events ...event.Event) {
	r.mux.Lock()

	for _, e := range events {
		_, ok := r.projects[e.AggregateId()]
		if !ok {
			r.projects[e.AggregateId()] = &Aggregate{}
		}

		r.projects[e.AggregateId()].ProcessEvent(e)
	}

	r.mux.Unlock()
}

var ErrNotFound = errors.New("project not found")

func (r *Repo) Project(id uuid.UUID) (Project, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	p, ok := r.projects[id]
	if !ok {
		return Project{}, ErrNotFound
	}

	return p.Project, nil
}

func (r *Repo) ProjectByName(name string) (Project, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	for _, p := range r.projects {
		if p.Name == name {
			return p.Project, nil
		}
	}

	return Project{}, ErrNotFound
}
