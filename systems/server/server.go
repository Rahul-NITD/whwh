package server

import (
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

func (t *TesterServer) Publish(sid string, data []byte) {
	t.sseServer.Publish(sid, &sse.Event{
		Data: data,
	})
}
