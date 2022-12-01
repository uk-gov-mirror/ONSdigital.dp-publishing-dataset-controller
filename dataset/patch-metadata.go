package dataset

import (
	"fmt"
	"net/http"

	dphandlers "github.com/ONSdigital/dp-net/handlers"
	dprequest "github.com/ONSdigital/dp-net/v2/request"
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

	//unmarshal request body to create patch array
	patches, err := dprequest.GetPatches(req.Body, []dprequest.PatchOp{dprequest.OpReplace})
	if err != nil {
		log.Error(ctx, "error obtaining patch from request body", err, log.Data(logInfo))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// create array of patch operations for both the dataset and version dataset update
	datasetPatchBody := []dprequest.Patch{}
	datasetVersionPatchBody := []dprequest.Patch{}

	for _, patch := range patches {
		switch patch.Path {
		case "/release_frequency", "/license", "/survey", "/next_release", "/canonical_topic", "/related_datasets", "/publications", "/methodologies", "/related_content", "/qmi", "/subtopics", "/national_statistic":
			datasetPatchBody = append(datasetPatchBody, patch)
		case "/alerts", "/latest_changes", "/release_date", "/usage_notes":
			datasetVersionPatchBody = append(datasetVersionPatchBody, patch)
		default:
			log.Error(ctx, fmt.Sprintf("invalid patch path: %s", patch.Path), err, log.Data(logInfo))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	}

	if len(datasetPatchBody) > 0 {
		err = dc.PatchDataset(ctx, userAccessToken, "", collectionID, datasetID, ifMatch, datasetPatchBody)
		if err != nil {
			log.Error(ctx, "error updating (patch) dataset", err, log.Data(logInfo))
			http.Error(w, "error updating (patch) dataset", http.StatusInternalServerError)
			return
		}
	}
	if len(datasetVersionPatchBody) > 0 {
		err = dc.PatchVersion(ctx, userAccessToken, "", collectionID, datasetID, edition, version, ifMatch, datasetVersionPatchBody)
		if err != nil {
			log.Error(ctx, "error updating (patch) dataset version", err, log.Data(logInfo))
			http.Error(w, "error updating (patch) dataset version", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	log.Info(ctx, "patch metadata: request successful", log.Data(logInfo))
}
