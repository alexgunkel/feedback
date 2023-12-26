package handler

import (
	"bytes"
	"github.com/alexgunkel/feedback/config"
	"github.com/alexgunkel/feedback/storage"
	"github.com/alexgunkel/feedback/util"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"strconv"
)

type Feedback struct {
	s    *storage.Storage
	cfg  *config.Config
	html *template.Template
}

func NewFeedback(cfg *config.Config, s *storage.Storage) *Feedback {
	t, err := template.ParseFiles("./templates/feedback.html")
	util.PanicOnError(err)
	return &Feedback{s: s, cfg: cfg, html: t}
}

func (f *Feedback) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	tmpl, err := template.ParseFiles("./templates/feedback.html")
	util.PanicOnError(err)
	v := mux.Vars(request)
	id, _ := strconv.Atoi(v["id"])
	survey := f.s.GetSurvey(id)
	key := v["key"]
	if key != survey.AccessKey {
		writer.WriteHeader(403)
		return
	}
	if survey == nil {
		writer.WriteHeader(404)
		return
	}
	buf := bytes.NewBufferString("")
	err = tmpl.Execute(buf, struct {
		Cfg    *config.Config
		Survey *storage.Survey
	}{
		Cfg:    f.cfg,
		Survey: survey,
	})
	util.PanicOnError(err)
	_, _ = writer.Write(buf.Bytes())
}

var _ http.Handler = &Feedback{}
