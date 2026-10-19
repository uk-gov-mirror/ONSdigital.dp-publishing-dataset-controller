package dataset

import (
	"encoding/json"
	"fmt"
	"net/http"

	datasetApiSdk "github.com/ONSdigital/dp-dataset-api/sdk"
	dphandlers "github.com/ONSdigital/dp-net/v3/handlers"
	"github.com/ONSdigital/dp-publishing-dataset-controller/mapper"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// GetEditions returns a mapped list of all editions
func GetEditions(dc DatasetAPIClient) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string) {
		getEditions(w, r, dc, accessToken, collectionID, lang)
	})
}

func getEditions(w http.ResponseWriter, req *http.Request, dc DatasetAPIClient, userAccessToken, collectionID, lang string) {
	ctx := req.Context()

	vars := mux.Vars(req)
	datasetID := vars[datasetIDKey]

	err := checkAccessTokenAndCollectionHeaders(userAccessToken, collectionID)
	if err != nil {
		log.Error(ctx, err.Error(), err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logInfo := map[string]interface{}{
		datasetIDKey:   datasetID,
		"collectionID": collectionID,
	}

	headers := datasetApiSdk.Headers{
		CollectionID: collectionID,
		AccessToken:  userAccessToken,
	}

	log.Info(ctx, "calling get editions", log.Data(logInfo))

	dataset, err := dc.GetDatasetCurrentAndNext(ctx, headers, datasetID)
	if err != nil {
		errMsg := fmt.Sprintf("error getting dataset from dataset API: %v", err.Error())
		log.Error(ctx, "error getting dataset from dataset API", err, log.Data(logInfo))
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	editions, err := dc.GetEditions(ctx, headers, datasetID, nil)
	if err != nil {
		errMsg := fmt.Sprintf("error getting editions from dataset API: %v", err.Error())
		log.Error(ctx, "error getting editions from dataset API", err, log.Data(logInfo))
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	latestVersionInEdition := make(map[string]string)
	for e := range editions.Items {
		_, _, versionID, err := getIDsFromURL(editions.Items[e].Links.LatestVersion.HRef)
		if err != nil {
			latestVersionInEdition[editions.Items[e].Edition] = ""
			continue
		}

		version, err := dc.GetVersion(ctx, headers, datasetID, editions.Items[e].Edition, versionID)
		if err != nil {
			latestVersionInEdition[editions.Items[e].Edition] = ""
			continue
		}
		latestVersionInEdition[editions.Items[e].Edition] = version.ReleaseDate
	}

	mapped := mapper.AllEditions(ctx, dataset, editions, latestVersionInEdition)

	b, err := json.Marshal(mapped)
	if err != nil {
		log.Error(ctx, "error marshalling editions response to json", err, log.Data(logInfo))
		http.Error(w, "error marshalling editions response to json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(b)
	if err != nil {
		log.Error(ctx, "error writing response", err)
		http.Error(w, "error writing response", http.StatusInternalServerError)
		return
	}

	log.Info(ctx, "get editions: request successful", log.Data(logInfo))
}
