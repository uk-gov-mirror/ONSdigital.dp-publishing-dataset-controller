package routes

import (
	"net/http"

	ds "github.com/ONSdigital/dp-api-clients-go/dataset"
	zc "github.com/ONSdigital/dp-api-clients-go/zebedee"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	bc "github.com/ONSdigital/dp-publishing-dataset-controller/clients/topics"
	"github.com/ONSdigital/dp-publishing-dataset-controller/config"
	"github.com/ONSdigital/dp-publishing-dataset-controller/dataset"
	"github.com/gorilla/mux"
)

// Init initialises routes for the service
func Init(router *mux.Router, cfg *config.Config, hc healthcheck.HealthCheck, dc *ds.Client, zc *zc.Client, bc *bc.Client) {
	router.StrictSlash(true).Path("/health").HandlerFunc(hc.Handler)

	router.StrictSlash(true).Path("/datasets").HandlerFunc(dataset.GetAll(dc, cfg.DatasetsBatchSize, cfg.DatasetsBatchWorkers)).Methods(http.MethodGet)
	router.StrictSlash(true).Path("/datasets/{datasetID}/create").HandlerFunc(dataset.GetTopics(bc)).Methods(http.MethodGet)
	router.StrictSlash(true).Path("/datasets/{datasetID}/editions/{editionID}/versions/{versionID}").HandlerFunc(dataset.GetMetadataHandler(dc, zc)).Methods(http.MethodGet)
	router.StrictSlash(true).Path("/datasets/{datasetID}/editions/{editionID}/versions/{versionID}").HandlerFunc(dataset.PutMetadata(dc, zc)).Methods(http.MethodPut)
}
