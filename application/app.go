package application

import (
	"github.com/rdnt/rdnt.dev/event"
	"github.com/rdnt/rdnt.dev/project"
)

type App struct {
	projects *project.Repo
}

type todoEvents struct{}

func (t todoEvents) Events() ([]event.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (t todoEvents) Subscribe(h func(e event.Event)) (dispose func(), err error) {
	//TODO implement me
	panic("implement me")
}

func New() *App {
	todo := &todoEvents{}
	return &App{
		projects: project.NewRepo(todo),
	}
}
