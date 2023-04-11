package dataset

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	datasetclient "github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-publishing-dataset-controller/model"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	metadataBody = `{"dataset":{"id":"test-dataset"},"version":{"id":"1"},"instance":{},"collection_id":"testcollection","collection_state":"InProgress"}`
)

func TestUnitPutMetadata(t *testing.T) {

	b := metadataBody

	Convey("test putMetadata", t, func() {
		Convey("on success", func() {

			mockDatasetClient := &DatasetClientMock{
				PutDatasetFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID string, d datasetclient.DatasetDetails) error {
					return nil
				},
				PutVersionFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID, edition, version string, v datasetclient.Version) error {
					return nil
				},
				PutInstanceFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, instanceID string, i datasetclient.UpdateInstance, ifMatch string) (string, error) {
					return "", nil
				},
			}

			mockZebedeeClient := &ZebedeeClientMock{
				PutDatasetInCollectionFunc: func(ctx context.Context, userAccessToken, collectionID, lang, datasetID, state string) error {
					return nil
				},
				PutDatasetVersionInCollectionFunc: func(ctx context.Context, userAccessToken, collectionID, lang, datasetID, edition, version, state string) error {
					return nil
				},
			}

			req := httptest.NewRequest("PUT", "/datasets/test-dataset/editions/test-edition/versions/1", bytes.NewBufferString(b))
			req.Header.Set("Collection-Id", "testcollection")
			req.Header.Set("X-Florence-Token", "testuser")
			rec := httptest.NewRecorder()
			router := mux.NewRouter()
			router.Path("/datasets/{datasetID}/editions/{editionID}/versions/{versionID}").HandlerFunc(PutMetadata(mockDatasetClient, mockZebedeeClient))

			Convey("returns 200 response", func() {
				router.ServeHTTP(rec, req)
				So(rec.Code, ShouldEqual, http.StatusOK)
			})
		})

		Convey("errors if no headers are passed", func() {

			mockDatasetClient := &DatasetClientMock{
				PutDatasetFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID string, d datasetclient.DatasetDetails) error {
					return nil
				},
				PutVersionFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID, edition, version string, v datasetclient.Version) error {
					return nil
				},
				PutInstanceFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, instanceID string, i datasetclient.UpdateInstance, ifMatch string) (string, error) {
					return "", nil
				},
			}

			mockZebedeeClient := &ZebedeeClientMock{
				PutDatasetInCollectionFunc: func(ctx context.Context, userAccessToken, collectionID, lang, datasetID, state string) error {
					return nil
				},
				PutDatasetVersionInCollectionFunc: func(ctx context.Context, userAccessToken, collectionID, lang, datasetID, edition, version, state string) error {
					return nil
				},
			}

			Convey("collection id not set", func() {
				req := httptest.NewRequest("PUT", "/datasets/test-dataset/editions/test-edition/versions/1", bytes.NewBufferString(b))
				req.Header.Set("X-Florence-Token", "testuser")
				rec := httptest.NewRecorder()
				router := mux.NewRouter()
				router.Path("/datasets/{datasetID}/editions/{editionID}/versions/{versionID}").HandlerFunc(PutMetadata(mockDatasetClient, mockZebedeeClient))

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
				req := httptest.NewRequest("PUT", "/datasets/test-dataset/editions/test-edition/versions/1", bytes.NewBufferString(b))
				req.Header.Set("Collection-Id", "testcollection")
				rec := httptest.NewRecorder()
				router := mux.NewRouter()
				router.Path("/datasets/{datasetID}/editions/{editionID}/versions/{versionID}").HandlerFunc(PutMetadata(mockDatasetClient, mockZebedeeClient))

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
				PutDatasetFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID string, d datasetclient.DatasetDetails) error {
					return errors.New("test dataset API error")
				},
				PutVersionFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID, edition, version string, v datasetclient.Version) error {
					return nil
				},
				PutInstanceFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, instanceID string, i datasetclient.UpdateInstance, ifMatch string) (string, error) {
					return "", nil
				},
			}

			mockZebedeeClient := &ZebedeeClientMock{
				PutDatasetInCollectionFunc: func(ctx context.Context, userAccessToken, collectionID, lang, datasetID, state string) error {
					return nil
				},
				PutDatasetVersionInCollectionFunc: func(ctx context.Context, userAccessToken, collectionID, lang, datasetID, edition, version, state string) error {
					return nil
				},
			}

			req := httptest.NewRequest("PUT", "/datasets/test-dataset/editions/test-edition/versions/1", bytes.NewBufferString(b))
			req.Header.Set("Collection-Id", "testcollection")
			req.Header.Set("X-Florence-Token", "testuser")
			rec := httptest.NewRecorder()
			router := mux.NewRouter()
			router.Path("/datasets/{datasetID}/editions/{editionID}/versions/{versionID}").HandlerFunc(PutMetadata(mockDatasetClient, mockZebedeeClient))

			Convey("returns 500 response and error body", func() {
				router.ServeHTTP(rec, req)
				So(rec.Code, ShouldEqual, http.StatusInternalServerError)
				response := rec.Body.String()
				So(response, ShouldResemble, "error updating dataset\n")
			})

		})
	})
}

