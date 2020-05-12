package handlers

import (
	"context"
	"encoding/json"
	"github.com/ONSdigital/dp-publishing-dataset-controller/mapper"
	"github.com/ONSdigital/go-ns/common"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
)

// ClientError implements error interface with additional code method
type ClientError interface {
	error
	Code() int
}

// GetEditMetadataHandler is a handler that wraps getEditMetadataHandler passing in addition arguments
func GetEditMetadataHandler(dc DatasetClient) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		getEditMetadataHandler(w, req, dc)
	}
}

// getEditMetadataHandler gets the Edit Metadata page information used on the edit metadata screens
func getEditMetadataHandler(w http.ResponseWriter, req *http.Request, dc DatasetClient) {
	vars := mux.Vars(req)
	datasetID := vars["datasetID"]
	edition := vars["editionID"]
	version := vars["versionID"]
	ctx := req.Context()
	userAccessToken := getUserAccessTokenFromContext(ctx)
	collectionID := getCollectionIDFromContext(ctx)

	v, err := dc.GetVersion(ctx, userAccessToken, "", "", collectionID, datasetID, edition, version)
	if err != nil {
		setStatusCode(req, w, err)
		return
	}

	d, err := dc.Get(ctx, userAccessToken, "", collectionID, datasetID)
	if err != nil {
		log.Event(ctx, "failed Get dataset details", log.Error(err))
		setStatusCode(req, w, err)
		return
	}
	p, err := mapper.EditDatasetVersionMetaData(d, v)
	if err != nil {
		err := errors.Wrap(err, "failed to map EditDatasetVersionMetaData")
		log.Event(ctx, "failed to map EditDatasetVersionMetaData", log.Error(err))
		setStatusCode(req, w, err)
		return
	}

	b, err := json.Marshal(p)
	if err != nil {

		setStatusCode(req, w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(b)
	if err != nil {

		log.Event(ctx, "failed to write bytes for http response", log.Error(err))
		setStatusCode(req, w, err)
		return
	}

}

func setStatusCode(req *http.Request, w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if err, ok := err.(ClientError); ok {
		if err.Code() == http.StatusNotFound {
			status = err.Code()
		}
	}
	log.Event(req.Context(), "client error", log.ERROR, log.Error(err), log.Data{"setting-response-status": status})
	w.WriteHeader(status)
}

func getUserAccessTokenFromContext(ctx context.Context) string {
	if ctx.Value(common.FlorenceIdentityKey) != nil {
		accessToken, ok := ctx.Value(common.FlorenceIdentityKey).(string)
		if !ok {
			log.Event(ctx, "error retrieving user access token", log.WARN, log.Error(errors.New("error casting access token context value to string")))
		}
		return accessToken
	}
	return ""
}

func getCollectionIDFromContext(ctx context.Context) string {
	if ctx.Value(common.CollectionIDHeaderKey) != nil {
		collectionID, ok := ctx.Value(common.CollectionIDHeaderKey).(string)
		if !ok {
			log.Event(ctx, "error retrieving collection ID", log.WARN, log.Error(errors.New("error casting collection ID context value to string")))
		}
		return collectionID
	}
	return ""
}
