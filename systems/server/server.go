package server

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/r3labs/sse/v2"
)

type TesterServer struct {
	sseServer *sse.Server
}

func NewTesterServer() *TesterServer {
	t := &TesterServer{
		sseServer: sse.New(),
	}
	t.sseServer.AutoReplay = false
	return t
}

func (t *TesterServer) CreateStream() string {
	sid := uuid.NewString()
	t.sseServer.CreateStream(sid)
	return sid
}

func (t *TesterServer) EventServe() func(http.ResponseWriter, *http.Request) {
	return t.sseServer.ServeHTTP
}

func (t *TesterServer) PublishRequest(sid string, r *http.Request) error {

	data, err := marshalRequest(r)
	if err != nil {
		return err
	}
	t.sseServer.Publish(sid, &sse.Event{
		Data: data,
	})
	return nil
}

func marshalRequest(r *http.Request) ([]byte, error) {
	var buf bytes.Buffer
	r.Write(&buf)
	return json.Marshal(buf.Bytes())
}
