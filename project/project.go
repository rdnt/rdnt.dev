package project

import "github.com/rdnt/rdnt.dev/uuid"

type Project struct {
	Id   uuid.UUID
	Name string
}

func New(id uuid.UUID, name string) Project {
	return Project{
		Id:   id,
		Name: name,
	}
}
