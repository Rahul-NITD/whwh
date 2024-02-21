package handlers

import (
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

	err := t.testerServer.PublishRequest(sid, r)
	if err != nil {
		println("Error in Publishing,", err)
	}

	defer t.deferhome()
}

func (t *TesterServerHandler) createStreamHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, t.testerServer.CreateStream())
}

func (t *TesterServerHandler) healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "All systems healthy")
}
