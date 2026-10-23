package dataset

import (
	"encoding/json"
	"net/http"

	datasetApiSdk "github.com/ONSdigital/dp-dataset-api/sdk"

	dphandlers "github.com/ONSdigital/dp-net/v3/handlers"
	"github.com/ONSdigital/dp-publishing-dataset-controller/mapper"
	"github.com/ONSdigital/log.go/v2/log"
)

const (
	datasetIDKey = "datasetID"
	editionKey   = "edition"
)

// GetAll returns a mapped list of all datasets
func GetAll(dc DatasetAPIClient, batchSize, maxWorkers int) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string) {
		getAll(w, r, dc, accessToken, collectionID, lang, batchSize, maxWorkers)
	})
}

func getAll(w http.ResponseWriter, req *http.Request, dc DatasetAPIClient, userAccessToken, collectionID, lang string, batchSize, maxWorkers int) {
	ctx := req.Context()

	err := checkAccessTokenAndCollectionHeaders(userAccessToken, collectionID)
	if err != nil {
		log.Error(ctx, err.Error(), err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	headers := datasetApiSdk.Headers{
		CollectionID: collectionID,
		AccessToken:  userAccessToken,
	}

	log.Info(ctx, "calling get datasets")

	datasets, err := dc.GetDatasetsInBatches(ctx, headers, batchSize, maxWorkers)
	if err != nil {
		log.Error(ctx, "error getting all datasets from dataset API", err)
		http.Error(w, "error getting all datasets from dataset API", http.StatusInternalServerError)
		return
	}

	mapped := mapper.AllDatasets(datasets)

	b, err := json.Marshal(mapped)
	if err != nil {
		log.Error(ctx, "error marshalling response to json", err)
		http.Error(w, "error marshalling response to json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(b)
	if err != nil {
		log.Error(ctx, "error writing response", err)
		http.Error(w, "error writing response", http.StatusInternalServerError)
		return
	}

	log.Info(ctx, "get all: request successful")
}
