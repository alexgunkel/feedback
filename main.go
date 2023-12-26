package main

import (
	"github.com/gorilla/mux"
	"net/http"

	"github.com/alexgunkel/feedback/config"
	"github.com/alexgunkel/feedback/handler"
	"github.com/alexgunkel/feedback/storage"
	"github.com/alexgunkel/feedback/util"
)

func main() {
	c := make(chan storage.Input, 10)
	cfg := config.ReadConfig()
	s := storage.NewStorage(cfg)
	go s.HandleRatings(c)
	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("./public/"))
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", fs))
	r.Handle("/result/{id:[0-9]+}/{key}", handler.NewStatistics(s, cfg))
	r.Handle("/admin", handler.NewAdmin(s))
	r.Handle("/{id:[0-9]+}/{key}", handler.NewFeedback(cfg, s)).Methods("GET")
	r.Handle("/qr/{id:[0-9]+}/{key}", handler.NewQrCodeHandler(cfg, s)).Methods("GET")
	r.Handle("/{id:[0-9]+}/{key}", handler.NewStore(cfg, c)).Methods("POST")

	srv := &http.Server{
		Addr:    cfg.Listen,
		Handler: r,
	}
	util.PanicOnError(srv.ListenAndServe())
}
