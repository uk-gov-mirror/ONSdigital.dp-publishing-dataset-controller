package dataset

import (
	"encoding/json"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/dataset"
	"github.com/ONSdigital/dp-api-clients-go/headers"
	"github.com/ONSdigital/dp-publishing-dataset-controller/mapper"
	"github.com/ONSdigital/log.go/log"
)

// GetAllDatasets returns a mapped list of all datasets
func GetAll(dc *dataset.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		getAll(w, req, dc)
	}
}

func getAll(w http.ResponseWriter, req *http.Request, dc *dataset.Client) {
	ctx := req.Context()

	userAccessToken, err := headers.GetUserAuthToken(req)
	if err != nil && err != headers.ErrHeaderNotFound {
		log.Event(ctx, "error getting user access token from header", log.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	collectionID, err := headers.GetCollectionID(req)
	if err != nil && err != headers.ErrHeaderNotFound {
		log.Event(ctx, "error getting collection ID from header", log.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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
