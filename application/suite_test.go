package application_test

import (
	"testing"
	"time"

	"gotest.tools/assert"

	"github.com/rdnt/rdnt.dev/application/commands"
	query "github.com/rdnt/rdnt.dev/application/queries"
	"github.com/rdnt/rdnt.dev/event"
	"github.com/rdnt/rdnt.dev/memstore"
	"github.com/rdnt/rdnt.dev/project"
)

type EventStore interface {
	Publish(event event.Event) error
	Subscribe(h func(e event.Event)) (dispose func(), err error)
	Events() (events []event.Event, err error)
}

type suite struct {
	bus         EventStore
	store       EventStore
	commands    *commands.Commands
	projectRepo *project.Repo
	queries     *query.Queries
	projectView *project.Repo
}

func newSuite(t *testing.T) *suite {
	eventStore := memstore.New()
	eventBus := eventStore

	projectRepo := project.NewRepo(eventStore)

	commandHandler := commands.New(
		eventStore,
		eventBus,
		projectRepo,
	)

	err := commandHandler.Start()
	assert.NilError(t, err)

	projectView := project.NewRepo(eventBus)

	queries := query.New(
		eventBus,
		projectView,
	)

	err = queries.Start()
	assert.NilError(t, err)

	return &suite{
		bus:         eventBus,
		store:       eventStore,
		projectRepo: projectRepo,
		commands:    commandHandler,
		projectView: projectView,
		queries:     queries,
	}
}

func eventually(t *testing.T, f func() bool) {
	t.Helper()

	ch := make(chan bool, 1)

	timer := time.NewTimer(1 * time.Second)
	defer timer.Stop()

	ticker := time.NewTicker(1 * time.Nanosecond)
	defer ticker.Stop()

	for tick := ticker.C; ; {
		select {
		case <-timer.C:
			t.Fail()
			return
		case <-tick:
			tick = nil
			go func() { ch <- f() }()
		case v := <-ch:
			if v {
				return
			}
			tick = ticker.C
		}
	}
}
