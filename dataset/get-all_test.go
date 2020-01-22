package dataset

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"

	datasetclient "github.com/ONSdigital/dp-api-clients-go/dataset"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitGetAllDatasets(t *testing.T) {

	mockedDatasetResponse := []datasetclient.Dataset{
		datasetclient.Dataset{
			ID: "id-1",
			Next: &datasetclient.DatasetDetails{
				Title: "Test title 1",
			},
		},
		datasetclient.Dataset{
			ID: "id-2",
			Next: &datasetclient.DatasetDetails{
				Title: "Test title 2",
			},
		},
	}

	expectedSuccessResponse := "[{\"id\":\"id-1\",\"title\":\"Test title 1\"},{\"id\":\"id-2\",\"title\":\"Test title 2\"}]"

	Convey("test getAllDatasets", t, func() {
		Convey("on success", func() {

			mockDatasetClient := &ClientMock{
				GetDatasetsFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, collectionID string) (datasetclient.List, error) {
					return datasetclient.List{Items: mockedDatasetResponse}, nil
				},
			}

			req := httptest.NewRequest("GET", "/datasets", nil)
			req.Header.Set("Collection-Id", "testcollection")
			req.Header.Set("X-Florence-Token", "testuser")
			rec := httptest.NewRecorder()
			router := mux.NewRouter()
			router.Path("/datasets").HandlerFunc(GetAll(mockDatasetClient))

			Convey("returns 200 response", func() {
				router.ServeHTTP(rec, req)
				So(rec.Code, ShouldEqual, http.StatusOK)
			})

			Convey("returns JSON response", func() {
				router.ServeHTTP(rec, req)
				response := rec.Body.String()
				spew.Dump(response)
				So(response, ShouldEqual, expectedSuccessResponse)
			})
		})

		Convey("errors if no headers are passed", func() {

			mockDatasetClient := &ClientMock{
				GetDatasetsFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, collectionID string) (datasetclient.List, error) {
					return datasetclient.List{}, nil
				},
			}

			Convey("collection id not set", func() {
				req := httptest.NewRequest("GET", "/datasets", nil)
				req.Header.Set("X-Florence-Token", "testuser")
				rec := httptest.NewRecorder()
				router := mux.NewRouter()
				router.Path("/datasets").HandlerFunc(GetAll(mockDatasetClient))

				Convey("returns 400 response", func() {
					router.ServeHTTP(rec, req)
					So(rec.Code, ShouldEqual, http.StatusBadRequest)
				})

				Convey("returns error body", func() {
					router.ServeHTTP(rec, req)
					response := rec.Body.String()
					spew.Dump(response)
					So(response, ShouldResemble, "no collection ID header set\n")
				})
			})

			Convey("user auth token not set", func() {
				req := httptest.NewRequest("GET", "/datasets", nil)
				req.Header.Set("Collection-Id", "testcollection")
				rec := httptest.NewRecorder()
				router := mux.NewRouter()
				router.Path("/datasets").HandlerFunc(GetAll(mockDatasetClient))

				Convey("returns 400 response", func() {
					router.ServeHTTP(rec, req)
					So(rec.Code, ShouldEqual, http.StatusBadRequest)
				})

				Convey("returns error body", func() {
					router.ServeHTTP(rec, req)
					response := rec.Body.String()
					spew.Dump(response)
					So(response, ShouldResemble, "no user access token header set\n")
				})
			})
		})

		Convey("handles error from dataset client", func() {

			mockDatasetClient := &ClientMock{
				GetDatasetsFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, collectionID string) (datasetclient.List, error) {
					return datasetclient.List{}, errors.New("test dataset API error")
				},
			}

			req := httptest.NewRequest("GET", "/datasets", nil)
			req.Header.Set("Collection-Id", "testcollection")
			req.Header.Set("X-Florence-Token", "testuser")
			rec := httptest.NewRecorder()
			router := mux.NewRouter()
			router.Path("/datasets").HandlerFunc(GetAll(mockDatasetClient))

			Convey("returns 500 response", func() {
				router.ServeHTTP(rec, req)
				So(rec.Code, ShouldEqual, http.StatusInternalServerError)
			})

			Convey("returns error body", func() {
				router.ServeHTTP(rec, req)
				response := rec.Body.String()
				spew.Dump(response)
				So(response, ShouldResemble, "error getting all datasets from dataset API\n")
			})

		})
	})
}
