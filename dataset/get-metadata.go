package dataset

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	zebedeeclient "github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	datasetApiModels "github.com/ONSdigital/dp-dataset-api/models"
	datasetApiSdk "github.com/ONSdigital/dp-dataset-api/sdk"
	dphandlers "github.com/ONSdigital/dp-net/v3/handlers"
	"github.com/ONSdigital/dp-publishing-dataset-controller/mapper"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// ClientError implements error interface with additional code method
type ClientError interface {
	error
	Code() int
}

const editionConfirmedState = "edition-confirmed"

// GetEditMetadataHandler is a handler that wraps getEditMetadataHandler passing in addition arguments
func GetMetadataHandler(dc DatasetAPIClient, zc ZebedeeClient) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string) {
		getEditMetadataHandler(w, r, dc, zc, accessToken, collectionID, lang)
	})
}

// getEditMetadataHandler gets the Edit Metadata page information used on the edit metadata screens
func getEditMetadataHandler(w http.ResponseWriter, req *http.Request, dc DatasetAPIClient, zc ZebedeeClient, userAccessToken, collectionID, lang string) {
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

	v, sdkheaders, err := dc.GetVersionWithHeaders(ctx, headers, datasetID, edition, version)
	if err != nil {
		log.Error(ctx, "failed Get version details", err, log.Data(logInfo))
		setErrorStatusCode(req, w, err, datasetID)
		return
	}

	// we get the next and current doc so that we have info relating to latest published version
	// on the current doc
	d, err := dc.GetDatasetCurrentAndNext(ctx, headers, datasetID)
	if err != nil {
		log.Error(ctx, "failed Get dataset details", err, log.Data(logInfo))
		setErrorStatusCode(req, w, err, datasetID)
		return
	}

	// if the version state is "edition-confirmed" it's in a pre-edited state so we get previously
	// published version's dimensions and return those so that they are pre-populated in the browser
	// to prevent the user having to fill these in again
	dims := []datasetApiModels.Dimension{}
	if v.State == editionConfirmedState && v.Version > 1 {
		dimensions := getLatestPublishedVersionDimensions(ctx, w, req, dc, headers, d.Current.Links.LatestVersion.HRef)
		dims = append(dims, dimensions...)
	}

	c, err := getCollectionDetails(ctx, zc, userAccessToken, d.Next.CollectionID)
	if err != nil {
		log.Error(ctx, "failed Get collection details", err, log.Data(logInfo))
		setErrorStatusCode(req, w, err, datasetID)
		return
	}

	editMetadata := mapper.EditMetadata(d.Next, v, dims, c)
	editMetadata.VersionEtag = sdkheaders.ETag

	b, err := json.Marshal(editMetadata)
	if err != nil {
		log.Error(ctx, "failed marshalling page into bytes", err)
		setErrorStatusCode(req, w, err, datasetID)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(b)
	if err != nil {
		log.Error(ctx, "failed to write bytes for http response", err, log.Data(logInfo))
		setErrorStatusCode(req, w, err, datasetID)
		return
	}
}

func getCollectionDetails(ctx context.Context, zc ZebedeeClient, userAccessToken, collectionID string) (zebedeeclient.Collection, error) {
	if collectionID != "" {
		c, err := zc.GetCollection(ctx, userAccessToken, collectionID)
		if err != nil {
			return zebedeeclient.Collection{}, err
		}
		return c, nil
	} else {
		return zebedeeclient.Collection{}, nil
	}
}

func getLatestPublishedVersionDimensions(ctx context.Context, w http.ResponseWriter, req *http.Request, dc DatasetAPIClient, headers datasetApiSdk.Headers, latestVersionURL string) []datasetApiModels.Dimension {
	datasetID, editionID, versionID, err := getIDsFromURL(latestVersionURL)
	if err != nil {
		log.Error(ctx, "failed to parse latest version url", err)
		return []datasetApiModels.Dimension{}
	}

	latestPublishedVersion, err := dc.GetVersion(ctx, headers, datasetID, editionID, versionID)
	if err != nil {
		log.Error(ctx, "failed Get latest published version details", err)
		setErrorStatusCode(req, w, err, datasetID)
		return []datasetApiModels.Dimension{}
	}

	return latestPublishedVersion.Dimensions
}

func getIDsFromURL(urlValue string) (datasetID, editionID, versionID string, err error) {
	parsedURL, err := url.Parse(urlValue)
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
	log.Error(req.Context(), "client error", err, log.Data{"setting-response-status": status, datasetIDKey: datasetID})
	w.WriteHeader(status)
}
