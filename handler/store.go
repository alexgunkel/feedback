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

type Store struct {
	c chan<- storage.Input

	html []byte
}

func NewStore(cfg *config.Config, c chan<- storage.Input) *Store {
	t, err := template.ParseFiles("./templates/thank-you.html")
	util.PanicOnError(err)
	buf := bytes.NewBufferString("")
	err = t.Execute(buf, &cfg)
	util.PanicOnError(err)
	return &Store{c: c, html: buf.Bytes()}
}

func (s *Store) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	survey, _ := strconv.Atoi(vars["id"])
	err := request.ParseForm()
	if err != nil {
		writer.WriteHeader(403)
		return
	}

	s.c <- storage.Input{
		Val:    request.Form,
		Survey: uint(survey),
	}

	writer.WriteHeader(201)
	_, err = writer.Write(s.html)
	util.PanicOnError(err)
}

var _ http.Handler = &Store{}
