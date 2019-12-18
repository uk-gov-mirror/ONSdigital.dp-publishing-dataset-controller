package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ONSdigital/dp-api-clients-go/dataset"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-publishing-dataset-controller/config"
	"github.com/ONSdigital/dp-publishing-dataset-controller/routes"
	"github.com/ONSdigital/go-ns/server"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

var version string

func main() {
	log.Namespace = "dp-publishing-dataset-controller"

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	cfg, err := config.Get()
	if err != nil {
		log.Event(nil, "error getting configuration", log.Error(err))
		os.Exit(1)
	}
	log.Event(nil, "config on startup", log.Data{"config": cfg})

	router := mux.NewRouter()

	dc := dataset.NewAPIClient(cfg.DatasetAPIURL)

	hc := healthcheck.Create(version, cfg.HealthCheckCritialTimeout, cfg.HealthCheckInterval)

	routes.Init(router, cfg, hc, dc)

	s := server.New(cfg.BindAddr, router)
	hc.Start(nil)

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Event(nil, "error starting http server", log.Error(err))
			os.Exit(1)
		}
	}()

	for {
		select {
		case <-signals:
			log.Event(nil, "os signal received")
			gracefulShutdown(cfg, s, hc)
		}
	}
}

func gracefulShutdown(cfg *config.Config, s *server.Server, hc healthcheck.HealthCheck) {
	log.Event(nil, fmt.Sprintf("shutdown with timeout: %s", cfg.GracefulShutdownTimeout))
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GracefulShutdownTimeout)
	if err := s.Shutdown(ctx); err != nil {
		log.Event(nil, "failed to gracefully shutdown http server", log.Error(err))
	}
	log.Event(nil, "graceful shutdown of http server complete", nil)
	hc.Stop()
	log.Event(nil, "shutdown complete", nil)
	cancel()
	os.Exit(1)
}
