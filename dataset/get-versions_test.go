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

	datasetID := "test-dataset"
	editionID := "test-edition"
	verionsBatchSize := 10
	versionsMaxWorkers := 3

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

	expectedSuccessResponse := "[{\"id\":\"version-2\",\"title\":\"Version: 2\",\"version\":2,\"release_date\":\"\"},{\"id\":\"version-1\",\"title\":\"Version: 1\",\"version\":1,\"release_date\":\"\"}]"

	Convey("test getAllDatasets", t, func() {
		Convey("on success", func() {

			mockDatasetClient := &DatasetClientMock{
				GetVersionsInBatchesFunc: func(ctx context.Context, userAuthToken string, serviceAuthToken string, downloadServiceAuthToken string, collectionID string, datasetID string, editionID string, batchSize int, maxWorkers int) (datasetclient.VersionsList, error) {
					return datasetclient.VersionsList{Items: mockedVersionsResponse}, nil
				},
			}

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

			mockDatasetClient := &DatasetClientMock{
				GetVersionsInBatchesFunc: func(ctx context.Context, userAuthToken string, serviceAuthToken string, downloadServiceAuthToken string, collectionID string, datasetID string, editionID string, batchSize int, maxWorkers int) (datasetclient.VersionsList, error) {
					return datasetclient.VersionsList{Items: mockedVersionsResponse}, nil
				},
			}

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
				So(response, ShouldResemble, "error getting all datasets from dataset API: test dataset API error\n")
			})

		})
	})
}
