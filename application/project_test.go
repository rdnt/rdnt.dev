package application_test

import (
	"testing"

	"gotest.tools/assert"

	"github.com/rdnt/rdnt.dev/uuid"
)

func TestServer(t *testing.T) {
	s := newSuite(t)

	pid := uuid.New()
	t.Run("create project", func(t *testing.T) {
		name := "first project"
		err := s.commands.CreateProject(pid, name)
		assert.NilError(t, err)

		p, err := s.projectRepo.Project(pid)
		assert.NilError(t, err)
		assert.Equal(t, p.Id, pid)
		assert.Equal(t, p.Name, name)

		p, err = s.projectRepo.ProjectByName(name)
		assert.NilError(t, err)
		assert.Equal(t, p.Id, pid)
		assert.Equal(t, p.Name, name)
	})
}
