package dataset

import (
	"encoding/json"
	"fmt"
	"net/http"

	dphandlers "github.com/ONSdigital/dp-net/handlers"
	"github.com/ONSdigital/dp-publishing-dataset-controller/mapper"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

// GetEditions returns a mapped list of all editions
func GetEditions(dc DatasetClient) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string) {
		getEditions(w, r, dc, accessToken, collectionID, lang)
	})
}

func getEditions(w http.ResponseWriter, req *http.Request, dc DatasetClient, userAccessToken, collectionID, lang string) {
	ctx := req.Context()

	vars := mux.Vars(req)
	datasetID := vars["datasetID"]

	err := checkAccessTokenAndCollectionHeaders(userAccessToken, collectionID)
	if err != nil {
		log.Event(ctx, err.Error(), log.ERROR)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logInfo := map[string]interface{}{
		"datasetID":    datasetID,
		"collectionID": collectionID,
	}

	log.Event(ctx, "calling get editions", log.INFO, log.Data(logInfo))

	dataset, err := dc.GetDatasetCurrentAndNext(ctx, userAccessToken, "", collectionID, datasetID)
	if err != nil {
		errMsg := fmt.Sprintf("error getting dataset from dataset API: %v", err.Error())
		log.Event(ctx, "error getting dataset from dataset API", log.ERROR, log.Error(err), log.Data(logInfo))
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	editions, err := dc.GetEditions(ctx, userAccessToken, "", collectionID, datasetID)
	if err != nil {
		errMsg := fmt.Sprintf("error getting editions from dataset API: %v", err.Error())
		log.Event(ctx, "error getting editions from dataset API", log.ERROR, log.Error(err), log.Data(logInfo))
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	latestVersionInEdition := make(map[string]string)
	for _, edition := range editions {
		_, _, versionID, err := getIDsFromURL(edition.Links.LatestVersion.URL)
		if err != nil {
			latestVersionInEdition[edition.Edition] = ""
			continue
		}
		version, err := dc.GetVersion(ctx, userAccessToken, "", "", collectionID, datasetID, edition.Edition, versionID)
		if err != nil {
			latestVersionInEdition[edition.Edition] = ""
			continue
		}
		latestVersionInEdition[edition.Edition] = version.ReleaseDate
	}

	mapped := mapper.AllEditions(ctx, dataset, editions, latestVersionInEdition)

	b, err := json.Marshal(mapped)
	if err != nil {
		log.Event(ctx, "error marshalling response to json", log.ERROR, log.Error(err), log.Data(logInfo))
		http.Error(w, "error marshalling response to json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)

	log.Event(ctx, "get versions: request successful", log.INFO, log.Data(logInfo))
}
