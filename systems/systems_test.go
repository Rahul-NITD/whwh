package systems_test

import (
	"context"
	"testing"

	"github.com/Rahul-NITD/whwh/drivers"
	"github.com/Rahul-NITD/whwh/specs"
	"github.com/alecthomas/assert/v2"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestSystem(t *testing.T) {
	specs.TesterSpecification(t, drivers.NewSysDriver())
}

func TestSystemDocker(t *testing.T) {

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:       "../.",
			Dockerfile:    "./drivers/Dockerfile",
			PrintBuildLog: true,
		},
		ExposedPorts: []string{"8000:8000"},
		WaitingFor:   wait.ForHTTP("/health").WithPort("8000"),
	}

	cxt := context.Background()

	container, err := testcontainers.GenericContainer(cxt, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	assert.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, container.Terminate(cxt))
	})

	specs.TesterSpecification(t, drivers.DocDriver{})
}
