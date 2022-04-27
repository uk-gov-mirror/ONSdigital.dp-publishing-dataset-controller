package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/health"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dpnethttp "github.com/ONSdigital/dp-net/http"
	"github.com/ONSdigital/dp-publishing-dataset-controller/clients/topics"
	"github.com/ONSdigital/dp-publishing-dataset-controller/config"
	"github.com/ONSdigital/dp-publishing-dataset-controller/routes"
	"github.com/ONSdigital/log.go/v2/log"
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
		log.Fatal(ctx, "error getting configuration", err)
		os.Exit(1)
	}
	log.Info(ctx, "config on startup", log.Data{"config": cfg})

	versionInfo, err := healthcheck.NewVersionInfo(
		BuildTime,
		GitCommit,
		Version,
	)
	if err != nil {
		log.Fatal(ctx, "failed to create service version information", err)
		os.Exit(1)
	}

	apiRouterCli := health.NewClient("api-router", cfg.APIRouterURL)
	dc := dataset.NewWithHealthClient(apiRouterCli)
	zc := zebedee.NewWithHealthClient(apiRouterCli)
	bc := topics.New(cfg.BabbageURL)

	hc := healthcheck.New(versionInfo, cfg.HealthCheckCritialTimeout, cfg.HealthCheckInterval)
	if err = hc.AddCheck("API router", apiRouterCli.Checker); err != nil {
		log.Fatal(ctx, "failed to add dataset API checker", err)
		os.Exit(1)
	}

	router := mux.NewRouter()
	routes.Init(router, cfg, hc, dc, zc, bc)

	s := dpnethttp.NewServer(cfg.BindAddr, router)

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Error(ctx, "error starting http server", err)
			return
		}
	}()

	hc.Start(ctx)

	// Block until a fatal error occurs
	select {
	case signal := <-signals:
		log.Info(ctx, "quitting after os signal received", log.Data{"signal": signal})
	}

	log.Info(ctx, fmt.Sprintf("shutdown with timeout: %s", cfg.GracefulShutdownTimeout))

	// give the app `Timeout` seconds to close gracefully before killing it.
	ctx, cancel := context.WithTimeout(context.Background(), cfg.GracefulShutdownTimeout)

	go func() {
		log.Info(ctx, "stop health checkers")
		hc.Stop()

		if err := s.Shutdown(ctx); err != nil {
			log.Error(ctx, "failed to gracefully shutdown http server", err)
		}

		cancel() // stop timer
	}()

	// wait for timeout or success (via cancel)
	<-ctx.Done()
	if ctx.Err() == context.DeadlineExceeded {
		log.Warn(ctx, "context deadline exceeded", log.FormatErrors([]error{ctx.Err()}))
	} else {
		log.Info(ctx, "graceful shutdown complete", log.Data{"context": ctx.Err()})
	}

	os.Exit(0)
}
