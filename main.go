package main

import (
	"os"

	"github.com/ONSdigital/dp-api-clients-go/dataset"
	zebedee "github.com/ONSdigital/dp-api-clients-go/zebedee"
	"github.com/ONSdigital/dp-publishing-dataset-controller/config"
	"github.com/ONSdigital/dp-publishing-dataset-controller/routes"
	"github.com/ONSdigital/go-ns/server"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

func main() {
	cfg, err := config.Get()
	if err != nil {
		log.Event(nil, "error getting configuration", log.Error(err))
		os.Exit(1)
	}

	log.Namespace = "dp-publishing-dataset-controller"
	log.Event(nil, "config on startup", log.Data{"config": cfg})

	router := mux.NewRouter()

	dc := dataset.NewAPIClient(cfg.DatasetAPIURL)
	zc := zebedee.NewZebedeeClient(cfg.ZebedeeURL)

	routes.Init(router, cfg, dc, zc)

	s := server.New(cfg.BindAddr, router)

	if err := s.ListenAndServe(); err != nil {
		log.Event(nil, "error starting http server", log.Error(err))
		os.Exit(1)
	}
}
