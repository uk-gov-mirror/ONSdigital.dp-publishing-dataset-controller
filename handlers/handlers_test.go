package handlers

import (
	"github.com/ONSdigital/dp-api-clients-go/dataset"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
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

	//target := fmt.Sprintf("/datasets/%s/editions/%s/versions/%s", mockDatasetID, mockEdition, mockVersionNum)

	Convey("test setStatusCode", t, func() {
		Convey("test status code handles 404 response from client", func() {
			req := httptest.NewRequest("GET", "http://localhost:24000", nil)
			w := httptest.NewRecorder()
			err := &testCliError{}
			setStatusCode(req, w, err)

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
			mockDatasetClient := NewMockDatasetClient(mockCtrl)
			mockDatasetClient.EXPECT().Get(gomock.Any(), mockUserAuthToken, mockServiceAuthToken, mockCollectionID, mockDatasetID).Return(mockDatasetDetails, nil)
			mockDatasetClient.EXPECT().GetVersion(gomock.Any(), mockUserAuthToken, mockServiceAuthToken, mockDownloadToken, mockCollectionID, mockDatasetID, mockEdition, mockVersionNum).Return(mockVersionDetails, nil)

			req := httptest.NewRequest("GET", "/datasets/bar/editions/baz/versions/1", nil)
			w := doTestRequest("/datasets/{datasetID}/editions/{editionID}/versions/{versionID}", req, GetEditMetadataHandler(mockDatasetClient), nil)

			So(w.Code, ShouldEqual, http.StatusOK)
			So(w.Body.String(), ShouldNotBeNil)
		})
	})
}
