package dataset

import (
	"encoding/json"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/dataset"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

// GetDataset returns a specfic dataset
func Get(dc *dataset.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		get(w, req, dc)
	}
}

func get(w http.ResponseWriter, req *http.Request, dc *dataset.Client) {
	ctx := req.Context()
	vars := mux.Vars(req)
	datasetID := vars["datasetID"]
	userAccessToken := ""
	collectionID := ""

	dataset, err := dc.Get(ctx, userAccessToken, "", collectionID, datasetID)
	if err != nil {
		log.Event(ctx, "error getting dataset", log.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(dataset)
	if err != nil {
		log.Event(ctx, "error marshalling json", log.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(b)
}
