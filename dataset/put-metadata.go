package dataset

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	datasetclient "github.com/ONSdigital/dp-api-clients-go/dataset"
	dphandlers "github.com/ONSdigital/dp-net/handlers"
	"github.com/ONSdigital/dp-publishing-dataset-controller/model"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

// PutMetadata updates dataset, version and dimension data
func PutMetadata(dc DatasetClient, zc ZebedeeClient) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string) {
		putMetadata(w, r, dc, zc, accessToken, collectionID, lang)
	})
}

func putMetadata(w http.ResponseWriter, req *http.Request, dc DatasetClient, zc ZebedeeClient, userAccessToken, collectionID, lang string) {
	ctx := req.Context()

	err := checkAccessTokenAndCollectionHeaders(userAccessToken, collectionID)
	if err != nil {
		log.Event(ctx, err.Error(), log.ERROR)
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
		log.Event(ctx, "putMetadata endpoint: error reading body", log.ERROR, log.Error(err), log.Data(logInfo))
		http.Error(w, "error reading body", http.StatusBadRequest)
		return
	}

	var body model.EditMetadata
	if err = json.Unmarshal(b, &body); err != nil {
		log.Event(ctx, "putMetadata endpoint: error unmarshalling body", log.ERROR, log.Error(err), log.Data(logInfo))
		http.Error(w, "error unmarshalling body", http.StatusBadRequest)
		return
	}

	err = dc.PutDataset(ctx, userAccessToken, "", collectionID, datasetID, body.Dataset)
	if err != nil {
		log.Event(ctx, "error updating dataset", log.ERROR, log.Error(err), log.Data(logInfo))
		http.Error(w, "error updating dataset", http.StatusInternalServerError)
		return
	}

	err = dc.PutVersion(ctx, userAccessToken, "", collectionID, datasetID, edition, version, body.Version)
	if err != nil {
		log.Event(ctx, "error updating version", log.ERROR, log.Error(err), log.Data(logInfo))
		http.Error(w, "error updating version", http.StatusInternalServerError)
		return
	}

	instance := datasetclient.Instance{}
	instance.Dimensions = body.Dimensions

	err = dc.PutInstance(ctx, userAccessToken, "", collectionID, body.Version.ID, instance)
	if err != nil {
		log.Event(ctx, "error updating dimensions", log.ERROR, log.Error(err), log.Data(logInfo))
		http.Error(w, "error updating dimensions", http.StatusInternalServerError)
		return
	}

	err = zc.PutDatasetInCollection(ctx, userAccessToken, "", collectionID, datasetID, body.CollectionState)
	if err != nil {
		log.Event(ctx, "error adding dataset to collection", log.ERROR, log.Error(err), log.Data(logInfo))
		http.Error(w, "error adding dataset to collection", http.StatusInternalServerError)
		return
	}

	err = zc.PutDatasetVersionInCollection(ctx, userAccessToken, "", collectionID, datasetID, edition, version, body.CollectionState)
	if err != nil {
		log.Event(ctx, "error adding version to collection", log.ERROR, log.Error(err), log.Data(logInfo))
		http.Error(w, "error adding version to collection", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)

	log.Event(ctx, "put metadata: request successful", log.INFO, log.Data(logInfo))
}
