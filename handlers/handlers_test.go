package handlers

import (
	"fmt"
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
	mockUserAuthToken := ""
	mockServiceAuthToken := ""
	mockDownloadToken := ""
	mockCollectionID := ""
	mockDatasetID := "bar"
	mockEdition := "baz"
	mockVersionNum := "1"
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	target := fmt.Sprintf("/datasets/%s/editions/%s/versions/%s", mockDatasetID, mockEdition, mockVersionNum)

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
			mockDatasetClient := NewMockDatasetClient(mockCtrl)
			mockDatasetClient.EXPECT().GetVersion(gomock.Any(), mockUserAuthToken, mockServiceAuthToken, mockDownloadToken, mockCollectionID, mockDatasetID, mockEdition, mockVersionNum).Return(dataset.Version{}, nil)
			mockDatasetClient.EXPECT().Get(gomock.Any(), mockUserAuthToken, mockServiceAuthToken, mockCollectionID, mockDatasetID).Return(dataset.DatasetDetails{}, nil)

			req := httptest.NewRequest("GET", target, nil)
			w := doTestRequest(target, req, GetEditMetadataHandler(mockDatasetClient), nil)

			So(w.Code, ShouldEqual, http.StatusOK)
			So(w.Body.String(), ShouldNotBeNil)
		})
	})
}
