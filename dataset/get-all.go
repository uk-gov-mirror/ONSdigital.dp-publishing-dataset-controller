package dataset

import (
	"encoding/json"
	"net/http"

	//datasetclient "github.com/ONSdigital/dp-api-clients-go/dataset"

	dphandlers "github.com/ONSdigital/dp-net/handlers"
	"github.com/ONSdigital/dp-publishing-dataset-controller/mapper"
	"github.com/ONSdigital/log.go/log"
)

// GetAll returns a mapped list of all datasets
func GetAll(dc DatasetClient, batchSize, maxWorkers int) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string) {
		getAll(w, r, dc, accessToken, collectionID, lang, batchSize, maxWorkers)
	})
}

func getAll(w http.ResponseWriter, req *http.Request, dc DatasetClient, userAccessToken, collectionID, lang string, batchSize, maxWorkers int) {
	ctx := req.Context()

	err := checkAccessTokenAndCollectionHeaders(userAccessToken, collectionID)
	if err != nil {
		log.Event(ctx, err.Error(), log.ERROR)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Event(ctx, "calling get datasets")

	datasets, err := dc.GetDatasetsInBatches(ctx, userAccessToken, "", collectionID, batchSize, maxWorkers)
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

	log.Event(ctx, "get all: request successful", log.INFO)
}
