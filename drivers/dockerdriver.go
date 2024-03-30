package drivers

import (
	"context"
	"net/http"

	sse "github.com/r3labs/sse/v2"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type DockerDriver struct {
	ServerUrl  string
	baseDriver *SysDriver
}

func NewDockerDriver(url string, cxt context.Context) (*DockerDriver, func(context.Context) error, error) {

	container, err := testcontainers.GenericContainer(cxt, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			FromDockerfile: testcontainers.FromDockerfile{
				Context:       "../.",
				Dockerfile:    "./Dockerfile",
				PrintBuildLog: true,
			},
			ExposedPorts: []string{"8000:8000"},
			WaitingFor:   wait.ForHTTP("/health").WithPort("8000"),
		},
		Started: true,
	})

	return &DockerDriver{
		ServerUrl:  url,
		baseDriver: NewSysDriver(url),
	}, container.Terminate, err
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
