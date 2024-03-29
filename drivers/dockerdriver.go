package drivers

import (
	"bytes"
	"net/http"

	sse "github.com/r3labs/sse/v2"
)

type DockerDriver struct {
	ServerUrl  string
	baseDriver *SysDriver
}

func NewDockerDriver(url string) *DockerDriver {
	return &DockerDriver{
		ServerUrl:  url,
		baseDriver: NewSysDriver(url),
	}
}

func (d *DockerDriver) GetServerUrl() string {
	return d.ServerUrl
}

// TesterServerStart implements specs.Tester.
func (d *DockerDriver) TesterServerStart() (serverUrl string, shutdown func(), err error) {
	return d.ServerUrl, func() {}, nil
}

// HookServerStart implements specs.Tester.
func (d *DockerDriver) HookServerStart(outputBuffer *bytes.Buffer) (hookUrl string, shutdown func(), err error) {
	return d.baseDriver.HookServerStart(outputBuffer)
}

// HealthCheck implements specs.Tester.
func (d *DockerDriver) HealthCheck(serverUrl string, hookUrl string) error {
	return d.baseDriver.HealthCheck(serverUrl, hookUrl)
}

// ClientConnect implements specs.Tester.
func (d *DockerDriver) ClientConnect(serverUrl string, hookUrl string) (client *sse.Client, sid string, err error) {
	return d.baseDriver.ClientConnect(serverUrl, hookUrl)
}

// ClientSubscribe implements specs.Tester.
func (d *DockerDriver) ClientSubscribe(client *sse.Client, sid string, hookUrl string) (unsubscribe func(), err error) {
	return d.baseDriver.ClientSubscribe(client, sid, hookUrl)
}

// MakeRequest implements specs.Tester.
func (d *DockerDriver) MakeRequest(req *http.Request) (res *http.Response, err error) {
	return d.baseDriver.MakeRequest(req)
}
