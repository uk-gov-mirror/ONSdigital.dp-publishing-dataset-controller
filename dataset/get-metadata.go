package dataset

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	datasetclient "github.com/ONSdigital/dp-api-clients-go/dataset"
	zebedeeclient "github.com/ONSdigital/dp-api-clients-go/zebedee"
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

const editionConfirmedState = "edition-confirmed"

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

	// we get the next and current doc so that we have info relating to latest published version
	// on the current doc
	d, err := dc.GetDatasetCurrentAndNext(ctx, userAccessToken, "", collectionID, datasetID)
	if err != nil {
		log.Event(ctx, "failed Get dataset details", log.Error(err), log.Data(logInfo))
		setErrorStatusCode(req, w, err, datasetID)
		return
	}

	// if the version state is "edition-confirmed" it's in a pre-edited state so we get previously
	// published version's dimensions and return those so that they are pre-populated in the browser
	// to prevent the user having to fill these in again
	dims := []datasetclient.VersionDimension{}
	if v.State == editionConfirmedState && v.Version > 1 {
		dimensions := getLatestPublishedVersionDimensions(ctx, w, req, dc, userAccessToken, collectionID, d.Current.Links.LatestVersion.URL)
		dims = append(dims, dimensions...)
	}

	c, err := getCollectionDetails(ctx, zc, userAccessToken, d.Next.CollectionID)
	if err != nil {
		log.Event(ctx, "failed Get collection details", log.Error(err), log.Data(logInfo))
		setErrorStatusCode(req, w, err, datasetID)
		return
	}

	p := mapper.EditMetadata(d.Next, v, dims, c)

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

func getCollectionDetails(ctx context.Context, zc ZebedeeClient, userAccessToken, collectionID string) (zebedeeclient.Collection, error) {
	if len(collectionID) > 0 {
		c, err := zc.GetCollection(ctx, userAccessToken, collectionID)
		if err != nil {
			return zebedeeclient.Collection{}, err
		}
		return c, nil
	} else {
		return zebedeeclient.Collection{}, nil
	}
}

func getLatestPublishedVersionDimensions(ctx context.Context, w http.ResponseWriter, req *http.Request, dc DatasetClient, userAccessToken, collectionID, latestVersionURL string) []datasetclient.VersionDimension {
	datasetID, editionID, versionID, err := getIDsFromURL(latestVersionURL)
	if err != nil {
		log.Event(ctx, "failed to parse latest version url", log.Error(err))
		return []datasetclient.VersionDimension{}
	}

	latestPublishedVersion, err := dc.GetVersion(ctx, userAccessToken, "", "", collectionID, datasetID, editionID, versionID)
	if err != nil {
		log.Event(ctx, "failed Get latest published version details", log.Error(err))
		setErrorStatusCode(req, w, err, datasetID)
		return []datasetclient.VersionDimension{}
	}

	return latestPublishedVersion.Dimensions
}

func getIDsFromURL(URL string) (datasetID, editionID, versionID string, err error) {
	parsedURL, err := url.Parse(URL)
	if err != nil {
		return "", "", "", err
	}

	s := strings.Split(parsedURL.Path, "/")
	if len(s) < 8 {
		return "", "", "", errors.New("not enough arguements in path")
	}
	datasetID = s[3]
	editionID = s[5]
	versionID = s[7]
	return datasetID, editionID, versionID, nil
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
