package event

import "github.com/rdnt/rdnt.dev/uuid"

const (
	ProjectCreated Type = "project-created"
)

type ProjectCreatedEvent struct {
	ProjectId uuid.UUID
	Name      string
}

func (ProjectCreatedEvent) Type() Type {
	return ProjectCreated
}

func (ProjectCreatedEvent) AggregateType() AggregateType {
	return Project
}

func (e ProjectCreatedEvent) AggregateId() uuid.UUID {
	return e.ProjectId
}
