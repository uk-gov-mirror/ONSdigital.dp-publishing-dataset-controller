package routes

import (
	ds "github.com/ONSdigital/dp-api-clients-go/dataset"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-publishing-dataset-controller/config"
	"github.com/ONSdigital/dp-publishing-dataset-controller/dataset"
	"github.com/ONSdigital/dp-publishing-dataset-controller/handlers"
	"github.com/gorilla/mux"
)

// Init initialises routes for the service
func Init(router *mux.Router, cfg *config.Config, hc healthcheck.HealthCheck, dc *ds.Client) {
	router.StrictSlash(true).Path("/health").HandlerFunc(hc.Handler)

	router.StrictSlash(true).Path("/datasets").HandlerFunc(dataset.GetAll(dc))
	router.StrictSlash(true).Path("/datasets/{datasetID}/editions/{editionID}/versions/{versionID}").HandlerFunc(handlers.GetMetadataHandler(dc, *cfg))
}
