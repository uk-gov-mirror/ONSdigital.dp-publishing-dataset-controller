package dataset

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/dataset"
	datasetclient "github.com/ONSdigital/dp-api-clients-go/dataset"
	"github.com/ONSdigital/dp-api-clients-go/zebedee"
	zebedeeclient "github.com/ONSdigital/dp-api-clients-go/zebedee"
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
		Convey("when working", func() {

			mockDatasetDetails := dataset.DatasetDetails{
				ID:                "",
				CollectionID:      "",
				Contacts:          nil,
				Description:       "",
				Keywords:          nil,
				License:           "",
				Links:             dataset.Links{},
				Methodologies:     nil,
				NationalStatistic: false,
				NextRelease:       "",
				Publications:      nil,
				Publisher:         nil,
				QMI:               dataset.Publication{},
				RelatedDatasets:   nil,
				ReleaseFrequency:  "",
				State:             "",
				Theme:             "",
				Title:             "",
				UnitOfMeasure:     "",
				URI:               "",
				UsageNotes:        nil,
			}

			mockVersionDetails := dataset.Version{
				Alerts:        nil,
				CollectionID:  "",
				Downloads:     nil,
				Edition:       "",
				Dimensions:    nil,
				ID:            "",
				InstanceID:    "",
				LatestChanges: nil,
				Links:         dataset.Links{},
				ReleaseDate:   "",
				State:         "",
				Temporal:      nil,
				Version:       0,
			}

			mockInstance := dataset.Instance{
				mockVersionDetails,
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
			mockDatasetClient := &DatasetClientMock{
				GetFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID string) (m datasetclient.DatasetDetails, err error) {
					return mockDatasetDetails, nil
				},
				GetVersionFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, downloadServiceAuthToken, collectionID, datasetID, edition, version string) (m datasetclient.Version, err error) {
					return mockVersionDetails, nil
				},
				GetInstanceFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, instanceID string) (i datasetclient.Instance, err error) {
					return mockInstance, nil
				},
			}

			mockZebedeeClient := &ZebedeeClientMock{
				GetCollectionFunc: func(ctx context.Context, userAccessToken, collectionID string) (c zebedeeclient.Collection, err error) {
					return mockCollection, nil
				},
			}

			req := httptest.NewRequest("GET", "/datasets/bar/editions/baz/versions/1", nil)
			req.Header.Set("Collection-Id", "testcollection")
			req.Header.Set("X-Florence-Token", "testuser")
			w := doTestRequest("/datasets/{datasetID}/editions/{editionID}/versions/{versionID}", req, GetMetadataHandler(mockDatasetClient, mockZebedeeClient), nil)

			So(w.Code, ShouldEqual, http.StatusOK)
			So(w.Body.String(), ShouldNotBeNil)
		})
	})
}
