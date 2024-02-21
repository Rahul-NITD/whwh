package drivers

import (
	"bytes"
	"net/http"

	sse "github.com/r3labs/sse/v2"
)

type DocDriver struct{}

// TesterServerStart implements specs.Tester.
func (DocDriver) TesterServerStart() (serverUrl string, shutdown func(), err error) {
	panic("unimplemented")
}

// HookServerStart implements specs.Tester.
func (DocDriver) HookServerStart(outputBuffer *bytes.Buffer) (hookUrl string, shutdown func(), err error) {
	panic("unimplemented")
}

// HealthCheck implements specs.Tester.
func (DocDriver) HealthCheck(serverUrl string, hookUrl string) error {
	panic("unimplemented")
}

// ClientConnect implements specs.Tester.
func (DocDriver) ClientConnect(eventUrl string, hookUrl string) (client *sse.Client, sid string, err error) {
	panic("unimplemented")
}

// ClientSubscribe implements specs.Tester.
func (DocDriver) ClientSubscribe(client *sse.Client, sid string, hookUrl string) (unsubscribe func(), err error) {
	panic("unimplemented")
}

// MakeRequest implements specs.Tester.
func (DocDriver) MakeRequest(req *http.Request) (res *http.Response, err error) {
	panic("unimplemented")
}
