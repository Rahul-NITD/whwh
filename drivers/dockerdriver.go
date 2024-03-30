package drivers

import (
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

// HealthCheck implements specs.Tester.
func (d *DockerDriver) HealthCheck() error {
	return d.baseDriver.HealthCheck()
}

// ClientConnect implements specs.Tester.
func (d *DockerDriver) ClientConnect(hookUrl string) (client *sse.Client, sid string, err error) {
	return d.baseDriver.ClientConnect(hookUrl)
}

// ClientSubscribe implements specs.Tester.
func (d *DockerDriver) ClientSubscribe(client *sse.Client, sid string, hookUrl string) (unsubscribe func(), err error) {
	return d.baseDriver.ClientSubscribe(client, sid, hookUrl)
}

// MakeRequest implements specs.Tester.
func (d *DockerDriver) MakeRequest(req *http.Request) (res *http.Response, err error) {
	return d.baseDriver.MakeRequest(req)
}
