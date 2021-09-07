package dataset

import (
	"encoding/json"
	"fmt"
	"net/http"

	dphandlers "github.com/ONSdigital/dp-net/handlers"
	"github.com/ONSdigital/dp-publishing-dataset-controller/mapper"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// GetVersions returns a mapped list of all versions
func GetVersions(dc DatasetClient, batchSize, maxWorkers int) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string) {
		getVersions(w, r, dc, accessToken, collectionID, lang, batchSize, maxWorkers)
	})
}

func getVersions(w http.ResponseWriter, req *http.Request, dc DatasetClient, userAccessToken, collectionID, lang string, batchSize, maxWorkers int) {
	ctx := req.Context()

	vars := mux.Vars(req)
	datasetID := vars["datasetID"]
	editionID := vars["editionID"]

	logInfo := map[string]interface{}{
		"datasetID": datasetID,
		"edition":   editionID,
	}

	err := checkAccessTokenAndCollectionHeaders(userAccessToken, collectionID)
	if err != nil {
		log.Error(ctx, err.Error(), err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Info(ctx, "calling get versions", log.Data(logInfo))

	dataset, err := dc.GetDatasetCurrentAndNext(ctx, userAccessToken, "", collectionID, datasetID)
	if err != nil {
		errMsg := fmt.Sprintf("error getting dataset from dataset API: %v", err.Error())
		log.Error(ctx, "error getting dataset from dataset API", err, log.Data(logInfo))
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	edition, err := dc.GetEdition(ctx, userAccessToken, "", collectionID, datasetID, editionID)
	if err != nil {
		errMsg := fmt.Sprintf("error getting edition from dataset API: %v", err.Error())
		log.Error(ctx, "error getting edition from dataset API", err, log.Data(logInfo))
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	versions, err := dc.GetVersionsInBatches(ctx, userAccessToken, "", "", collectionID, datasetID, editionID, batchSize, maxWorkers)
	if err != nil {
		errMsg := fmt.Sprintf("error getting all versions from dataset API: %v", err.Error())
		log.Error(ctx, "error getting all versions from dataset API", err, log.Data(logInfo))
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	mapped := mapper.AllVersions(ctx, dataset, edition, versions)

	b, err := json.Marshal(mapped)
	if err != nil {
		log.Error(ctx, "error marshalling response to json", err, log.Data(logInfo))
		http.Error(w, "error marshalling response to json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)

	log.Info(ctx, "get versions: request successful", log.Data(logInfo))
}
