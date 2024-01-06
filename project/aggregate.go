package project

import (
	"fmt"

	"github.com/rdnt/rdnt.dev/event"
)

type Aggregate struct {
	Project
}

func (p *Project) ProcessEvent(e event.Event) {
	switch e := e.(type) {
	case event.ProjectCreatedEvent:
		p.Id = e.ProjectId
		p.Name = e.Name
	default:
		fmt.Println("project: unknown event", e)
	}
}
