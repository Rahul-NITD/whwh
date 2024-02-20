package drivers

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"

	sse "github.com/r3labs/sse/v2"
)

type SysDriver struct {
	serverUrl string
	hookUrl   string
	mu        sync.Mutex
	done      chan struct{}
}

func NewSysDriver() *SysDriver {
	return &SysDriver{
		done: make(chan struct{}),
	}
}

// TesterServerStart implements specs.Tester.
func (d *SysDriver) TesterServerStart() (serverUrl string, shutdown func(), err error) {
	r := http.NewServeMux()

	sseserver := sse.New()

	testsid := "ccb04aa8-1b6f-4f24-b9de-d274b39c5128"

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		r.Write(&buf)

		val, err := json.Marshal(buf.Bytes())
		if err != nil {
			println("Error in Marshalling,", err)
		}

		sseserver.Publish(testsid, &sse.Event{
			Data: val,
		})
		<-d.done
	})
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "All systems healthy")
	})
	r.HandleFunc("/createstream", func(w http.ResponseWriter, r *http.Request) {
		sid := testsid
		sseserver.CreateStream(sid)
		fmt.Fprint(w, sid)
	})
	r.HandleFunc("/events", sseserver.ServeHTTP)

	svr := httptest.NewServer(r)
	return svr.URL, svr.Close, nil
}

// HookServerStart implements specs.Tester.
func (d *SysDriver) HookServerStart(outputBuffer *bytes.Buffer) (hookUrl string, shutdown func(), err error) {

	r := http.NewServeMux()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		d.mu.Lock()
		defer d.mu.Unlock()
		outputBuffer.Reset()
		r.Write(outputBuffer)
	})

	svr := httptest.NewServer(r)
	return svr.URL, svr.Close, nil
}

// HealthCheck implements specs.Tester.
func (*SysDriver) HealthCheck(serverUrl string, hookUrl string) error {
	if err := makeGetRequest(serverUrl); err != nil {
		return err
	}
	return nil
}

func makeGetRequest(url string) error {
	res, err := http.Get(url + "/health")
	if err != nil {
		return fmt.Errorf("could not make request to server, %s, %s", url+"/health", err.Error())
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server not healthy, got %d", res.StatusCode)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if string(body) != "All systems healthy" {
		return fmt.Errorf(`did not receive "All systems healthy"`)
	}
	return nil
}

// ClientConnect implements specs.Tester.
func (s *SysDriver) ClientConnect(serverUrl string, hookUrl string) (*sse.Client, string, error) {
	s.serverUrl = serverUrl
	s.hookUrl = hookUrl

	res, err := http.Get(serverUrl + "/createstream")
	if err != nil {
		return nil, "", err
	}
	sid, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, "", err
	}
	cli := sse.NewClient(serverUrl + "/events")
	return cli, string(sid), nil
}

// ClientSubscribe implements specs.Tester.
func (s *SysDriver) ClientSubscribe(client *sse.Client, sid string, hookUrl string) (unsubscribe func(), err error) {

	cxt, cancel := context.WithCancel(context.Background())
	go client.SubscribeWithContext(cxt, sid, func(msg *sse.Event) {

		var dec []byte
		json.Unmarshal(msg.Data, &dec)
		req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(dec)))
		if err != nil {
			println("Err in reading request, ", err.Error())
			close(s.done)
			return
		}
		req.URL, err = url.Parse(hookUrl)
		req.Host = req.URL.Host
		if err != nil {
			println("Err in reading request, ", err.Error())
			close(s.done)
			return
		}

		req.RequestURI = ""

		_, err = http.DefaultClient.Do(req)
		if err != nil {
			println("Err in reading request, ", err.Error())
			close(s.done)
			return
		}
		close(s.done)
	})
	return cancel, nil
}

// MakeRequest implements specs.Tester.
func (*SysDriver) MakeRequest(req *http.Request) (res *http.Response, err error) {
	return http.DefaultClient.Do(req)
}
