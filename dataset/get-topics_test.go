package dataset

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	babbageclient "github.com/ONSdigital/dp-publishing-dataset-controller/clients/topics"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitGetAllTopics(t *testing.T) {

	mockTopics := babbageclient.TopicsResult{
		Topics: babbageclient.Topic{
			Results: []babbageclient.Result{{
				Description: babbageclient.Description{
					Title: "test 1",
				},
				URI:  "/test/uri/1",
				Type: "page",
			},
				{
					Description: babbageclient.Description{
						Title: "test 2",
					},
					URI:  "/test/uri/2",
					Type: "page",
				}},
		},
	}

	expectedSuccessResponse := "[{\"title\":\"test 1\"},{\"title\":\"test 2\"}]"

	Convey("test getTopics", t, func() {
		Convey("on success", func() {

			mockBabbageClient := &BabbageClientMock{
				GetTopicsFunc: func(ctx context.Context, userAuthToken string) (babbageclient.TopicsResult, error) {
					return mockTopics, nil
				},
			}

			req := httptest.NewRequest("GET", "/datasets/123/create", nil)
			req.Header.Set("Collection-Id", "testcollection")
			req.Header.Set("X-Florence-Token", "testuser")
			rec := httptest.NewRecorder()
			router := mux.NewRouter()
			router.Path("/datasets/123/create").HandlerFunc(GetTopics(mockBabbageClient))
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

			mockBabbageClient := &BabbageClientMock{
				GetTopicsFunc: func(ctx context.Context, userAuthToken string) (babbageclient.TopicsResult, error) {
					return babbageclient.TopicsResult{}, nil
				},
			}

			Convey("collection id not set", func() {
				req := httptest.NewRequest("GET", "/datasets/123/create", nil)
				req.Header.Set("X-Florence-Token", "testuser")
				rec := httptest.NewRecorder()
				router := mux.NewRouter()
				router.Path("/datasets/123/create").HandlerFunc(GetTopics(mockBabbageClient))

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
				req := httptest.NewRequest("GET", "/datasets/123/create", nil)
				req.Header.Set("Collection-Id", "testcollection")
				rec := httptest.NewRecorder()
				router := mux.NewRouter()
				router.Path("/datasets/123/create").HandlerFunc(GetTopics(mockBabbageClient))

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

		Convey("handles error from babbage client", func() {

			mockBabbageClient := &BabbageClientMock{
				GetTopicsFunc: func(ctx context.Context, userAuthToken string) (babbageclient.TopicsResult, error) {
					return babbageclient.TopicsResult{}, errors.New("test babbage API error")
				},
			}

			req := httptest.NewRequest("GET", "/datasets/123/create", nil)
			req.Header.Set("Collection-Id", "testcollection")
			req.Header.Set("X-Florence-Token", "testuser")
			rec := httptest.NewRecorder()
			router := mux.NewRouter()
			router.Path("/datasets/123/create").HandlerFunc(GetTopics(mockBabbageClient))

			Convey("returns 500 response", func() {
				router.ServeHTTP(rec, req)
				So(rec.Code, ShouldEqual, http.StatusInternalServerError)
			})

			Convey("returns error body", func() {
				router.ServeHTTP(rec, req)
				response := rec.Body.String()
				So(response, ShouldResemble, "error getting topics\n")
			})

		})
	})
}
