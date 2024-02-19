package drivers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Rahul-NITD/whwh/specs"
)

type SysDriver struct{}

// TesterServerStart implements specs.Tester.
func (*SysDriver) TesterServerStart() (serverUrl string, shutdown func(), err error) {
	return "", func() {}, fmt.Errorf("yet to implement")
}

// HookServerStart implements specs.Tester.
func (*SysDriver) HookServerStart() (hookUrl string, outputBuffer io.Reader, shutdown func(), err error) {
	return "", nil, func() {}, fmt.Errorf("yet to implement")
}

// HealthCheck implements specs.Tester.
func (*SysDriver) HealthCheck(serverUrl string, hookUrl string) error {
	return fmt.Errorf("yet to implement")
}

// ClientConnect implements specs.Tester.
func (*SysDriver) ClientConnect(serverUrl string, hookUrl string) (client specs.ConnectedClient, sid string, err error) {
	return nil, "", fmt.Errorf("yet to implement")
}

// ClientSubscribe implements specs.Tester.
func (*SysDriver) ClientSubscribe(client specs.ConnectedClient, eventUrl string, sid string) (unsubscribe func(), err error) {
	return func() {}, fmt.Errorf("yet to implement")
}

// MakeRequest implements specs.Tester.
func (*SysDriver) MakeRequest(req *http.Request) (res *http.Response, err error) {
	return nil, fmt.Errorf("yet to implement")
}
