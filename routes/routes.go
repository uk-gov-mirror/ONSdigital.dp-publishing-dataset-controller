package routes

import (
	"github.com/ONSdigital/dp-api-clients-go/dataset"
	zebedee "github.com/ONSdigital/dp-api-clients-go/zebedee"
	"github.com/ONSdigital/dp-publishing-dataset-controller/config"
	"github.com/ONSdigital/dp-publishing-dataset-controller/handler"
	"github.com/gorilla/mux"
)

// Init initialises routes for the service
func Init(router *mux.Router, cfg *config.Config, dc *dataset.Client, zc *zebedee.ZebedeeClient) {
	router.StrictSlash(true).Path("/datasets").HandlerFunc(handler.GetAllDatasets(dc))
	router.StrictSlash(true).Path("/datasets/{datasetID}").HandlerFunc(handler.GetDataset(dc))
}
