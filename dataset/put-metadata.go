package dataset

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	datasetclient "github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	dphandlers "github.com/ONSdigital/dp-net/handlers"
	"github.com/ONSdigital/dp-publishing-dataset-controller/mapper"
	"github.com/ONSdigital/dp-publishing-dataset-controller/model"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// PutMetadata updates all the dataset, version and dimension object fields
func PutMetadata(dc DatasetClient, zc ZebedeeClient) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string) {
		putMetadata(w, r, dc, zc, accessToken, collectionID, lang)
	})
}

func putMetadata(w http.ResponseWriter, req *http.Request, dc DatasetClient, zc ZebedeeClient, userAccessToken, collectionID, lang string) {
	ctx := req.Context()

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

	err = dc.PutDataset(ctx, userAccessToken, "", collectionID, datasetID, body.Dataset)
	if err != nil {
		log.Error(ctx, "error updating dataset", err, log.Data(logInfo))
		http.Error(w, "error updating dataset", http.StatusInternalServerError)
		return
	}

	err = dc.PutVersion(ctx, userAccessToken, "", collectionID, datasetID, edition, version, body.Version)
	if err != nil {
		log.Error(ctx, "error updating version", err, log.Data(logInfo))
		http.Error(w, "error updating version", http.StatusInternalServerError)
		return
	}

	instance := datasetclient.UpdateInstance{}
	instance.InstanceID = body.Version.ID
	instance.Dimensions = body.Dimensions

	_, err = dc.PutInstance(ctx, userAccessToken, "", collectionID, body.Version.ID, instance, "")
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
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)

	log.Info(ctx, "put metadata: request successful", log.Data(logInfo))
}

// PutEditableMetadata updates a given list of metadata fields, agreed as being editable for both a dataset and a version object
// This new endpoint makes an unique call to the dataset api updating only the relevant metadata fields in a transactional way
// It also calls zebedee to update the collection
func PutEditableMetadata(dc DatasetClient, zc ZebedeeClient) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string) {
		putEditableMetadata(w, r, dc, zc, accessToken, collectionID, lang)
	})
}

func putEditableMetadata(w http.ResponseWriter, req *http.Request, dc DatasetClient, zc ZebedeeClient, userAccessToken, collectionID, lang string) {
	ctx := req.Context()

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

	err = dc.PutMetadata(ctx, userAccessToken, "", collectionID, datasetID, edition, version, editableMetadata, versionEtag)
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
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)

	log.Info(ctx, "put metadata: request successful", log.Data(logInfo))
}
