package drivers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/aargeee/whwh/systems"
	"github.com/aargeee/whwh/systems/client"
	"github.com/aargeee/whwh/systems/hook"
	sse "github.com/r3labs/sse/v2"
)

type SysDriver struct {
	done      chan struct{}
	serverUrl string
	isServer  bool
}

func NewSysDriver(svr string) *SysDriver {
	return &SysDriver{
		done:      make(chan struct{}),
		isServer:  false,
		serverUrl: svr,
	}
}

func (d *SysDriver) GetServerUrl() string {
	return d.serverUrl
}

// HookServerStart implements specs.Tester.
func (d *SysDriver) HookServerStart(outputBuffer *bytes.Buffer) (hookUrl string, shutdown func(), err error) {
	svr := httptest.NewServer(hook.NewHook(outputBuffer))
	return svr.URL, svr.Close, nil
}

// HealthCheck implements specs.Tester.
func (d *SysDriver) HealthCheck() error {
	if err := makeHealthRequest(d.serverUrl); err != nil {
		return err
	}
	return nil
}

func makeHealthRequest(url string) error {
	res, err := http.Get(url + systems.HEALTHPATH)
	if err != nil {
		return fmt.Errorf("could not make request to server, %s, %s", url+systems.HEALTHPATH, err.Error())
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server not healthy, got %d", res.StatusCode)
	}

	defer res.Body.Close()

	var response systems.HealthReport
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return err
	}

	if response.Status != "HEALTHY" {
		return fmt.Errorf("systems returned %s", response.Status)
	}

	return nil
}

// ClientConnect implements specs.Tester.
func (s *SysDriver) ClientConnect(hookUrl string) (*sse.Client, string, error) {
	return client.ClientConnect(s.serverUrl, hookUrl)
}

// ClientSubscribe implements specs.Tester.
func (s *SysDriver) ClientSubscribe(clientD *sse.Client, sid string, hookUrl string) (unsubscribe func(), err error) {
	return client.ClientSubscribe(clientD, sid, hookUrl, func() { close(s.done) }) // To notify homehandler that request is complete
}

// MakeRequest implements specs.Tester.
func (d *SysDriver) MakeRequest(req *http.Request) (res *http.Response, err error) {
	res, err = http.DefaultClient.Do(req)
	if d.isServer {
		<-d.done
	}
	d.isServer = true
	return
}
