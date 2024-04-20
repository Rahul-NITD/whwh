package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path"

	"github.com/aargeee/whwh/systems"
	"github.com/aargeee/whwh/systems/server"
)

type TesterServerHandler struct {
	http.Handler
	testerServer *server.TesterServer
	templ        *template.Template
}

func NewTesterServerHandler(deferfunc ...func()) *TesterServerHandler {

	t := &TesterServerHandler{}

	t.testerServer = server.NewTesterServer()

	r := http.NewServeMux()
	r.HandleFunc(systems.HOMEPATH, t.homeHandler(deferfunc...))
	r.HandleFunc(systems.HEALTHPATH, t.healthHandler)
	r.HandleFunc(systems.CREATESTREAMPATH, t.createStreamHandler)
	r.HandleFunc(systems.EVENTSPATH, t.testerServer.EventServe())
	r.HandleFunc(systems.HOW_TO_GUIDE, t.howToHandler)

	t.Handler = r
	templ, err := template.ParseFiles(path.Join("templ", "how_to_guide.html"))
	if err != nil {
		log.Fatal(err)
	}
	t.templ = templ

	return t
}

func (t *TesterServerHandler) howToHandler(w http.ResponseWriter, r *http.Request) {
	if err := t.templ.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (t *TesterServerHandler) homeHandler(deferfunc ...func()) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		sid := r.URL.Query().Get("stream")

		if sid == "" {
			http.Redirect(w, r, "http://whwh.rahulgoel.dev/how_to", http.StatusSeeOther)
		}

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
