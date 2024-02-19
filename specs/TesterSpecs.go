package specs

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/r3labs/sse/v2"
)

type ConnectedClient interface {
	SubscribeChanWithContext(cxt context.Context, stream string, eventChannel chan *sse.Event)
}

type Tester interface {
	TesterServerStart() (serverUrl string, shutdown func(), err error)
	HookServerStart() (hookUrl string, outputBuffer io.Reader, shutdown func(), err error)

	HealthCheck(serverUrl, hookUrl string) error

	ClientConnect(eventUrl, hookUrl string) (client ConnectedClient, sid string, err error)
	ClientSubscribe(client ConnectedClient, sid string) (unsubscribe func(), err error)
	MakeRequest(req *http.Request) (res *http.Response, err error)
}

func TesterSpecification(t *testing.T, tester Tester) {
	t.Helper()
	// For server
	serverUrl, shutdown, err := tester.TesterServerStart()
	assert.NoError(t, err, "Tester Server could not be started")
	t.Cleanup(shutdown)

	// For hook
	hookUrl, outputBuffer, shutdown, err := tester.HookServerStart()
	assert.NoError(t, err, "Tester Server could not be started")
	t.Cleanup(shutdown)

	// HealthCheck
	assert.NoError(t, tester.HealthCheck(serverUrl, hookUrl), "Server or Hook not healthy")

	// Client Side
	eventUrl := serverUrl + "/events"
	client, sid, err := tester.ClientConnect(eventUrl, hookUrl)
	assert.NoError(t, err, "Client could not establish connection")

	// Subscribe to event
	ubsubscribe, err := tester.ClientSubscribe(client, sid)
	assert.NoError(t, err, "Could not subscribe to event")
	t.Cleanup(ubsubscribe)

	// Test Request
	t.Run("Test POST without body", func(t *testing.T) {
		afterHook := makeRequestGetHookOutput(t, tester, hookUrl, outputBuffer)
		afterServer := makeRequestGetHookOutput(t, tester, serverUrl, outputBuffer)
		assert.Equal[[]byte](t, afterHook, afterServer)
	})
}

func makeRequestGetHookOutput(t *testing.T, tester Tester, url string, outputBuffer io.Reader) []byte {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, url, http.NoBody)
	assert.NoError(t, err, "Error while creating POST request "+url)
	res, err := tester.MakeRequest(req)
	assert.NoError(t, err, "Error in making request to hook "+url)
	assert.Equal[int](t, http.StatusOK, res.StatusCode)

	// save output buffer response
	afterHook, err := io.ReadAll(outputBuffer)
	assert.NoError(t, err)
	assert.NotEqual[[]byte](t, []byte{}, afterHook)
	return afterHook
}
