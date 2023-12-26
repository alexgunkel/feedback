package handler

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"github.com/alexgunkel/feedback/config"
	"github.com/alexgunkel/feedback/storage"
	"github.com/alexgunkel/feedback/util"
	"html/template"
	"net/http"
	"time"
)

type Admin struct {
	db *storage.Storage
}

func (a *Admin) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	t, err := template.ParseFiles("./templates/admin.html")
	util.PanicOnError(err)
	buf := bytes.NewBufferString("")

	if request.Method == "POST" {
		dt := request.PostFormValue("datetime")
		newSurvey := &storage.Survey{
			Title:         request.PostFormValue("new-survey"),
			AccessKey:     randomString(),
			EvaluationKey: randomString(),
		}
		t, err := time.Parse("2006-01-02T15:04", dt)
		if err == nil {
			newSurvey.DateTime = t
		} else {
			println(err.Error())
		}

		a.db.Add(newSurvey)
	}

	err = t.Execute(buf, struct {
		Surveys []storage.Survey
		Domain  string
	}{
		Surveys: a.db.GetSurveys(),
		Domain:  config.ReadConfig().Domain,
	})
	util.PanicOnError(err)
	_, err = writer.Write(buf.Bytes())
	util.PanicOnError(err)
}

func randomString() string {
	res := make([]byte, 10)
	_, err := rand.Read(res)
	if err != nil {
		panic(err)
	}
	return base32.StdEncoding.EncodeToString(res)
}

func NewAdmin(db *storage.Storage) *Admin {
	return &Admin{db: db}
}

var _ http.Handler = &Admin{}
