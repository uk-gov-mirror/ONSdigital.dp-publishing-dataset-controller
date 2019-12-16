package routes

import (
	ds "github.com/ONSdigital/dp-api-clients-go/dataset"
	zebedee "github.com/ONSdigital/dp-api-clients-go/zebedee"
	"github.com/ONSdigital/dp-publishing-dataset-controller/config"
	"github.com/ONSdigital/dp-publishing-dataset-controller/dataset"
	"github.com/gorilla/mux"
)

// Init initialises routes for the service
func Init(router *mux.Router, cfg *config.Config, dc *ds.Client, zc *zebedee.ZebedeeClient) {
	router.StrictSlash(true).Path("/datasets").HandlerFunc(dataset.GetAll(dc))
	router.StrictSlash(true).Path("/datasets/{datasetID}").HandlerFunc(dataset.Get(dc))
}
