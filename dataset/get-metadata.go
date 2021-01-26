package dataset

import (
	"encoding/json"
	"net/http"

	dphandlers "github.com/ONSdigital/dp-net/handlers"
	"github.com/ONSdigital/dp-publishing-dataset-controller/mapper"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

// ClientError implements error interface with additional code method
type ClientError interface {
	error
	Code() int
}

// GetEditMetadataHandler is a handler that wraps getEditMetadataHandler passing in addition arguments
func GetMetadataHandler(dc DatasetClient, zc ZebedeeClient) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string) {
		getEditMetadataHandler(w, r, dc, zc, accessToken, collectionID, lang)
	})
}

// getEditMetadataHandler gets the Edit Metadata page information used on the edit metadata screens
func getEditMetadataHandler(w http.ResponseWriter, req *http.Request, dc DatasetClient, zc ZebedeeClient, userAccessToken, collectionID, lang string) {
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

	v, err := dc.GetVersion(ctx, userAccessToken, "", "", collectionID, datasetID, edition, version)
	if err != nil {
		log.Event(ctx, "failed Get version details", log.Error(err), log.Data(logInfo))
		setErrorStatusCode(req, w, err, datasetID)
		return
	}

	d, err := dc.Get(ctx, userAccessToken, "", collectionID, datasetID)
	if err != nil {
		log.Event(ctx, "failed Get dataset details", log.Error(err), log.Data(logInfo))
		setErrorStatusCode(req, w, err, datasetID)
		return
	}

	i, err := dc.GetInstance(ctx, userAccessToken, "", collectionID, v.InstanceID)
	if err != nil {
		log.Event(ctx, "failed Get instance details", log.Error(err), log.Data(logInfo))
		setErrorStatusCode(req, w, err, datasetID)
		return
	}

	c, err := zc.GetCollection(ctx, userAccessToken, collectionID)
	if err != nil {
		log.Event(ctx, "failed Get collection details", log.Error(err), log.Data(logInfo))
		setErrorStatusCode(req, w, err, datasetID)
		return
	}

	p := mapper.EditMetadata(d, v, i, c)

	b, err := json.Marshal(p)
	if err != nil {
		log.Event(ctx, "failed marshalling page into bytes", log.Error(err))
		setErrorStatusCode(req, w, err, datasetID)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(b)
	if err != nil {
		log.Event(ctx, "failed to write bytes for http response", log.Error(err), log.Data(logInfo))
		setErrorStatusCode(req, w, err, datasetID)
		return
	}

}

func setErrorStatusCode(req *http.Request, w http.ResponseWriter, err error, datasetID string) {
	status := http.StatusInternalServerError
	if err, ok := err.(ClientError); ok {
		if err.Code() == http.StatusNotFound {
			status = err.Code()
		}
	}
	log.Event(req.Context(), "client error", log.ERROR, log.Error(err), log.Data{"setting-response-status": status, "datasetID": datasetID})
	w.WriteHeader(status)
}
