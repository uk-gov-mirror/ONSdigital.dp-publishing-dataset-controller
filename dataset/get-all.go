package dataset

import (
	"encoding/json"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/dataset"
	"github.com/ONSdigital/dp-publishing-dataset-controller/mapper"
	"github.com/ONSdigital/log.go/log"
)

// GetAllDatasets returns a mapped list of all datasets
func GetAllDatasets(dc *dataset.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		getAllDatasets(w, req, dc)
	}
}

func getAllDatasets(w http.ResponseWriter, req *http.Request, dc *dataset.Client) {
	ctx := req.Context()
	// userAccessToken := getUserAccessTokenFromContent(ctx)
	// collectionID := getCollectionIDFromContext(ctx)
	userAccessToken := ""
	collectionID := ""

	log.Event(ctx, "calling get datasets")

	datasets, err := dc.GetDatasets(ctx, userAccessToken, "", collectionID)
	if err != nil {
		log.Event(ctx, "error getting all datasets", log.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	mapped := mapper.AllDatasets(datasets)

	b, err := json.Marshal(mapped)
	if err != nil {
		log.Event(ctx, "error marshalling json", log.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(b)
}
