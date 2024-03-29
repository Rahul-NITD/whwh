package systems_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/aargeee/whwh/drivers"
	"github.com/aargeee/whwh/handlers"
	"github.com/aargeee/whwh/specs"
	"github.com/alecthomas/assert/v2"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TesterServerStart implements specs.Tester.
func TesterServerStart() (serverUrl string, shutdown func(), err error) {
	handler := handlers.NewTesterServerHandler()
	svr := httptest.NewServer(handler)
	return svr.URL, svr.Close, nil
}

func TestSystem(t *testing.T) {

	svr, shutdown, err := TesterServerStart()
	assert.NoError(t, err, "Could not start TesterServer")
	t.Cleanup(shutdown)

	specs.TesterSpecification(t, drivers.NewSysDriver(svr))
}

func TestSystemDocker(t *testing.T) {

	if testing.Short() {
		t.Skip()
	}

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:       "../.",
			Dockerfile:    "./Dockerfile",
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

	specs.TesterSpecification(t, drivers.NewDockerDriver("http://localhost:8000"))
}
