package dataset

import (
	"encoding/json"
	"net/http"

	dphandlers "github.com/ONSdigital/dp-net/handlers"
	"github.com/ONSdigital/dp-publishing-dataset-controller/mapper"
	"github.com/ONSdigital/log.go/v2/log"
)

// GetTopics returns a mapped list of topics
func GetTopics(bc BabbageClient) http.HandlerFunc {
	return dphandlers.ControllerHandler(func(w http.ResponseWriter, r *http.Request, lang, collectionID, accessToken string) {
		getTopics(w, r, bc, accessToken, collectionID, lang)
	})
}

func getTopics(w http.ResponseWriter, req *http.Request, bc BabbageClient, userAccessToken, collectionID, lang string) {
	ctx := req.Context()

	err := checkAccessTokenAndCollectionHeaders(userAccessToken, collectionID)
	if err != nil {
		log.Error(ctx, err.Error(), err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Info(ctx, "calling get topics")

	topics, err := bc.GetTopics(ctx, userAccessToken)
	if err != nil {
		log.Error(ctx, "error getting topics", err)
		http.Error(w, "error getting topics", http.StatusInternalServerError)
		return
	}

	mapped := mapper.Topics(topics)

	b, err := json.Marshal(mapped)
	if err != nil {
		log.Error(ctx, "error marshalling response to json", err)
		http.Error(w, "error marshalling response to json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)

	log.Info(ctx, "get topics: request successful")
}
