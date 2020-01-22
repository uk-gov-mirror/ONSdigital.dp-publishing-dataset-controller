package dataset

import (
	"context"

	datasetclient "github.com/ONSdigital/dp-api-clients-go/dataset"
)

//go:generate moq -out mocks_test.go -pkg dataset . Client

type Client interface {
	//healthcheck.Client
	GetDatasets(ctx context.Context, userAuthToken, serviceAuthToken, collectionID string) (m datasetclient.List, err error)
}
