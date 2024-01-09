package commands

import (
	"errors"

	"github.com/rdnt/rdnt.dev/event"
	"github.com/rdnt/rdnt.dev/project"
	"github.com/rdnt/rdnt.dev/uuid"
)

func (s *Commands) CreateProject(id uuid.UUID, name string) error {
	_, err := s.projects.ProjectByName(name)
	if err == nil {
		return errors.New("project already exists")
	} else if !errors.Is(err, project.ErrNotFound) && err != nil {
		return err
	}

	e := event.ProjectCreatedEvent{
		ProjectId: id,
		Name:      name,
	}

	err = s.publish(e)
	if err != nil {
		return err
	}

	return nil
}
