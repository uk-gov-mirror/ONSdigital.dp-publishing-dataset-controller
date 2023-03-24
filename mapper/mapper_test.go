package mapper

import (
	"context"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	zebedee "github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	babbage "github.com/ONSdigital/dp-publishing-dataset-controller/clients/topics"
	"github.com/ONSdigital/dp-publishing-dataset-controller/model"
	. "github.com/smartystreets/goconvey/convey"
)

var ctx = context.Background()

func TestUnitMapper(t *testing.T) {
	t.Parallel()
	Convey("test AllDatasets", t, func() {
		ds := dataset.List{
			Items: []dataset.Dataset{},
		}
		ds.Items = append(ds.Items, dataset.Dataset{
			ID: "test-id-1",
			Next: &dataset.DatasetDetails{
				Title: "test title 1",
			},
		}, dataset.Dataset{
			ID: "test-id-2",
			Next: &dataset.DatasetDetails{
				Title: "test title 2",
			},
		}, dataset.Dataset{
			ID: "test-id-3",
		})

		mapped := AllDatasets(ds)

		So(mapped[0].ID, ShouldEqual, "test-id-1")
		So(mapped[0].Title, ShouldEqual, "test title 1")
		So(mapped[1].ID, ShouldEqual, "test-id-2")
		So(mapped[1].Title, ShouldEqual, "test title 2")
		So(len(mapped), ShouldEqual, 2)
	})

	Convey("that datasets are ordered alphabetically by Title", t, func() {
		ds := dataset.List{
			Items: []dataset.Dataset{},
		}
		ds.Items = append(ds.Items, dataset.Dataset{
			ID: "test-id-3",
			Next: &dataset.DatasetDetails{
				Title: "3rd Title",
			},
		}, dataset.Dataset{
			ID: "test-id-1",
			Next: &dataset.DatasetDetails{
				Title: "1st Title",
			},
		}, dataset.Dataset{
			ID: "test-id-2",
			Next: &dataset.DatasetDetails{
				Title: "2nd Title",
			},
		})

		mapped := AllDatasets(ds)

		So(mapped[0].ID, ShouldEqual, "test-id-1")
		So(mapped[1].ID, ShouldEqual, "test-id-2")
		So(mapped[2].ID, ShouldEqual, "test-id-3")
		So(len(mapped), ShouldEqual, 3)
	})

	Convey("that datasets with an empty title are still sorted alphabetically using their ID instead", t, func() {
		ds := dataset.List{
			Items: []dataset.Dataset{},
		}
		ds.Items = append(ds.Items, dataset.Dataset{
			ID: "test-id-4",
			Next: &dataset.DatasetDetails{
				Title: "DFG",
			},
		}, dataset.Dataset{
			ID: "test-id-1",
			Next: &dataset.DatasetDetails{
				Title: "",
			},
		}, dataset.Dataset{
			ID: "test-id-2",
			Next: &dataset.DatasetDetails{
				Title: "",
			},
		}, dataset.Dataset{
			ID: "test-id-3",
			Next: &dataset.DatasetDetails{
				Title: "ABC",
			},
		})

		mapped := AllDatasets(ds)

		So(mapped[0].ID, ShouldEqual, "test-id-3")
		So(mapped[1].ID, ShouldEqual, "test-id-4")
		So(mapped[2].ID, ShouldEqual, "test-id-1")
		So(mapped[3].ID, ShouldEqual, "test-id-2")
		So(len(mapped), ShouldEqual, 4)
	})

	Convey("that datasets are ordered correctly regardless of casing in the ID or Title fields", t, func() {
		ds := dataset.List{
			Items: []dataset.Dataset{},
		}
		ds.Items = append(ds.Items, dataset.Dataset{
			ID: "test-id-4",
			Next: &dataset.DatasetDetails{
				Title: "dfg",
			},
		}, dataset.Dataset{
			ID: "Test-id-1",
			Next: &dataset.DatasetDetails{
				Title: "",
			},
		}, dataset.Dataset{
			ID: "test-id-2",
			Next: &dataset.DatasetDetails{
				Title: "ABC",
			},
		}, dataset.Dataset{
			ID: "test-id-3",
			Next: &dataset.DatasetDetails{
				Title: "123",
			},
		})

		mapped := AllDatasets(ds)

		So(mapped[0].ID, ShouldEqual, "test-id-3")
		So(mapped[1].ID, ShouldEqual, "test-id-2")
		So(mapped[2].ID, ShouldEqual, "test-id-4")
		So(mapped[3].ID, ShouldEqual, "Test-id-1")
		So(len(mapped), ShouldEqual, 4)
	})

	mockTopics := babbage.TopicsResult{
		Topics: babbage.Topic{
			Results: []babbage.Result{{
				Description: babbage.Description{
					Title: "test 1",
				},
				URI:  "/test/uri/1",
				Type: "page",
			},
				{
					Description: babbage.Description{
						Title: "test 2",
					},
					URI:  "/test/uri/2",
					Type: "page",
				}},
		},
	}

	mockEmptyTopics := babbage.TopicsResult{
		Topics: babbage.Topic{
			Results: []babbage.Result{},
		},
	}

	expectedTopics := []model.Topics{{Title: "test 1"}, {Title: "test 2"}}

	Convey("test Topics", t, func() {
		Convey("maps correctly if results have topics", func() {
			outcome := Topics(mockTopics)
			So(outcome, ShouldResemble, expectedTopics)
		})
		Convey("retruns empty slice and doesn't error if no results", func() {
			outcome := Topics(mockEmptyTopics)
			So(outcome, ShouldResemble, []model.Topics(nil))
		})
	})

	mockedAllVersions := dataset.VersionsList{
		Items: []dataset.Version{},
	}
	mockedAllVersions.Items = append(mockedAllVersions.Items, dataset.Version{
		ID:          "test-id-3",
		Version:     3,
		ReleaseDate: "",
		State:       "edition-confirmed",
	}, dataset.Version{
		ID:          "test-id-1",
		Version:     1,
		ReleaseDate: "2020-11-07T00:00:00.000Z",
		State:       "published",
	}, dataset.Version{
		ID:          "test-id-2",
		Version:     2,
		ReleaseDate: "2020-11-20T00:00:00.000Z",
		State:       "published",
	})

	mockedDataset := dataset.Dataset{
		Next: &dataset.DatasetDetails{
			Title: "Test title",
		},
	}

	mockedEdition := dataset.Edition{
		Edition: "edition-1",
	}

	expectedAllVersions := []model.Version{{ID: "test-id-3", Title: "Version: 3", Version: 3, ReleaseDate: "", State: "edition-confirmed"}, {ID: "test-id-2", Title: "Version: 2 (published)", Version: 2, ReleaseDate: "20 November 2020", State: "published"}, {ID: "test-id-1", Title: "Version: 1 (published)", Version: 1, ReleaseDate: "07 November 2020", State: "published"}}

	expectedVersionsPage := model.VersionsPage{DatasetName: "Test title", EditionName: "edition-1", Versions: expectedAllVersions}

	Convey("test AllVersions", t, func() {
		Convey("maps correctly", func() {
			mapped := AllVersions(ctx, mockedDataset, mockedEdition, mockedAllVersions)
			So(mapped, ShouldResemble, expectedVersionsPage)
		})
	})
}

