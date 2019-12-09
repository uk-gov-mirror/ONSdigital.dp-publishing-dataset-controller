package handler

import (
	"net/http"
	"encoding/json"

	"github.com/ONSdigital/dp-api-clients-go/dataset"
	"github.com/ONSdigital/log.go/log"
)

// GetAllDatasets returns a mapped list of all datasets
func GetAllDatasets(dc *dataset.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		getAllDatasets(w, req, dc)
	}
}

func getAllDatasets(w http.ResponseWriter, req *http.Request, dc *dataset.Client) {
	ctx := req.Context()
	// userAccessToken := getUserAccessTokenFromContent(ctx)
	// collectionID := getCollectionIDFromContext(ctx)
	userAccessToken := "39db1923182338dc096c662eae572a5d8debff2869e619908d99ac2bbd07aa4c"
	collectionID := "00test-906d2af2c0c50a81b8bc944327971d852ed6b7f89f5a1841911854f59478033f"

	datasets, err := dc.GetDatasets(ctx, userAccessToken, "", collectionID)
	if err != nil {
		log.Event(nil, "error getting all datasets", log.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	b, _ := json.Marshal(datasets)
	w.Write(b)
}
