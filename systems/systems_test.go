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

	t.Run("Test Without Docker", func(t *testing.T) {
		specs.TesterSpecification(t, drivers.NewSysDriver(svr), specs.TestArguments{
			ServerUrl:    svr,
			HookUrl:      hookUrl,
			OutputBuffer: outputBuffer,
		})
	})

	t.Run("With Docker", DockerTest(hookUrl, outputBuffer))
}

func DockerTest(hookUrl string, outputBuffer *bytes.Buffer) func(t *testing.T) {
	return func(t *testing.T) {

		if testing.Short() {
			t.Skip()
		}

		cxt := context.Background()

		driver, term, err := drivers.NewDockerDriver("http://localhost:8000", cxt)
		assert.NoError(t, err)
		t.Cleanup(func() {
			assert.NoError(t, term(cxt), "Error terminating")
		})

		specs.TesterSpecification(t, driver, specs.TestArguments{
			ServerUrl:    "http://localhost:8000",
			HookUrl:      hookUrl,
			OutputBuffer: outputBuffer,
		})
	}
}
