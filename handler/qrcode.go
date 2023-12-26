package handler

import (
	"fmt"
	"github.com/alexgunkel/feedback/config"
	"github.com/alexgunkel/feedback/storage"
	"github.com/alexgunkel/feedback/util"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

type QrCode struct {
	cfg *config.Config
	db  *storage.Storage
}

func (q *QrCode) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	v := mux.Vars(request)
	id, _ := strconv.Atoi(v["id"])
	survey := q.db.GetSurvey(id)
	key := v["key"]
	if key != survey.AccessKey {
		writer.WriteHeader(403)
		return
	}

	size := util.QrCodeSmall
	if sz := request.URL.Query().Get("size"); sz != "" {
		switch sz {
		case "big":
			size = util.QrCodeBig
		case "medium":
			size = util.QrCodeMedium
		}
	}

	address := q.cfg.Domain + survey.GetAccessPath()
	writer.Header().Set("Content-Type", "image/png")

	writer.Header().Set("Cache-Control", "no-store, no-cache")
	writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.png\"", strings.Replace(survey.Title, " ", "_", -1)))
	_, err := writer.Write(util.CreateCode(address, size))
	if err != nil {
		return
	}
}

func NewQrCodeHandler(cfg *config.Config, db *storage.Storage) *QrCode {
	return &QrCode{cfg: cfg, db: db}
}

var _ http.Handler = &QrCode{}
