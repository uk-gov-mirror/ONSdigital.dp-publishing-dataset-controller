package routes

import (
	"github.com/ONSdigital/dp-publishing-dataset-controller/config"
	"github.com/ONSdigital/dp-publishing-dataset-controller/handler"
	"github.com/ONSdigital/dp-publishing-dataset-controller/model"
	"github.com/gorilla/mux"
)

// Init initialises routes for the service
func Init(router *mux.Router, cfg *config.Config, cli model.Clients) {
	router.StrictSlash(true).Path("/datasets").HandlerFunc(handler.GetAllDatasets(cli.Dc))
	router.StrictSlash(true).Path("/datasets/{datasetID}").HandlerFunc(handler.GetDataset(cli.Dc))
}
