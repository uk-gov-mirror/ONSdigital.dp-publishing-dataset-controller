package dataset

import (
	"io/ioutil"
	"net/http"

	dphandlers "github.com/ONSdigital/dp-net/handlers"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// PatchMetadata updates dataset and version data
func PatchMetadata(dc DatasetClient) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string) {
		patchMetadata(w, r, dc, accessToken, collectionID, lang)
	})
}

func patchMetadata(w http.ResponseWriter, req *http.Request, dc DatasetClient, userAccessToken, collectionID, lang string) {
	ctx := req.Context()
	ifMatch := req.Header.Get("If-Match")

	err := checkAccessTokenAndCollectionHeaders(userAccessToken, collectionID)
	if err != nil {
		log.Error(ctx, err.Error(), err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(req)
	datasetID := vars["datasetID"]
	edition := vars["editionID"]
	version := vars["versionID"]

	logInfo := map[string]interface{}{
		"datasetID": datasetID,
		"edition":   edition,
		"version":   version,
	}

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error(ctx, "patchMetadata endpoint: error reading body", err, log.Data(logInfo))
		http.Error(w, "error reading body", http.StatusBadRequest)
		return
	}

	err = dc.PatchDataset(ctx, userAccessToken, "", collectionID, datasetID, ifMatch, req.Body)
	if err != nil {
		log.Error(ctx, "error updating (patch) dataset", err, log.Data(logInfo))
		http.Error(w, "error updating (patch) dataset", http.StatusInternalServerError)
		return
	}

	err = dc.PatchVersion(ctx, userAccessToken, "", collectionID, datasetID, edition, version, ifMatch, req.Body)
	if err != nil {
		log.Error(ctx, "error updating (patch) version", err, log.Data(logInfo))
		http.Error(w, "error updating (patch) version", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)

	log.Info(ctx, "patch metadata: request successful", log.Data(logInfo))
}
