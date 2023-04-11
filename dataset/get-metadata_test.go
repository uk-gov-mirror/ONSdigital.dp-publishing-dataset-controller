package dataset

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	datasetclient "github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	zebedeeclient "github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	"github.com/ONSdigital/dp-publishing-dataset-controller/model"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

type testCliError struct{}

func (e *testCliError) Error() string { return "client error" }
func (e *testCliError) Code() int     { return http.StatusNotFound }

// doTestRequest helper function that creates a router and mocks requests
func doTestRequest(target string, req *http.Request, handlerFunc http.HandlerFunc, w *httptest.ResponseRecorder) *httptest.ResponseRecorder {
	if w == nil {
		w = httptest.NewRecorder()
	}
	router := mux.NewRouter()
	router.HandleFunc(target, handlerFunc)
	router.ServeHTTP(w, req)
	return w
}

func TestUnitHandlers(t *testing.T) {
	t.Parallel()
	const mockUserAuthToken = "testuser"
	const mockDatasetID = "bar"
	const mockEdition = "baz"
	const mockVersionNum = "1"
	const mockCollectionId = "test-collection"
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	Convey("test setStatusCode", t, func() {
		Convey("test status code handles 404 response from client", func() {
			req := httptest.NewRequest("GET", "http://localhost:24000", nil)
			w := httptest.NewRecorder()
			err := &testCliError{}
			setErrorStatusCode(req, w, err, mockDatasetID)

			So(w.Code, ShouldEqual, http.StatusNotFound)
		})
	})

	Convey("test getEditMetadataHandler", t, func() {

		mockDatasetDetails := dataset.DatasetDetails{
			ID:           "test-dataset",
			CollectionID: mockCollectionId,
			Links:        dataset.Links{LatestVersion: dataset.Link{URL: "/v1/datasets/test/editions/test/version/1"}},
		}

		mockDataset := dataset.Dataset{
			Current: &mockDatasetDetails,
			Next:    &mockDatasetDetails,
		}

		mockVersionDetails := dataset.Version{
			ID:      "test-version",
			Version: 1,
		}

		datasetCollectionItem := zebedee.CollectionItem{
			ID:           mockDatasetDetails.ID,
			State:        "inProgress",
			LastEditedBy: "an-user",
		}

		mockCollection := zebedee.Collection{
			ID: mockCollectionId,
			Datasets: []zebedee.CollectionItem{

				{
					ID:           "foo",
					State:        "reviewed",
					LastEditedBy: "other-user",
				},
				datasetCollectionItem,
			},
		}

		responseHeaders := dataset.ResponseHeaders{ETag: "version-etag"}

		mockZebedeeClient := &ZebedeeClientMock{
			GetCollectionFunc: func(ctx context.Context, userAccessToken, collectionID string) (c zebedeeclient.Collection, err error) {
				if collectionID == mockCollectionId {
					return mockCollection, nil
				} else {
					return c, errors.New("collection not found")
				}
			},
		}

		mockDatasetClient := &DatasetClientMock{
			GetDatasetCurrentAndNextFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID string) (m datasetclient.Dataset, err error) {
				return mockDataset, nil
			},
			GetVersionFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, downloadServiceAuthToken, collectionID, datasetID, edition, version string) (datasetclient.Version, error) {
				return mockVersionDetails, nil
			},
			GetVersionWithHeadersFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, downloadServiceAuthToken, collectionID, datasetID, edition, version string) (datasetclient.Version, datasetclient.ResponseHeaders, error) {
				return mockVersionDetails, responseHeaders, nil
			},
		}

		Convey("when Version.State is NOT edition-confirmed returns correctly with empty dimensions struct", func() {
			mockVersionDetails.State = "associated"

			req := httptest.NewRequest("GET", "/datasets/bar/editions/baz/versions/1", nil)
			req.Header.Set("Collection-Id", mockCollectionId)
			req.Header.Set("X-Florence-Token", mockUserAuthToken)
			w := doTestRequest("/datasets/{datasetID}/editions/{editionID}/versions/{versionID}", req, GetMetadataHandler(mockDatasetClient, mockZebedeeClient), nil)

			So(w.Code, ShouldEqual, http.StatusOK)
			So(w.Body.String(), ShouldNotBeNil)

			var body model.EditMetadata
			err := json.Unmarshal(w.Body.Bytes(), &body)
			So(err, ShouldBeNil)
			So(body.Version, ShouldResemble, mockVersionDetails)
			So(body.Dataset, ShouldResemble, *mockDataset.Next)
			So(body.Dimensions, ShouldBeEmpty)
			So(body.VersionEtag, ShouldEqual, responseHeaders.ETag)
			So(body.CollectionID, ShouldEqual, mockCollectionId)
			So(body.CollectionState, ShouldEqual, datasetCollectionItem.State)
			So(body.CollectionLastEditedBy, ShouldEqual, datasetCollectionItem.LastEditedBy)
		})

		Convey("when Version.State is edition-confirmed returns correctly with populated dimensions struct", func() {
			mockVersionDetails.State = "edition-confirmed"
			mockVersionDetails.Version = 2

			mockVersionDetails.Dimensions = []datasetclient.VersionDimension{{ID: "dim001", Label: "Test dimension"}}

			req := httptest.NewRequest("GET", "/datasets/bar/editions/baz/versions/1", nil)
			req.Header.Set("Collection-Id", mockCollectionId)
			req.Header.Set("X-Florence-Token", mockUserAuthToken)
			w := doTestRequest("/datasets/{datasetID}/editions/{editionID}/versions/{versionID}", req, GetMetadataHandler(mockDatasetClient, mockZebedeeClient), nil)

			So(w.Code, ShouldEqual, http.StatusOK)
			So(w.Body.String(), ShouldNotBeNil)

			var body model.EditMetadata
			err := json.Unmarshal(w.Body.Bytes(), &body)
			So(err, ShouldBeNil)
			So(body.Version, ShouldResemble, mockVersionDetails)
			So(body.Dataset, ShouldResemble, *mockDataset.Next)
			So(body.Dimensions, ShouldResemble, mockVersionDetails.Dimensions)
			So(body.VersionEtag, ShouldEqual, responseHeaders.ETag)
			So(body.CollectionID, ShouldEqual, mockCollectionId)
			So(body.CollectionState, ShouldEqual, datasetCollectionItem.State)
			So(body.CollectionLastEditedBy, ShouldEqual, datasetCollectionItem.LastEditedBy)
		})
	})

	Convey("test getIDsFromURL", t, func() {
		expectedErr := errors.New("not enough arguements in path")
		Convey("returns error if url doesn't have enough path elements", func() {
			_, _, _, err := getIDsFromURL("https://test.ons.gov.uk/this/isnt/enough")

			So(err, ShouldResemble, expectedErr)
		})

		Convey("returns correct values", func() {
			datasetID, editionID, versionID, err := getIDsFromURL("https://test.ons.gov.uk/v1/datasets/ds1/editions/ed2/versions/1")

			So(datasetID, ShouldEqual, "ds1")
			So(editionID, ShouldEqual, "ed2")
			So(versionID, ShouldEqual, "1")
			So(err, ShouldBeNil)
		})
	})
}
