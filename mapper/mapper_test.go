package mapper

import (
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/dataset"
	zebedee "github.com/ONSdigital/dp-api-clients-go/zebedee"
	babbage "github.com/ONSdigital/dp-publishing-dataset-controller/clients/topics"
	"github.com/ONSdigital/dp-publishing-dataset-controller/model"
	. "github.com/smartystreets/goconvey/convey"
)

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

	mockContacts := []dataset.Contact{
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
	}
	mockMethodology := []dataset.Methodology{
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
	}
	mockPublications := []dataset.Publication{
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
	}
	mockRelatedDataset := []dataset.RelatedDataset{
		{
			URL:   "foo",
			Title: "Bar",
		},
		{
			URL:   "bAz",
			Title: "qux",
		},
	}
	mockKeywords := []string{"foo", "Bar", "bAz"}
	mockUsageNotes := []dataset.UsageNote{
		{
			Title: "foo",
			Note:  "Bar",
		},
		{
			Title: "bAz",
			Note:  "qux",
		},
	}
	mockDatasetDetails := &dataset.DatasetDetails{
		ID:                "foo",
		CollectionID:      "Bar",
		Contacts:          &mockContacts,
		Description:       "bAz",
		Keywords:          &mockKeywords,
		License:           "qux",
		Links:             dataset.Links{},
		Methodologies:     &mockMethodology,
		NationalStatistic: false,
		NextRelease:       "quux",
		Publications:      &mockPublications,
		Publisher:         &dataset.Publisher{},
		QMI: dataset.Publication{
			Description: "foo",
			URL:         "Bar",
			Title:       "bAz",
		},
		RelatedDatasets:  &mockRelatedDataset,
		ReleaseFrequency: "grault",
		State:            "garply",
		Theme:            "waldo",
		Title:            "fred",
		UnitOfMeasure:    "plugh",
		URI:              "xyzzy",
		UsageNotes:       &mockUsageNotes,
	}
	mockAlerts := []dataset.Alert{
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

	mockLatestChanges := []dataset.Change{
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
	}
	mockVersion := dataset.Version{
		Alerts:        &mockAlerts,
		CollectionID:  "foo",
		Downloads:     nil,
		Edition:       "Bar",
		Dimensions:    mockDimensions,
		ID:            "bAz",
		InstanceID:    "qux",
		LatestChanges: mockLatestChanges,
		Links:         dataset.Links{},
		ReleaseDate:   "quux",
		State:         "grault",
		Temporal:      nil,
		Version:       1,
	}
	mockDimensions = []dataset.VersionDimension{}

	mockCollection := zebedee.Collection{
		ID: "test-collection",
		Datasets: []zebedee.CollectionItem{
			{
				ID:    "foo",
				State: "inProgress",
			},
		},
	}

	expectedEditMetadata := model.EditMetadata{
		Dataset:         *mockDatasetDetails,
		Version:         mockVersion,
		Dimensions:      mockDimensions,
		CollectionID:    "test-collection",
		CollectionState: "inProgress",
	}

	Convey("test EditMetadata", t, func() {
		outcome := EditMetadata(mockDatasetDetails, mockVersion, mockDimensions, mockCollection)
		So(outcome, ShouldResemble, expectedEditMetadata)
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
}
