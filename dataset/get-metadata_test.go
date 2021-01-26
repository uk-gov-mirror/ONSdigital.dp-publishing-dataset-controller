package dataset

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/dataset"
	datasetclient "github.com/ONSdigital/dp-api-clients-go/dataset"
	"github.com/ONSdigital/dp-api-clients-go/zebedee"
	zebedeeclient "github.com/ONSdigital/dp-api-clients-go/zebedee"
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
	const mockUserAuthToken = ""
	const mockServiceAuthToken = ""
	const mockDownloadToken = ""
	const mockCollectionID = ""
	const mockDatasetID = "bar"
	const mockEdition = "baz"
	const mockVersionNum = "1"
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
			ID:    "test-dataset",
			Links: dataset.Links{LatestVersion: dataset.Link{URL: "/v1/datasets/test/editions/test/version/1"}},
		}

		mockVersionDetails := dataset.Version{
			ID:      "test-version",
			Version: 1,
		}

		mockCollection := zebedee.Collection{
			ID: "test-collection",
			Datasets: []zebedee.CollectionItem{
				{
					ID:    "foo",
					State: "inProgress",
				},
			},
		}

		mockZebedeeClient := &ZebedeeClientMock{
			GetCollectionFunc: func(ctx context.Context, userAccessToken, collectionID string) (c zebedeeclient.Collection, err error) {
				return mockCollection, nil
			},
		}

		Convey("when Version.State is NOT edition-confirmed returns correctly with empty dimensions struct", func() {

			mockDatasetClient := &DatasetClientMock{
				GetFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID string) (m datasetclient.DatasetDetails, err error) {
					return mockDatasetDetails, nil
				},
				GetVersionFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, downloadServiceAuthToken, collectionID, datasetID, edition, version string) (m datasetclient.Version, err error) {
					return mockVersionDetails, nil
				},
			}

			req := httptest.NewRequest("GET", "/datasets/bar/editions/baz/versions/1", nil)
			req.Header.Set("Collection-Id", "testcollection")
			req.Header.Set("X-Florence-Token", "testuser")
			w := doTestRequest("/datasets/{datasetID}/editions/{editionID}/versions/{versionID}", req, GetMetadataHandler(mockDatasetClient, mockZebedeeClient), nil)

			So(w.Code, ShouldEqual, http.StatusOK)
			So(w.Body.String(), ShouldNotBeNil)
		})

		Convey("when Version.State is edition-confirmed returns correctly with populated dimensions struct", func() {

			mockVersion := mockVersionDetails
			mockVersion.State = "edition-confirmed"
			mockVersion.Version = 2

			var count int
			mockDatasetClient := &DatasetClientMock{
				GetFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID string) (m datasetclient.DatasetDetails, err error) {
					return mockDatasetDetails, nil
				},
				GetVersionFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, downloadServiceAuthToken, collectionID, datasetID, edition, version string) (m datasetclient.Version, err error) {
					var data datasetclient.Version
					if count == 0 {
						data = mockVersion
					}
					if count == 1 {
						data = datasetclient.Version{Dimensions: []datasetclient.VersionDimension{{ID: "dim001", Label: "Test dimension"}}}
					}
					count++
					return data, nil
				},
			}

			req := httptest.NewRequest("GET", "/datasets/bar/editions/baz/versions/1", nil)
			req.Header.Set("Collection-Id", "testcollection")
			req.Header.Set("X-Florence-Token", "testuser")
			w := doTestRequest("/datasets/{datasetID}/editions/{editionID}/versions/{versionID}", req, GetMetadataHandler(mockDatasetClient, mockZebedeeClient), nil)

			So(w.Code, ShouldEqual, http.StatusOK)
			var body model.EditMetadata
			b, _ := ioutil.ReadAll(w.Body)
			_ = json.Unmarshal(b, &body)
			So(body.Dataset.ID, ShouldEqual, "test-dataset")
			So(body.Version.ID, ShouldEqual, "test-version")
			So(len(body.Dimensions), ShouldBeGreaterThan, 0)
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
