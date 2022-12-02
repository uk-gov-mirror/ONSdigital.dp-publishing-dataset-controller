package dataset

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	dphandlers "github.com/ONSdigital/dp-net/handlers"
	"github.com/ONSdigital/dp-publishing-dataset-controller/model"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// PatchDatasetMetadata updates dataset metadata
func UpdateDatasetMetadata(dc DatasetClient, zc ZebedeeClient) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string) {
		updateDatasetMetadata(w, r, dc, zc, accessToken, collectionID, lang)
	})
}

func updateDatasetMetadata(w http.ResponseWriter, req *http.Request, dc DatasetClient, zc ZebedeeClient, userAccessToken, collectionID, lang string) {
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
	logData := log.Data{"dataset_id": datasetID}

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error(ctx, "datasetPatchMetadata endpoint: error reading body", err, logData)
		http.Error(w, "error reading body", http.StatusBadRequest)
		return
	}

	var body model.UpdateMetadata
	if err = json.Unmarshal(b, &body); err != nil {
		log.Error(ctx, "datasetPatchMetadata endpoint: error unmarshalling body", err, logData)
		http.Error(w, "error unmarshalling body", http.StatusBadRequest)
		return
	}

	err = dc.PatchDataset(ctx, userAccessToken, "", collectionID, datasetID, ifMatch, body.Patches)
	if err != nil {
		log.Error(ctx, "error updating dataset metadata via patch endpoint", err, logData)
		http.Error(w, "error updating dataset metadata via patch endpoint", http.StatusInternalServerError)
		return
	}
	err = zc.PutDatasetInCollection(ctx, userAccessToken, collectionID, "", datasetID, body.CollectionState)
	if err != nil {
		log.Error(ctx, "error adding dataset to collection", err, logData)
		http.Error(w, "error adding dataset to collection", http.StatusInternalServerError)
		return
	}

	datasetPatches, err := json.Marshal(body.Patches)
	if err != nil {
		log.Error(ctx, "error marshalling dataset patches to json", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(datasetPatches)
	log.Info(ctx, "patch dataset metadata: request successful", logData)
}

// UpdateDatasetVersionMetadata updates dataset version metadata
func UpdateDatasetVersionMetadata(dc DatasetClient, zc ZebedeeClient) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string) {
		updateDatasetVersionMetadata(w, r, dc, zc, accessToken, collectionID, lang)
	})
}

func updateDatasetVersionMetadata(w http.ResponseWriter, req *http.Request, dc DatasetClient, zc ZebedeeClient, userAccessToken, collectionID, lang string) {
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

	logData := log.Data{
		"datasetID": datasetID,
		"edition":   edition,
		"version":   version,
	}

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error(ctx, "datasetVersionPatchMetadata endpoint: error reading body", err, logData)
		http.Error(w, "error reading body", http.StatusBadRequest)
		return
	}

	var body model.UpdateMetadata
	if err = json.Unmarshal(b, &body); err != nil {
		log.Error(ctx, "datasetVersionPatchMetadata endpoint: error unmarshalling body", err, logData)
		http.Error(w, "error unmarshalling body", http.StatusBadRequest)
		return
	}

	err = dc.PatchVersion(ctx, userAccessToken, "", collectionID, datasetID, edition, version, ifMatch, body.Patches)
	if err != nil {
		log.Error(ctx, "error updating dataset version metadata via patch endpoint", err, logData)
		http.Error(w, "error updating dataset version metadata via patch endpoint", http.StatusInternalServerError)
		return
	}

	err = zc.PutDatasetVersionInCollection(ctx, userAccessToken, collectionID, "", datasetID, edition, version, body.CollectionState)
	if err != nil {
		log.Error(ctx, "error adding version to collection", err, logData)
		http.Error(w, "error adding version to collection", http.StatusInternalServerError)
		return
	}

	datasetVersionPatches, err := json.Marshal(body.Patches)
	if err != nil {
		log.Error(ctx, "error marshalling dataset version patches to json", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(datasetVersionPatches)
	log.Info(ctx, "patch dataset version metadata: request successful", logData)
}