func TestUnitPutEditableMetadata(t *testing.T) {
	Convey("Given a metadata object", t, func() {
		mockDatasetId := "test-dataset"
		mockEdition := "test-edition"
		mockVersionNumber := "1"
		mockCollectionId := "testcollection"
		etag := "versionEtag"

		metadata := model.EditMetadata{
			Dataset: datasetclient.DatasetDetails{
				ID:           mockDatasetId,
				CollectionID: mockCollectionId,
				Contacts: &[]datasetclient.Contact{{
					Name:      "contact",
					Email:     "contact@ons.gov.uk",
					Telephone: "029",
				}},
				Description: "dataset description",
				Keywords:    &[]string{"one", "two"},
				License:     "license",
				Methodologies: &[]datasetclient.Methodology{
					{
						Title:       "methodology",
						Description: "methodology description",
						URL:         "methodology url",
					},
				},
				NationalStatistic: true,
				NextRelease:       "tomorrow",
				Publications:      &[]datasetclient.Publication{},
				QMI:               datasetclient.Publication{},
				RelatedDatasets:   &[]datasetclient.RelatedDataset{},
				ReleaseFrequency:  "daily",
				Title:             "dataset title",
				UnitOfMeasure:     "unit",
				UsageNotes:        &[]datasetclient.UsageNote{},
				CanonicalTopic:    "topic",
				Subtopics:         []string{"three"},
				Survey:            "census",
				RelatedContent:    &[]datasetclient.GeneralDetails{},
			},
			Version: datasetclient.Version{
				Alerts:        &[]datasetclient.Alert{},
				CollectionID:  mockCollectionId,
				Dimensions:    []datasetclient.VersionDimension{},
				ID:            "version-id",
				LatestChanges: []datasetclient.Change{},
				Version:       1,
				UsageNotes:    &[]datasetclient.UsageNote{},
			},
			CollectionState: "in-progress",
			VersionEtag:     etag,
		}

		Convey("And a router using the PutEditableMetadata handler", func() {

			florenceToken := "testuser"

			datasetClient := &DatasetClientMock{
				PutMetadataFunc: func(ctx context.Context, userAuthToken, serviceAuthToken, collectionID, datasetID, edition, version string, editableMetadata datasetclient.EditableMetadata, versionEtag string) error {
					if userAuthToken != florenceToken || serviceAuthToken != "" {
						return errors.New("Function called with unexpected tokens")
					}
					if collectionID != mockCollectionId || datasetID != mockDatasetId || edition != mockEdition || version != mockVersionNumber {
						return errors.New("Function called with unexpected parameters")
					}
					if versionEtag != etag {
						return errors.New("Function called with invalid version etag")
					}
					return nil
				},
			}

			zebedeeClient := &ZebedeeClientMock{
				PutDatasetInCollectionFunc: func(ctx context.Context, userAccessToken, collectionID, lang, datasetID, state string) error {
					if userAccessToken != florenceToken || collectionID != mockCollectionId || lang != "" || datasetID != mockDatasetId || state != metadata.CollectionState {
						return errors.New("Function called with unexpected parameters")
					}
					return nil
				},
				PutDatasetVersionInCollectionFunc: func(ctx context.Context, userAccessToken, collectionID, lang, datasetID, edition, version, state string) error {
					if userAccessToken != florenceToken || collectionID != mockCollectionId || lang != "" || datasetID != mockDatasetId || edition != mockEdition || version != mockVersionNumber || state != metadata.CollectionState {
						return errors.New("Function called with unexpected parameters")
					}
					return nil
				},
			}

			router := mux.NewRouter()
			router.Path("/datasets/{datasetID}/editions/{editionID}/versions/{versionID}/metadata").HandlerFunc(PutEditableMetadata(datasetClient, zebedeeClient))

			rec := httptest.NewRecorder()

			body, _ := json.Marshal(metadata)
			url := fmt.Sprintf("/datasets/%s/editions/%s/versions/%s/metadata", mockDatasetId, mockEdition, mockVersionNumber)

			req := httptest.NewRequest("PUT", url, bytes.NewBuffer(body))

			Convey("When a request without a florence token header is made", func() {
				req.Header.Set("Collection-Id", mockCollectionId)

				router.ServeHTTP(rec, req)

				Convey("Then we receive a 400 response", func() {
					So(rec.Code, ShouldEqual, http.StatusBadRequest)
					So(rec.Body.String(), ShouldEqual, "no user access token header set\n")

					So(len(datasetClient.PutMetadataCalls()), ShouldEqual, 0)
					So(len(zebedeeClient.PutDatasetInCollectionCalls()), ShouldEqual, 0)
					So(len(zebedeeClient.PutDatasetVersionInCollectionCalls()), ShouldEqual, 0)
				})
			})

			Convey("When a request without a collection id header is made", func() {
				req.Header.Set("X-Florence-Token", florenceToken)

				router.ServeHTTP(rec, req)

				Convey("Then we receive a 400 response", func() {
					So(rec.Code, ShouldEqual, http.StatusBadRequest)
					So(rec.Body.String(), ShouldEqual, "no collection ID header set\n")

					So(len(datasetClient.PutMetadataCalls()), ShouldEqual, 0)
					So(len(zebedeeClient.PutDatasetInCollectionCalls()), ShouldEqual, 0)
					So(len(zebedeeClient.PutDatasetVersionInCollectionCalls()), ShouldEqual, 0)
				})
			})

			Convey("And all headers are set", func() {
				req.Header.Set("Collection-Id", mockCollectionId)
				req.Header.Set("X-Florence-Token", florenceToken)

				Convey("And the version etag is wrong", func() {
					etag = "wrong"

					Convey("When a PUT metadata request is made", func() {
						router.ServeHTTP(rec, req)

						Convey("Then we receive a 500 response", func() {
							So(rec.Code, ShouldEqual, http.StatusInternalServerError)
							So(rec.Body.String(), ShouldEqual, "error updating metadata\n")

							So(len(datasetClient.PutMetadataCalls()), ShouldEqual, 1)
							So(len(zebedeeClient.PutDatasetInCollectionCalls()), ShouldEqual, 0)
							So(len(zebedeeClient.PutDatasetVersionInCollectionCalls()), ShouldEqual, 0)
						})
					})
				})

				Convey("When a PUT metadata request is made", func() {
					router.ServeHTTP(rec, req)

					Convey("Then we receive a 200 response", func() {
						So(rec.Code, ShouldEqual, http.StatusOK)

						So(len(datasetClient.PutMetadataCalls()), ShouldEqual, 1)
						So(len(zebedeeClient.PutDatasetInCollectionCalls()), ShouldEqual, 1)
						So(len(zebedeeClient.PutDatasetVersionInCollectionCalls()), ShouldEqual, 1)
					})
				})
			})
		})
	})
}
