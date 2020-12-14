package dataset

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	datasetclient "github.com/ONSdigital/dp-api-clients-go/dataset"
	"github.com/ONSdigital/dp-publishing-dataset-controller/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitPutMetadata(t *testing.T) {

	putBody, _ := json.Marshal(model.EditMetadata{
		Dataset: datasetclient.DatasetDetails{
			ID: "test-dataset",
		},
		Version: datasetclient.Version{
			ID: "1",
		},
		Instance:        datasetclient.Instance{},
		CollectionID:    "testcollection",
		CollectionState: "InProgress",
	})

	mockedPutBody := bytes.NewReader(putBody)

	Convey("test putMetadata", t, func() {
		Convey("on success", func() {

			mockDatasetClient := &DatasetClientMock{
				PutDatasetFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID string, d datasetclient.DatasetDetails) error {
					return nil
				},
				PutVersionFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID, edition, version string, v datasetclient.Version) error {
					return nil
				},
				PutInstanceFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, instanceID string, i datasetclient.Instance) error {
					return nil
				},
			}

			mockZebedeeClient := &ZebedeeClientMock{
				PutDatasetInCollectionFunc: func(ctx context.Context, userAccessToken, collectionID, lang, datasetID, state string) error {
					return nil
				},
				PutDatasetVersionInCollectionFunc: func(ctx context.Context, userAccessToken, collectionID, lang, datasetID, edition, version, state string) error {
					return nil
				},
			}

			req := httptest.NewRequest("PUT", "/datasets/test-dataset/editions/test-edition/versions/1", mockedPutBody)
			req.Header.Set("Collection-Id", "testcollection")
			req.Header.Set("X-Florence-Token", "testuser")
			rec := httptest.NewRecorder()
			router := mux.NewRouter()
			router.Path("/datasets/{datasetID}/editions/{editionID}/versions/{versionID}").HandlerFunc(PutMetadata(mockDatasetClient, mockZebedeeClient))

			Convey("returns 200 response", func() {
				router.ServeHTTP(rec, req)
				So(rec.Code, ShouldEqual, http.StatusOK)
			})
		})

		Convey("errors if no headers are passed", func() {

			mockDatasetClient := &DatasetClientMock{
				PutDatasetFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID string, d datasetclient.DatasetDetails) error {
					return nil
				},
				PutVersionFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID, edition, version string, v datasetclient.Version) error {
					return nil
				},
				PutInstanceFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, instanceID string, i datasetclient.Instance) error {
					return nil
				},
			}

			mockZebedeeClient := &ZebedeeClientMock{
				PutDatasetInCollectionFunc: func(ctx context.Context, userAccessToken, collectionID, lang, datasetID, state string) error {
					return nil
				},
				PutDatasetVersionInCollectionFunc: func(ctx context.Context, userAccessToken, collectionID, lang, datasetID, edition, version, state string) error {
					return nil
				},
			}

			Convey("collection id not set", func() {
				req := httptest.NewRequest("PUT", "/datasets/test-dataset/editions/test-edition/versions/1", mockedPutBody)
				req.Header.Set("X-Florence-Token", "testuser")
				rec := httptest.NewRecorder()
				router := mux.NewRouter()
				router.Path("/datasets/{datasetID}/editions/{editionID}/versions/{versionID}").HandlerFunc(PutMetadata(mockDatasetClient, mockZebedeeClient))

				Convey("returns 400 response", func() {
					router.ServeHTTP(rec, req)
					So(rec.Code, ShouldEqual, http.StatusBadRequest)
				})

				Convey("returns error body", func() {
					router.ServeHTTP(rec, req)
					response := rec.Body.String()
					So(response, ShouldResemble, "no collection ID header set\n")
				})
			})

			Convey("user auth token not set", func() {
				req := httptest.NewRequest("PUT", "/datasets/test-dataset/editions/test-edition/versions/1", mockedPutBody)
				req.Header.Set("Collection-Id", "testcollection")
				rec := httptest.NewRecorder()
				router := mux.NewRouter()
				router.Path("/datasets/{datasetID}/editions/{editionID}/versions/{versionID}").HandlerFunc(PutMetadata(mockDatasetClient, mockZebedeeClient))

				Convey("returns 400 response", func() {
					router.ServeHTTP(rec, req)
					So(rec.Code, ShouldEqual, http.StatusBadRequest)
				})

				Convey("returns error body", func() {
					router.ServeHTTP(rec, req)
					response := rec.Body.String()
					So(response, ShouldResemble, "no user access token header set\n")
				})
			})
		})

		Convey("handles error from dataset client", func() {

			mockDatasetClient := &DatasetClientMock{
				PutDatasetFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID string, d datasetclient.DatasetDetails) error {
					return nil
				},
				PutVersionFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID, edition, version string, v datasetclient.Version) error {
					return errors.New("test dataset API error")
				},
				PutInstanceFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, instanceID string, i datasetclient.Instance) error {
					return nil
				},
			}

			mockZebedeeClient := &ZebedeeClientMock{
				PutDatasetInCollectionFunc: func(ctx context.Context, userAccessToken, collectionID, lang, datasetID, state string) error {
					return nil
				},
				PutDatasetVersionInCollectionFunc: func(ctx context.Context, userAccessToken, collectionID, lang, datasetID, edition, version, state string) error {
					return nil
				},
			}

			req := httptest.NewRequest("PUT", "/datasets/test-dataset/editions/test-edition/versions/1", mockedPutBody)
			req.Header.Set("Collection-Id", "testcollection")
			req.Header.Set("X-Florence-Token", "testuser")
			rec := httptest.NewRecorder()
			router := mux.NewRouter()
			router.Path("/datasets/{datasetID}/editions/{editionID}/versions/{versionID}").HandlerFunc(PutMetadata(mockDatasetClient, mockZebedeeClient))

			Convey("returns 500 response and error body", func() {
				router.ServeHTTP(rec, req)
				So(rec.Code, ShouldEqual, http.StatusInternalServerError)
				response := rec.Body.String()
				So(response, ShouldResemble, "error updating version\n")
			})

		})
	})
}