func TestMetadata(t *testing.T) {
	Convey("Given a dataset and version objects", t, func() {
		mockDatasetDetails := &dataset.DatasetDetails{
			ID:           "foo",
			CollectionID: "Bar",
			Contacts: &[]dataset.Contact{
				{
					Name:      "foo",
					Telephone: "Bar",
					Email:     "bAz",
				},
				{
					Name:      "bad-foo",
					Telephone: "bad-Bar",
					Email:     "bad-bAz",
				},
			},
			Description: "bAz",
			Keywords:    &[]string{"foo", "Bar", "bAz"},
			License:     "qux",
			Links:       dataset.Links{},
			Methodologies: &[]dataset.Methodology{
				{
					Description: "foo",
					URL:         "Bar",
					Title:       "bAz",
				},
				{
					Description: "qux",
					URL:         "quux",
					Title:       "grault",
				},
			},
			NationalStatistic: false,
			NextRelease:       "quux",
			Publications: &[]dataset.Publication{
				{
					Description: "Bar",
					URL:         "bAz",
					Title:       "foo",
				},
				{
					Description: "quux",
					URL:         "grault",
					Title:       "qux",
				},
			},
			Publisher: &dataset.Publisher{},
			QMI: dataset.Publication{
				Description: "foo",
				URL:         "Bar",
				Title:       "bAz",
			},
			RelatedDatasets: &[]dataset.RelatedDataset{
				{
					URL:   "foo",
					Title: "Bar",
				},
				{
					URL:   "bAz",
					Title: "qux",
				},
			},
			ReleaseFrequency: "grault",
			State:            "garply",
			Theme:            "waldo",
			Title:            "fred",
			UnitOfMeasure:    "plugh",
			URI:              "xyzzy",
			CanonicalTopic:   "1234",
			Subtopics:        []string{"5678", "9012"},
			RelatedContent: &[]dataset.GeneralDetails{
				{
					Description: "foo",
					HRef:        "Bar",
					Title:       "baz",
				},
				{
					Description: "foo",
					HRef:        "Bar",
					Title:       "baz",
				},
			},
			Survey: "census",
		}
		mockDimensions := []dataset.VersionDimension{
			{
				Links: dataset.Links{},
				Label: "bAz",
			},
			{
				Links: dataset.Links{},
				Label: "plaugh",
			},
		}

		mockVersion := dataset.Version{
			Alerts: &[]dataset.Alert{
				{
					Date:        "2020-02-04T11:05:06.000Z",
					Description: "Bar",
					Type:        "bAz",
				},
				{
					Date:        "2001-04-02T23:04:02.000Z",
					Description: "quux",
					Type:        "grault",
				},
			},
			CollectionID: "foo",
			Downloads:    nil,
			Edition:      "Bar",
			Dimensions:   mockDimensions,
			ID:           "bAz",
			InstanceID:   "qux",
			LatestChanges: []dataset.Change{
				{
					Description: "foo",
					Name:        "Bar",
					Type:        "bAz",
				},
				{
					Description: "qux",
					Name:        "quux",
					Type:        "grault",
				},
			},
			Links:       dataset.Links{},
			ReleaseDate: "grault",
			State:       "grault",
			Temporal:    nil,
			Version:     1,
			UsageNotes: &[]dataset.UsageNote{
				{
					Title: "foo",
					Note:  "Bar",
				},
				{
					Title: "bAz",
					Note:  "qux",
				},
			},
		}

		Convey("And a zebedee collection", func() {

			datasetCollectionItem := zebedee.CollectionItem{
				ID:           mockDatasetDetails.ID,
				State:        "inProgress",
				LastEditedBy: "User",
			}
			mockCollection := zebedee.Collection{
				ID: "test-collection",
				Datasets: []zebedee.CollectionItem{
					{
						ID:           "other dataset id",
						State:        "reviewd",
						LastEditedBy: "Other user",
					},
					datasetCollectionItem,
				},
			}
			Convey("When we call EditMetadata", func() {
				outcome := EditMetadata(mockDatasetDetails, mockVersion, mockDimensions, mockCollection)
				Convey("Then it returns an object with all the EditMetadata fields populated", func() {
					expectedEditMetadata := model.EditMetadata{
						Dataset:                *mockDatasetDetails,
						Version:                mockVersion,
						Dimensions:             mockDimensions,
						CollectionID:           mockCollection.ID,
						CollectionState:        datasetCollectionItem.State,
						CollectionLastEditedBy: datasetCollectionItem.LastEditedBy,
					}
					So(outcome, ShouldResemble, expectedEditMetadata)
				})
			})
		})

		Convey("And an empty EditMetadata", func() {
			editMetadata := model.EditMetadata{}
			Convey("When we call PutMetadata", func() {

				editableMetadataObj := PutMetadata(editMetadata)

				Convey("Then it returns an object with all the editable metadata fields populated", func() {
					So(editableMetadataObj.Description, ShouldBeEmpty)
					So(editableMetadataObj.Keywords, ShouldBeEmpty)
					So(editableMetadataObj.Title, ShouldBeEmpty)
					So(editableMetadataObj.UnitOfMeasure, ShouldBeEmpty)
					So(editableMetadataObj.Contacts, ShouldBeEmpty)
					So(editableMetadataObj.QMI.Description, ShouldBeEmpty)
					So(editableMetadataObj.QMI.Title, ShouldBeEmpty)
					So(editableMetadataObj.QMI.URL, ShouldBeEmpty)
					So(editableMetadataObj.RelatedContent, ShouldBeEmpty)
					So(editableMetadataObj.CanonicalTopic, ShouldBeEmpty)
					So(editableMetadataObj.Subtopics, ShouldBeEmpty)
					So(editableMetadataObj.License, ShouldBeEmpty)
					So(editableMetadataObj.Methodologies, ShouldBeEmpty)
					So(*editableMetadataObj.NationalStatistic, ShouldBeFalse)
					So(editableMetadataObj.NextRelease, ShouldBeEmpty)
					So(editableMetadataObj.Publications, ShouldBeEmpty)
					So(editableMetadataObj.RelatedDatasets, ShouldBeEmpty)
					So(editableMetadataObj.ReleaseFrequency, ShouldBeEmpty)
					So(editableMetadataObj.Survey, ShouldBeEmpty)

					So(editableMetadataObj.Dimensions, ShouldBeEmpty)
					So(editableMetadataObj.ReleaseDate, ShouldBeEmpty)
					So(editableMetadataObj.Alerts, ShouldBeNil)
					So(editableMetadataObj.LatestChanges, ShouldBeEmpty)
					So(editableMetadataObj.UsageNotes, ShouldBeNil)
				})
			})
		})

		Convey("And an EditMetadata object with full dataset and version", func() {
			editMetadata := model.EditMetadata{
				Dataset: *mockDatasetDetails,
				Version: mockVersion,
			}
			Convey("When we call PutMetadata", func() {

				editableMetadataObj := PutMetadata(editMetadata)

				Convey("Then it returns an object with all the editable metadata fields populated", func() {
					So(editableMetadataObj.Description, ShouldEqual, editMetadata.Dataset.Description)
					So(editableMetadataObj.Keywords, ShouldResemble, *editMetadata.Dataset.Keywords)
					So(editableMetadataObj.Title, ShouldEqual, editMetadata.Dataset.Title)
					So(editableMetadataObj.UnitOfMeasure, ShouldEqual, editMetadata.Dataset.UnitOfMeasure)
					So(editableMetadataObj.Contacts, ShouldResemble, *editMetadata.Dataset.Contacts)
					So(editableMetadataObj.QMI, ShouldResemble, &editMetadata.Dataset.QMI)
					So(editableMetadataObj.RelatedContent, ShouldResemble, *editMetadata.Dataset.RelatedContent)
					So(editableMetadataObj.CanonicalTopic, ShouldEqual, editMetadata.Dataset.CanonicalTopic)
					So(editableMetadataObj.Subtopics, ShouldResemble, editMetadata.Dataset.Subtopics)
					So(editableMetadataObj.License, ShouldResemble, editMetadata.Dataset.License)
					So(editableMetadataObj.Methodologies, ShouldResemble, *editMetadata.Dataset.Methodologies)
					So(editableMetadataObj.NationalStatistic, ShouldResemble, &editMetadata.Dataset.NationalStatistic)
					So(editableMetadataObj.NextRelease, ShouldResemble, editMetadata.Dataset.NextRelease)
					So(editableMetadataObj.Publications, ShouldResemble, *editMetadata.Dataset.Publications)
					So(editableMetadataObj.RelatedDatasets, ShouldResemble, *editMetadata.Dataset.RelatedDatasets)
					So(editableMetadataObj.ReleaseFrequency, ShouldResemble, editMetadata.Dataset.ReleaseFrequency)
					So(editableMetadataObj.Survey, ShouldEqual, editMetadata.Dataset.Survey)

					So(editableMetadataObj.Dimensions, ShouldResemble, editMetadata.Version.Dimensions)
					So(editableMetadataObj.ReleaseDate, ShouldEqual, editMetadata.Version.ReleaseDate)
					So(editableMetadataObj.Alerts, ShouldEqual, editMetadata.Version.Alerts)
					So(editableMetadataObj.LatestChanges, ShouldResemble, &editMetadata.Version.LatestChanges)
					So(editableMetadataObj.UsageNotes, ShouldEqual, editMetadata.Version.UsageNotes)
				})
			})
		})
	})
}
