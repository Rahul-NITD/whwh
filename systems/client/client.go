package client

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/r3labs/sse/v2"
)

type TesterClient struct {
	deferfunc func()
}

func ClientConnect(serverUrl string, hookUrl string) (*sse.Client, string, error) {

	res, err := http.Get(serverUrl + "/createstream")
	if err != nil {
		return nil, "", err
	}
	sid, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, "", err
	}
	cli := sse.NewClient(serverUrl + "/events")

	return cli, string(sid), nil
}

func ClientSubscribe(client *sse.Client, sid string, hookUrl string, deferfunc ...func()) (unsubscribe func(), err error) {

	cxt, cancel := context.WithCancel(context.Background())
	go client.SubscribeWithContext(cxt, sid, func(msg *sse.Event) {
		defer func() {
			for _, f := range deferfunc {
				f()
			}
		}()
		var dec []byte
		json.Unmarshal(msg.Data, &dec)
		req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(dec)))
		if err != nil {
			println("Err in reading request, ", err.Error())
			// close(s.done)
			return
		}
		req.URL, err = url.Parse(hookUrl)
		req.Host = req.URL.Host
		if err != nil {
			println("Err in reading request, ", err.Error())
			// close(s.done)
			return
		}

		req.RequestURI = ""

		_, err = http.DefaultClient.Do(req)
		if err != nil {
			println("Err in reading request, ", err.Error())
			// close(s.done)
			return
		}
		// close(s.done)
	})
	return cancel, nil
}
