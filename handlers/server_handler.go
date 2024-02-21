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
}

func NewTesterServerHandler(deferfunc ...func()) *TesterServerHandler {

	t := &TesterServerHandler{}

	t.testerServer = server.NewTesterServer()

	r := http.NewServeMux()
	r.HandleFunc(systems.HOMEPATH, t.homeHandler(deferfunc...))
	r.HandleFunc(systems.HEALTHPATH, t.healthHandler)
	r.HandleFunc(systems.CREATESTREAMPATH, t.createStreamHandler)
	r.HandleFunc(systems.EVENTSPATH, t.testerServer.EventServe())

	t.Handler = r

	return t
}

func (t *TesterServerHandler) homeHandler(deferfunc ...func()) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		sid := r.URL.Query().Get("stream")

		err := t.testerServer.PublishRequest(sid, r)
		if err != nil {
			println("Error in Publishing,", err)
		}
		defer dodefers(deferfunc...)
	}
}

func dodefers(deferfunc ...func()) {
	for _, f := range deferfunc {
		f()
	}
}

func (t *TesterServerHandler) createStreamHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, t.testerServer.CreateStream())
}

func (t *TesterServerHandler) healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "All systems healthy")
}
