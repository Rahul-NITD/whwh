package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/aargeee/whwh/systems"
	"github.com/aargeee/whwh/systems/server"
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
	sid := t.testerServer.CreateStream()

	response := systems.StreamPayloadResponse{
		Event: "CREATE_STREAM",
		Payload: systems.StreamPayload{
			StreamID: sid,
		},
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (t *TesterServerHandler) healthHandler(w http.ResponseWriter, r *http.Request) {
	err := t.testerServer.ReportHealth()
	report := systems.HealthReport{
		Event: "HEALTHZ",
	}
	if err != nil {
		report.Status = "UNHEALTHY"
		report.Message = err.Error()
	} else {
		report.Status = "HEALTHY"
		report.Message = "All Systems are Healthy"
	}
	if err := json.NewEncoder(w).Encode(report); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
