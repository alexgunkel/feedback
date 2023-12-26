package handler

import (
	"github.com/alexgunkel/feedback/config"
	"github.com/alexgunkel/feedback/storage"
	"github.com/alexgunkel/feedback/util"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"strconv"
)

type Result struct {
	storage *storage.Storage
	cfg     *config.Config
}

func (s *Result) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	tmpl, err := template.ParseFiles("./templates/result.html")
	util.PanicOnError(err)
	v := mux.Vars(request)
	id, _ := strconv.Atoi(v["id"])
	survey := s.storage.GetSurvey(id)
	key := v["key"]
	if key != survey.EvaluationKey {
		writer.WriteHeader(403)
		return
	}

	results := s.storage.GetResult(*survey)

	err = tmpl.Execute(writer, struct {
		Res storage.Result
		Cfg *config.Config
	}{
		Res: results,
		Cfg: s.cfg,
	})
	util.PanicOnError(err)
}

var _ http.Handler = &Result{}

func NewStatistics(st *storage.Storage, cfg *config.Config) *Result {
	return &Result{storage: st, cfg: cfg}
}
