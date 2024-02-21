package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Rahul-NITD/whwh/systems"
	"github.com/Rahul-NITD/whwh/systems/server"
)

type TesterServerHandler struct {
	http.Handler
	testerServer *server.TesterServer
	deferhome    func()
}

func NewTesterServerHandler(deferhome ...func()) *TesterServerHandler {

	t := &TesterServerHandler{}
	if len(deferhome) > 0 {
		t.deferhome = deferhome[0]
	}

	t.testerServer = server.NewTesterServer()

	r := http.NewServeMux()
	r.HandleFunc(systems.HOMEPATH, t.homeHandler)
	r.HandleFunc(systems.HEALTHPATH, t.healthHandler)
	r.HandleFunc(systems.CREATESTREAMPATH, t.createStreamHandler)
	r.HandleFunc(systems.EVENTSPATH, t.testerServer.EventServe())

	t.Handler = r

	return t
}

func (t *TesterServerHandler) homeHandler(w http.ResponseWriter, r *http.Request) {

	sid := r.URL.Query().Get("stream")

	val, err := marshalRequest(r)
	if err != nil {
		println("Error in Marshalling,", err)
	}

	t.testerServer.Publish(sid, val)

	defer t.deferhome()
}

func marshalRequest(r *http.Request) ([]byte, error) {
	var buf bytes.Buffer
	r.Write(&buf)
	return json.Marshal(buf.Bytes())
}

func (t *TesterServerHandler) createStreamHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, t.testerServer.CreateStream())
}

func (t *TesterServerHandler) healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "All systems healthy")
}
