package dataset

import (
	"encoding/json"
	"io"
	"net/http"

	datasetApiSdk "github.com/ONSdigital/dp-dataset-api/sdk"
	dphandlers "github.com/ONSdigital/dp-net/v3/handlers"
	"github.com/ONSdigital/dp-publishing-dataset-controller/mapper"
	"github.com/ONSdigital/dp-publishing-dataset-controller/model"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// PutMetadata updates all the dataset, version and dimension object fields
func PutMetadata(dc DatasetAPIClient, zc ZebedeeClient) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string) {
		putMetadata(w, r, dc, zc, accessToken, collectionID, lang)
	})
}

func putMetadata(w http.ResponseWriter, req *http.Request, dc DatasetAPIClient, zc ZebedeeClient, userAccessToken, collectionID, lang string) {
	ctx := req.Context()

	err := checkAccessTokenAndCollectionHeaders(userAccessToken, collectionID)
	if err != nil {
		log.Error(ctx, err.Error(), err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(req)
	datasetID := vars[datasetIDKey]
	edition := vars["editionID"]
	version := vars["versionID"]

	logInfo := map[string]interface{}{
		datasetIDKey: datasetID,
		editionKey:   edition,
		"version":    version,
	}

	headers := datasetApiSdk.Headers{
		CollectionID: collectionID,
		AccessToken:  userAccessToken,
	}

	b, err := io.ReadAll(req.Body)
	if err != nil {
		log.Error(ctx, "putMetadata endpoint: error reading body", err, log.Data(logInfo))
		http.Error(w, "error reading body", http.StatusBadRequest)
		return
	}

	var body model.EditMetadata
	if err = json.Unmarshal(b, &body); err != nil {
		log.Error(ctx, "putMetadata endpoint: error unmarshalling body", err, log.Data(logInfo))
		http.Error(w, "error unmarshalling body", http.StatusBadRequest)
		return
	}

	err = dc.PutDataset(ctx, headers, datasetID, body.Dataset)
	if err != nil {
		log.Error(ctx, "error updating dataset", err, log.Data(logInfo))
		http.Error(w, "error updating dataset", http.StatusInternalServerError)
		return
	}

	_, err = dc.PutVersion(ctx, headers, datasetID, edition, version, body.Version)
	if err != nil {
		log.Error(ctx, "error updating version", err, log.Data(logInfo))
		http.Error(w, "error updating version", http.StatusInternalServerError)
		return
	}

	instance := datasetApiSdk.UpdateInstance{}
	instance.InstanceID = body.Version.ID
	instance.Dimensions = body.Dimensions

	_, err = dc.PutInstance(ctx, headers, body.Version.ID, instance, "")
	if err != nil {
		log.Error(ctx, "error updating dimensions", err, log.Data(logInfo))
		http.Error(w, "error updating dimensions", http.StatusInternalServerError)
		return
	}

	err = zc.PutDatasetInCollection(ctx, userAccessToken, collectionID, "", datasetID, body.CollectionState)
	if err != nil {
		log.Error(ctx, "error adding dataset to collection", err, log.Data(logInfo))
		http.Error(w, "error adding dataset to collection", http.StatusInternalServerError)
		return
	}

	err = zc.PutDatasetVersionInCollection(ctx, userAccessToken, collectionID, "", datasetID, edition, version, body.CollectionState)
	if err != nil {
		log.Error(ctx, "error adding version to collection", err, log.Data(logInfo))
		http.Error(w, "error adding version to collection", http.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(body)
	if err != nil {
		log.Error(ctx, "error marshalling response", err, log.Data(logInfo))
		http.Error(w, "error marshalling response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseBody)
	if err != nil {
		log.Error(ctx, "error writing response", err)
		http.Error(w, "error writing response", http.StatusInternalServerError)
		return
	}

	log.Info(ctx, "put metadata: request successful", log.Data(logInfo))
}

// PutEditableMetadata updates a given list of metadata fields, agreed as being editable for both a dataset and a version object
// This new endpoint makes a unique call to the dataset api updating only the relevant metadata fields in a transactional way
// It also calls zebedee to update the collection
func PutEditableMetadata(dc DatasetAPIClient, zc ZebedeeClient) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string) {
		putEditableMetadata(w, r, dc, zc, accessToken, collectionID)
	})
}

func putEditableMetadata(w http.ResponseWriter, req *http.Request, dc DatasetAPIClient, zc ZebedeeClient, userAccessToken, collectionID string) {
	ctx := req.Context()

	err := checkAccessTokenAndCollectionHeaders(userAccessToken, collectionID)
	if err != nil {
		log.Error(ctx, err.Error(), err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(req)
	datasetID := vars[datasetIDKey]
	edition := vars["editionID"]
	version := vars["versionID"]

	logInfo := map[string]interface{}{
		datasetIDKey: datasetID,
		editionKey:   edition,
		"version":    version,
	}

	headers := datasetApiSdk.Headers{
		CollectionID: collectionID,
		AccessToken:  userAccessToken,
	}

	b, err := io.ReadAll(req.Body)
	if err != nil {
		log.Error(ctx, "putMetadata endpoint: error reading body", err, log.Data(logInfo))
		http.Error(w, "error reading body", http.StatusBadRequest)
		return
	}

	var body model.EditMetadata
	if err = json.Unmarshal(b, &body); err != nil {
		log.Error(ctx, "putMetadata endpoint: error unmarshalling body", err, log.Data(logInfo))
		http.Error(w, "error unmarshalling body", http.StatusBadRequest)
		return
	}

	versionEtag := body.VersionEtag

	editableMetadata := mapper.PutMetadata(body)

	err = dc.PutMetadata(ctx, headers, datasetID, edition, version, editableMetadata, versionEtag)
	if err != nil {
		log.Error(ctx, "error updating metadata", err, log.Data(logInfo))
		http.Error(w, "error updating metadata", http.StatusInternalServerError)
		return
	}

	err = zc.PutDatasetInCollection(ctx, userAccessToken, collectionID, "", datasetID, body.CollectionState)
	if err != nil {
		log.Error(ctx, "error adding dataset to collection", err, log.Data(logInfo))
		http.Error(w, "error adding dataset to collection", http.StatusInternalServerError)
		return
	}

	err = zc.PutDatasetVersionInCollection(ctx, userAccessToken, collectionID, "", datasetID, edition, version, body.CollectionState)
	if err != nil {
		log.Error(ctx, "error adding version to collection", err, log.Data(logInfo))
		http.Error(w, "error adding version to collection", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(b); err != nil {
		log.Error(ctx, "failed to write response body", err, log.Data(logInfo))
		http.Error(w, "failed to write response body", http.StatusInternalServerError)
		return
	}

	log.Info(ctx, "put metadata: request successful", log.Data(logInfo))
}
