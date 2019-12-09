package main

import (
	"os"

	"github.com/ONSdigital/dp-publishing-dataset-controller/config"
	"github.com/ONSdigital/go-ns/log"
)

func main() {
	log.Namespace = "dp-publishing-dataset-controller"

	cfg, err := config.Get()
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}

	log.Info("config on startup", log.Data{"config": cfg})
}
