package dataset

import (
	"context"

	datasetclient "github.com/ONSdigital/dp-api-clients-go/dataset"
	zebedeeclient "github.com/ONSdigital/dp-api-clients-go/zebedee"
)

//go:generate moq -out mocks_test.go -pkg dataset . DatasetClient ZebedeeClient

type DatasetClient interface {
	//healthcheck.Client
	GetDatasets(ctx context.Context, userAuthToken, serviceAuthToken, collectionID string) (m datasetclient.List, err error)
	Get(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID string) (m datasetclient.DatasetDetails, err error)
	GetDatasetCurrentAndNext(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID string) (m datasetclient.Dataset, err error)
	GetVersion(ctx context.Context, userAuthToken, serviceAuthToken, downloadServiceAuthToken, collectionID, datasetID, edition, version string) (m datasetclient.Version, err error)
	GetInstance(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, instanceID string) (i datasetclient.Instance, err error)
	PutDataset(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID string, d datasetclient.DatasetDetails) error
	PutVersion(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID, edition, version string, v datasetclient.Version) error
	PutInstance(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, instanceID string, i datasetclient.Instance) error
}

type ZebedeeClient interface {
	GetCollection(ctx context.Context, userAccessToken, collectionID string) (c zebedeeclient.Collection, err error)
	PutDatasetInCollection(ctx context.Context, userAccessToken, collectionID, lang, datasetID, state string) error
	PutDatasetVersionInCollection(ctx context.Context, userAccessToken, collectionID, lang, datasetID, edition, version, state string) error
}
