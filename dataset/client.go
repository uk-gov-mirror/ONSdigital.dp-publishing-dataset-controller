package dataset

import (
	"context"

	datasetclient "github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	zebedeeclient "github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	babbageclient "github.com/ONSdigital/dp-publishing-dataset-controller/clients/topics"
)

//go:generate moq -out mocks_test.go -pkg dataset . DatasetClient ZebedeeClient BabbageClient

type DatasetClient interface {
	GetDatasetsInBatches(ctx context.Context, userAuthToken, serviceAuthToken, collectionID string, batchSize, maxWorkers int) (datasetclient.List, error)
	Get(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID string) (m datasetclient.DatasetDetails, err error)
	GetVersionsInBatches(ctx context.Context, userAuthToken, serviceAuthToken, downloadServiceAuthToken, collectionID, datasetID, edition string, batchSize, maxWorkers int) (m datasetclient.VersionsList, err error)
	GetDatasetCurrentAndNext(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID string) (m datasetclient.Dataset, err error)
	GetEdition(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID, edition string) (m datasetclient.Edition, err error)
	GetEditions(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID string) (m []datasetclient.Edition, err error)
	GetVersion(ctx context.Context, userAuthToken, serviceAuthToken, downloadServiceAuthToken, collectionID, datasetID, edition, version string) (m datasetclient.Version, err error)
	GetInstance(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, instanceID, ifMatch string) (i datasetclient.Instance, eTag string, err error)
	PutDataset(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID string, d datasetclient.DatasetDetails) error
	PutVersion(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID, edition, version string, v datasetclient.Version) error
	PutInstance(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, instanceID string, i datasetclient.UpdateInstance, ifMatch string) (eTag string, err error)
	PutMetadata(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID, edition, version string, metadata datasetclient.EditableMetadata, versionEtag string) error
}

type ZebedeeClient interface {
	GetCollection(ctx context.Context, userAccessToken, collectionID string) (c zebedeeclient.Collection, err error)
	PutDatasetInCollection(ctx context.Context, userAccessToken, collectionID, lang, datasetID, state string) error
	PutDatasetVersionInCollection(ctx context.Context, userAccessToken, collectionID, lang, datasetID, edition, version, state string) error
}

type BabbageClient interface {
	GetTopics(ctx context.Context, userAccessToken string) (result babbageclient.TopicsResult, err error)
}
