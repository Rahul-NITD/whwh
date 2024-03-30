package systems_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aargeee/whwh/drivers"
	"github.com/aargeee/whwh/handlers"
	"github.com/aargeee/whwh/specs"
	"github.com/aargeee/whwh/systems/hook"
	"github.com/alecthomas/assert/v2"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TesterServerStart implements specs.Tester.
func SpinServer(handler http.Handler) (serverUrl string, shutdown func(), err error) {
	svr := httptest.NewServer(handler)
	return svr.URL, svr.Close, nil
}

func TestSystem(t *testing.T) {
	svr, shutdown, err := SpinServer(handlers.NewTesterServerHandler())
	assert.NoError(t, err, "Could not start TesterServer")
	t.Cleanup(shutdown)

	outputBuffer := &bytes.Buffer{}
	hookUrl, shutdown, err := SpinServer(hook.NewHook(outputBuffer))
	assert.NoError(t, err, "Could not start hook service")
	t.Cleanup(shutdown)

	specs.TesterSpecification(t, drivers.NewSysDriver(svr), specs.TestArguments{
		ServerUrl:    svr,
		HookUrl:      hookUrl,
		OutputBuffer: outputBuffer,
	})
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

	outputBuffer := &bytes.Buffer{}
	hookUrl, shutdown, err := SpinServer(hook.NewHook(outputBuffer))
	assert.NoError(t, err, "Could not start hook service")
	t.Cleanup(shutdown)

	specs.TesterSpecification(t, drivers.NewDockerDriver("http://localhost:8000"), specs.TestArguments{
		ServerUrl:    "http://localhost:8000",
		HookUrl:      hookUrl,
		OutputBuffer: outputBuffer,
	})
}
