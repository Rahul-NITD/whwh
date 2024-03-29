package specs

import (
	"bufio"
	"bytes"
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/r3labs/sse/v2"
)

type ConnectedClient interface {
	SubscribeChanWithContext(cxt context.Context, stream string, eventChannel chan *sse.Event) error
}

type Tester interface {
	GetServerUrl() string
	HookServerStart(outputBuffer *bytes.Buffer) (hookUrl string, shutdown func(), err error)

	HealthCheck(serverUrl, hookUrl string) error

	ClientConnect(eventUrl, hookUrl string) (client *sse.Client, sid string, err error)
	ClientSubscribe(client *sse.Client, sid string, hookUrl string) (unsubscribe func(), err error)
	MakeRequest(req *http.Request) (res *http.Response, err error)
}

func TesterSpecification(t *testing.T, subject Tester) {
	t.Helper()
	// For server
	serverUrl := subject.GetServerUrl()

	// For hook
	outputBuffer := &bytes.Buffer{}
	hookUrl, shutdown, err := subject.HookServerStart(outputBuffer)
	assert.NoError(t, err, "Tester Server could not be started")
	t.Cleanup(shutdown)

	// HealthCheck
	assert.NoError(t, subject.HealthCheck(serverUrl, hookUrl), "Server not healthy")

	// Client Side
	client, sid, err := subject.ClientConnect(serverUrl, hookUrl)
	assert.NoError(t, err, "Client could not establish connection")

	// Subscribe to event
	ubsubscribe, err := subject.ClientSubscribe(client, sid, hookUrl)
	assert.NoError(t, err, "Could not subscribe to event")
	t.Cleanup(ubsubscribe)

	// Test Request
	afterHook := makeRequestGetHookOutput(t, subject, hookUrl, sid, outputBuffer)
	assert.NoError(t, assertLineToRequest(afterHook), "Could not form a request from the line representations returned, "+string(afterHook))

	afterServer := makeRequestGetHookOutput(t, subject, serverUrl, sid, outputBuffer)
	assert.NoError(t, assertLineToRequest(afterServer), "Could not form a request from the line representations returned, "+string(afterHook))
	assert.Equal[string](t, afterHook, afterServer)
}

func assertLineToRequest(p string) error {
	_, err := http.ReadRequest(bufio.NewReader(strings.NewReader(p)))
	return err
}

func makeRequestGetHookOutput(t *testing.T, subject Tester, url string, sid string, outputBuffer *bytes.Buffer) string {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, url, http.NoBody)
	assert.NoError(t, err, "Error while creating POST request "+url)
	q := req.URL.Query()
	q.Add("stream", sid)
	req.URL.RawQuery = q.Encode()
	res, err := subject.MakeRequest(req)
	assert.NoError(t, err, "Error in making request to hook "+url)
	assert.Equal[int](t, http.StatusOK, res.StatusCode)

	// save output buffer response
	afterHook := outputBuffer.Bytes()
	assert.NotEqual[[]byte](t, []byte{}, afterHook)
	defer outputBuffer.Reset()
	return string(afterHook)
}
