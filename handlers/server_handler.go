package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Rahul-NITD/whwh/systems"
	"github.com/google/uuid"
	"github.com/r3labs/sse/v2"
)

type TesterServerHandler struct {
	http.Handler
	sseServer *sse.Server
	deferhome func()
}

func NewTesterServerHandler(deferhome ...func()) *TesterServerHandler {

	t := &TesterServerHandler{}
	if len(deferhome) > 0 {
		t.deferhome = deferhome[0]
	}

	r := http.NewServeMux()

	t.sseServer = sse.New()
	t.sseServer.AutoReplay = false

	r.HandleFunc("/", t.homeHandler)
	r.HandleFunc("/health", t.healthHandler)
	r.HandleFunc(systems.CREATESTREAMPATH, t.createStreamHandler)
	r.HandleFunc("/events", t.sseServer.ServeHTTP)

	t.Handler = r

	return t
}

func (t *TesterServerHandler) homeHandler(w http.ResponseWriter, r *http.Request) {

	sid := r.URL.Query().Get("stream")

	var buf bytes.Buffer
	r.Write(&buf)

	val, err := json.Marshal(buf.Bytes())
	if err != nil {
		println("Error in Marshalling,", err)
	}

	t.sseServer.Publish(sid, &sse.Event{
		Data: val,
	})

	defer t.deferhome()
}

func (t *TesterServerHandler) createStreamHandler(w http.ResponseWriter, r *http.Request) {
	sid := uuid.NewString()
	t.sseServer.CreateStream(sid)
	fmt.Fprint(w, sid)
}

func (t *TesterServerHandler) healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "All systems healthy")
}
