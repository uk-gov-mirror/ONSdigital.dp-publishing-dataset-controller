package dataset

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	datasetclient "github.com/ONSdigital/dp-api-clients-go/dataset"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitGetVersions(t *testing.T) {
	t.Parallel()

	datasetID := "test-dataset"
	editionID := "test-edition"
	verionsBatchSize := 10
	versionsMaxWorkers := 3

	mockedDatasetResponse := datasetclient.Dataset{
		Next: &datasetclient.DatasetDetails{
			Title: "Test title",
		},
	}

	mockedEditionResponse := datasetclient.Edition{
		Edition: "edition-1",
	}

	mockedVersionsResponse := []datasetclient.Version{
		{
			ID:         "version-1",
			InstanceID: "instance-001",
			Version:    1,
		},
		{
			ID:         "version-2",
			InstanceID: "instance-002",
			Version:    2,
		},
	}

	expectedSuccessResponse := "{\"dataset_name\":\"Test title\",\"edition_name\":\"edition-1\",\"versions\":[{\"id\":\"version-2\",\"title\":\"Version: 2\",\"version\":2,\"release_date\":\"\",\"state\":\"\"},{\"id\":\"version-1\",\"title\":\"Version: 1\",\"version\":1,\"release_date\":\"\",\"state\":\"\"}]}"

	Convey("test getAllVersions", t, func() {

		mockDatasetClient := &DatasetClientMock{
			GetDatasetCurrentAndNextFunc: func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, datasetID string) (datasetclient.Dataset, error) {
				return mockedDatasetResponse, nil
			},
			GetEditionFunc: func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, datasetID string, editionID string) (datasetclient.Edition, error) {
				return mockedEditionResponse, nil
			},
			GetVersionsInBatchesFunc: func(ctx context.Context, userAuthToken string, serviceAuthToken string, downloadServiceAuthToken string, collectionID string, datasetID string, editionID string, batchSize int, maxWorkers int) (datasetclient.VersionsList, error) {
				return datasetclient.VersionsList{Items: mockedVersionsResponse}, nil
			},
		}

		Convey("on success", func() {
			reqURL := fmt.Sprintf("/datasets/%v/editions/%v/versions", datasetID, editionID)
			req := httptest.NewRequest("GET", reqURL, nil)
			req.Header.Set("Collection-Id", "testcollection")
			req.Header.Set("X-Florence-Token", "testuser")
			rec := httptest.NewRecorder()
			router := mux.NewRouter()
			router.Path(reqURL).HandlerFunc(GetVersions(mockDatasetClient, verionsBatchSize, versionsMaxWorkers))

			Convey("returns 200 response", func() {
				router.ServeHTTP(rec, req)
				So(rec.Code, ShouldEqual, http.StatusOK)
			})

			Convey("returns JSON response", func() {
				router.ServeHTTP(rec, req)
				response := rec.Body.String()
				So(response, ShouldEqual, expectedSuccessResponse)
			})
		})

		Convey("errors if no headers are passed", func() {
			Convey("collection id not set", func() {
				reqURL := fmt.Sprintf("/datasets/%v/editions/%v/versions", datasetID, editionID)
				req := httptest.NewRequest("GET", reqURL, nil)
				req.Header.Set("X-Florence-Token", "testuser")
				rec := httptest.NewRecorder()
				router := mux.NewRouter()
				router.Path(reqURL).HandlerFunc(GetVersions(mockDatasetClient, verionsBatchSize, versionsMaxWorkers))

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
				reqURL := fmt.Sprintf("/datasets/%v/editions/%v/versions", datasetID, editionID)
				req := httptest.NewRequest("GET", reqURL, nil)
				req.Header.Set("Collection-Id", "testcollection")
				rec := httptest.NewRecorder()
				router := mux.NewRouter()
				router.Path(reqURL).HandlerFunc(GetVersions(mockDatasetClient, verionsBatchSize, versionsMaxWorkers))

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
				GetDatasetCurrentAndNextFunc: func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, datasetID string) (datasetclient.Dataset, error) {
					return mockedDatasetResponse, nil
				},
				GetEditionFunc: func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, datasetID string, editionID string) (datasetclient.Edition, error) {
					return mockedEditionResponse, nil
				},
				GetVersionsInBatchesFunc: func(ctx context.Context, userAuthToken string, serviceAuthToken string, downloadServiceAuthToken string, collectionID string, datasetID string, editionID string, batchSize int, maxWorkers int) (datasetclient.VersionsList, error) {
					return datasetclient.VersionsList{}, errors.New("test dataset API error")
				},
			}

			reqURL := fmt.Sprintf("/datasets/%v/editions/%v/versions", datasetID, editionID)
			req := httptest.NewRequest("GET", reqURL, nil)
			req.Header.Set("Collection-Id", "testcollection")
			req.Header.Set("X-Florence-Token", "testuser")
			rec := httptest.NewRecorder()
			router := mux.NewRouter()
			router.Path(reqURL).HandlerFunc(GetVersions(mockDatasetClient, verionsBatchSize, versionsMaxWorkers))

			Convey("returns 500 response", func() {
				router.ServeHTTP(rec, req)
				So(rec.Code, ShouldEqual, http.StatusInternalServerError)
			})

			Convey("returns error body", func() {
				router.ServeHTTP(rec, req)
				response := rec.Body.String()
				So(response, ShouldResemble, "error getting all versions from dataset API: test dataset API error\n")
			})

		})
	})
}
