package main

import (
	"os"

	"github.com/ONSdigital/dp-publishing-dataset-controller/config"
	"github.com/ONSdigital/go-ns/log"
	"github.com/ONSdigital/go-ns/server"
	"github.com/gorilla/mux"
)

func main() {
	cfg, err := config.Get()
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}

	log.Namespace = "dp-publishing-dataset-controller"
	log.Info("config on startup", log.Data{"config": cfg})

	router := mux.NewRouter()

	s := server.New(cfg.BindAddr, router)

	if err := s.ListenAndServe(); err != nil {
		log.Error(err, nil)
		os.Exit(2)
	}
}
