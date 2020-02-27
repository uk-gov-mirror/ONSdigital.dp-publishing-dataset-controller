package dataset

import (
	"encoding/json"
	"net/http"

	//datasetclient "github.com/ONSdigital/dp-api-clients-go/dataset"
	"github.com/ONSdigital/dp-api-clients-go/headers"
	"github.com/ONSdigital/dp-publishing-dataset-controller/mapper"
	"github.com/ONSdigital/log.go/log"
)

// GetAll returns a mapped list of all datasets
func GetAll(dc Client) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		getAll(w, req, dc)
	}
}

func getAll(w http.ResponseWriter, req *http.Request, dc Client) {
	ctx := req.Context()

	userAccessToken, err := headers.GetUserAuthToken(req)
	if err == headers.ErrHeaderNotFound {
		log.Event(ctx, "no user access token header set", log.ERROR, log.Error(err))
		http.Error(w, "no user access token header set", http.StatusBadRequest)
		return
	}
	if err != nil {
		log.Event(ctx, "error getting user access token from header", log.ERROR, log.Error(err))
		http.Error(w, "error getting user access token from header", http.StatusBadRequest)
		return
	}

	collectionID, err := headers.GetCollectionID(req)
	if err == headers.ErrHeaderNotFound {
		log.Event(ctx, "no collection ID header set", log.ERROR, log.Error(err))
		http.Error(w, "no collection ID header set", http.StatusBadRequest)
		return
	}
	if err != nil {
		log.Event(ctx, "error getting collection ID from header", log.ERROR, log.Error(err))
		http.Error(w, "error getting collection ID from header", http.StatusBadRequest)
		return
	}
	log.Event(ctx, "calling get datasets")

	datasets, err := dc.GetDatasets(ctx, userAccessToken, "", collectionID)
	if err != nil {
		log.Event(ctx, "error getting all datasets from dataset API", log.ERROR, log.Error(err))
		http.Error(w, "error getting all datasets from dataset API", http.StatusInternalServerError)
		return
	}

	mapped := mapper.AllDatasets(datasets)

	b, err := json.Marshal(mapped)
	if err != nil {
		log.Event(ctx, "error marshalling response to json", log.ERROR, log.Error(err))
		http.Error(w, "error marshalling response to json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
