package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ONSdigital/dp-api-clients-go/dataset"
	"github.com/ONSdigital/dp-api-clients-go/health"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dpnethttp "github.com/ONSdigital/dp-net/http"
	"github.com/ONSdigital/dp-publishing-dataset-controller/config"
	"github.com/ONSdigital/dp-publishing-dataset-controller/routes"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

// App version informaton retrieved on runtime
var (
	// BuildTime represents the time in which the service was built
	BuildTime string
	// GitCommit represents the commit (SHA-1) hash of the service that is running
	GitCommit string
	// Version represents the version of the service that is running
	Version string
)

func main() {
	ctx := context.Background()

	log.Namespace = "dp-publishing-dataset-controller"

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	cfg, err := config.Get()
	if err != nil {
		log.Event(ctx, "error getting configuration", log.FATAL, log.Error(err))
		os.Exit(1)
	}
	log.Event(ctx, "config on startup", log.INFO, log.Data{"config": cfg})

	versionInfo, err := healthcheck.NewVersionInfo(
		BuildTime,
		GitCommit,
		Version,
	)
	if err != nil {
		log.Event(ctx, "failed to create service version information", log.FATAL, log.Error(err))
		os.Exit(1)
	}

	apiRouterCli := health.NewClient("api-router", cfg.APIRouterURL)
	dc := dataset.NewWithHealthClient(apiRouterCli)

	hc := healthcheck.New(versionInfo, cfg.HealthCheckCritialTimeout, cfg.HealthCheckInterval)
	if err = hc.AddCheck("dataset API", dc.Checker); err != nil {
		log.Event(ctx, "failed to add dataset API checker", log.FATAL, log.Error(err))
		os.Exit(1)
	}

	router := mux.NewRouter()
	routes.Init(router, cfg, hc, dc)

	s := dpnethttp.NewServer(cfg.BindAddr, router)

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Event(ctx, "error starting http server", log.ERROR, log.Error(err))
			return
		}
	}()

	hc.Start(ctx)

	// Block until a fatal error occurs
	select {
	case signal := <-signals:
		log.Event(ctx, "quitting after os signal received", log.INFO, log.Data{"signal": signal})
	}

	log.Event(ctx, fmt.Sprintf("shutdown with timeout: %s", cfg.GracefulShutdownTimeout), log.INFO)

	// give the app `Timeout` seconds to close gracefully before killing it.
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GracefulShutdownTimeout)

	go func() {
		log.Event(ctx, "stop health checkers", log.INFO)
		hc.Stop()

		if err := s.Shutdown(ctx); err != nil {
			log.Event(ctx, "failed to gracefully shutdown http server", log.ERROR, log.Error(err))
		}

		cancel() // stop timer
	}()

	// wait for timeout or success (via cancel)
	<-ctx.Done()
	if ctx.Err() == context.DeadlineExceeded {
		log.Event(ctx, "context deadline exceeded", log.WARN, log.Error(ctx.Err()))
	} else {
		log.Event(ctx, "graceful shutdown complete", log.INFO, log.Data{"context": ctx.Err()})
	}

	os.Exit(0)
}
