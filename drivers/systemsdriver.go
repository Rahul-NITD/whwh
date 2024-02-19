package drivers

import (
	"io"
	"net/http"

	"github.com/Rahul-NITD/whwh/specs"
)

type SysDriver struct{}

// TesterServerStart implements specs.Tester.
func (*SysDriver) TesterServerStart() (serverUrl string, shutdown func(), err error) {
	panic("unimplemented")
}

// HookServerStart implements specs.Tester.
func (*SysDriver) HookServerStart() (hookUrl string, outputBuffer io.Reader, shutdown func(), err error) {
	panic("unimplemented")
}

// HealthCheck implements specs.Tester.
func (*SysDriver) HealthCheck(serverUrl string, hookUrl string) error {
	panic("unimplemented")
}

// ClientConnect implements specs.Tester.
func (*SysDriver) ClientConnect(serverUrl string, hookUrl string) (client specs.ConnectedClient, sid string, err error) {
	panic("unimplemented")
}

// ClientSubscribe implements specs.Tester.
func (*SysDriver) ClientSubscribe(client specs.ConnectedClient, eventUrl string, sid string) (unsubscribe func(), err error) {
	panic("unimplemented")
}

// MakeRequest implements specs.Tester.
func (*SysDriver) MakeRequest(req *http.Request) (res *http.Response, err error) {
	panic("unimplemented")
}
