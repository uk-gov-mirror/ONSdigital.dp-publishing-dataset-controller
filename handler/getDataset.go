package handler

import (
	"net/http"
	"encoding/json"

	"github.com/ONSdigital/dp-api-clients-go/dataset"
	"github.com/ONSdigital/log.go/log"
    "github.com/gorilla/mux"
)

// GetDataset returns a specfic dataset
func GetDataset(dc *dataset.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		getDataset(w, req, dc)
	}
}

func getDataset(w http.ResponseWriter, req *http.Request, dc *dataset.Client) {
	ctx := req.Context()
	vars := mux.Vars(req)
    datasetID := vars["datasetID"]
	userAccessToken := ""
	collectionID := ""

	datasets, err := dc.Get(ctx, userAccessToken, "", collectionID, datasetID)
	if err != nil {
		log.Event(nil, "error getting dataset", log.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	b, _ := json.Marshal(datasets)
	w.Write(b)
}
