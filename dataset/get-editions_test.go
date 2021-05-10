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

func TestUnitGetEditions(t *testing.T) {
	t.Parallel()

	datasetID := "test-dataset"

	mockedDatasetResponse := datasetclient.Dataset{
		Next: &datasetclient.DatasetDetails{
			Title: "Test title",
		},
	}

	mockedEditionResponse := []datasetclient.Edition{
		{
			Edition: "edition-1",
		},
		{
			Edition: "edition-2",
		},
	}

	mockedVersionResponse := datasetclient.Version{
		ID:          "version-1",
		InstanceID:  "instance-001",
		Version:     1,
		ReleaseDate: "2020-11-07T00:00:00.000Z",
	}

	expectedSuccessResponse := "{\"dataset_name\":\"Test title\",\"editions\":[{\"id\":\"edition-1\",\"title\":\"edition-1\",\"release_date\":\"\"},{\"id\":\"edition-2\",\"title\":\"edition-2\",\"release_date\":\"\"}]}"

	Convey("test getAllEditions", t, func() {

		mockDatasetClient := &DatasetClientMock{
			GetDatasetCurrentAndNextFunc: func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, datasetID string) (datasetclient.Dataset, error) {
				return mockedDatasetResponse, nil
			},
			GetEditionsFunc: func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, datasetID string) ([]datasetclient.Edition, error) {
				return mockedEditionResponse, nil
			},
			GetVersionFunc: func(ctx context.Context, userAuthToken string, serviceAuthToken string, downloadServiceAuthToken string, collectionID string, datasetID string, editionID string, versionID string) (datasetclient.Version, error) {
				return mockedVersionResponse, nil
			},
		}

		Convey("on success", func() {
			reqURL := fmt.Sprintf("/datasets/%v/editions", datasetID)
			req := httptest.NewRequest("GET", reqURL, nil)
			req.Header.Set("Collection-Id", "testcollection")
			req.Header.Set("X-Florence-Token", "testuser")
			rec := httptest.NewRecorder()
			router := mux.NewRouter()
			router.Path(reqURL).HandlerFunc(GetEditions(mockDatasetClient))

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
				reqURL := fmt.Sprintf("/datasets/%v/editions", datasetID)
				req := httptest.NewRequest("GET", reqURL, nil)
				req.Header.Set("X-Florence-Token", "testuser")
				rec := httptest.NewRecorder()
				router := mux.NewRouter()
				router.Path(reqURL).HandlerFunc(GetEditions(mockDatasetClient))

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
				reqURL := fmt.Sprintf("/datasets/%v/editions", datasetID)
				req := httptest.NewRequest("GET", reqURL, nil)
				req.Header.Set("Collection-Id", "testcollection")
				rec := httptest.NewRecorder()
				router := mux.NewRouter()
				router.Path(reqURL).HandlerFunc(GetEditions(mockDatasetClient))

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
				GetEditionsFunc: func(ctx context.Context, userAuthToken string, serviceAuthToken string, collectionID string, datasetID string) ([]datasetclient.Edition, error) {
					return mockedEditionResponse, errors.New("test dataset API error")
				},
				GetVersionFunc: func(ctx context.Context, userAuthToken string, serviceAuthToken string, downloadServiceAuthToken string, collectionID string, datasetID string, editionID string, versionID string) (datasetclient.Version, error) {
					return mockedVersionResponse, nil
				},
			}

			reqURL := fmt.Sprintf("/datasets/%v/editions", datasetID)
			req := httptest.NewRequest("GET", reqURL, nil)
			req.Header.Set("Collection-Id", "testcollection")
			req.Header.Set("X-Florence-Token", "testuser")
			rec := httptest.NewRecorder()
			router := mux.NewRouter()
			router.Path(reqURL).HandlerFunc(GetEditions(mockDatasetClient))

			Convey("returns 500 response", func() {
				router.ServeHTTP(rec, req)
				So(rec.Code, ShouldEqual, http.StatusInternalServerError)
			})

			Convey("returns error body", func() {
				router.ServeHTTP(rec, req)
				response := rec.Body.String()
				So(response, ShouldResemble, "error getting editions from dataset API: test dataset API error\n")
			})

		})
	})
}
