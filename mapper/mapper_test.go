package mapper

import (
	"github.com/ONSdigital/dp-publishing-dataset-controller/model"
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/dataset"
	"github.com/davecgh/go-spew/spew"
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

		spew.Dump(ds)

		mapped := AllDatasets(ds)

		spew.Dump(mapped)

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
	mockDatasetDetails := dataset.DatasetDetails{
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

	expectedNotices := []model.Notice{
		{
			ID:                    0,
			Date:                  "04 Feb 2020",
			Description:           "Bar",
			SimpleListHeading:     "bAz (04 Feb 2020)",
			SimpleListDescription: "Bar",
			Type:                  "bAz",
		},
		{
			ID:                    1,
			Date:                  "02 Apr 2001",
			Description:           "quux",
			SimpleListHeading:     "grault (02 Apr 2001)",
			SimpleListDescription: "quux",
			Type:                  "grault",
		},
	}

	expectedUsageNotes := []model.UsageNote{
		{
			ID:                    0,
			Title:                 mockUsageNotes[0].Title,
			Note:                  mockUsageNotes[0].Note,
			SimpleListHeading:     mockUsageNotes[0].Title,
			SimpleListDescription: mockUsageNotes[0].Note,
		},
		{
			ID:                    1,
			Title:                 mockUsageNotes[1].Title,
			Note:                  mockUsageNotes[1].Note,
			SimpleListHeading:     mockUsageNotes[1].Title,
			SimpleListDescription: mockUsageNotes[1].Note,
		},
	}

	expectedLatestChanges := []model.LatestChanges{
		{
			ID:                    0,
			Title:                 mockLatestChanges[0].Name,
			Description:           mockLatestChanges[0].Description,
			SimpleListHeading:     mockLatestChanges[0].Name,
			SimpleListDescription: mockLatestChanges[0].Description,
		},
		{
			ID:                    1,
			Title:                 mockLatestChanges[1].Name,
			Description:           mockLatestChanges[1].Description,
			SimpleListHeading:     mockLatestChanges[1].Name,
			SimpleListDescription: mockLatestChanges[1].Description,
		},
	}

	expectedRelatedDataset := []model.RelatedContent{
		{
			ID:                    0,
			Title:                 mockRelatedDataset[0].Title,
			Description:           "",
			Href:                  mockRelatedDataset[0].URL,
			SimpleListHeading:     mockRelatedDataset[0].Title,
			SimpleListDescription: "",
		},
		{
			ID:                    1,
			Title:                 mockRelatedDataset[1].Title,
			Description:           "",
			Href:                  mockRelatedDataset[1].URL,
			SimpleListHeading:     mockRelatedDataset[1].Title,
			SimpleListDescription: "",
		},
	}

	expectedRelatedMethodologies := []model.RelatedContent{
		{
			ID:                    0,
			Title:                 mockMethodology[0].Title,
			Description:           mockMethodology[0].Description,
			Href:                  mockMethodology[0].URL,
			SimpleListHeading:     mockMethodology[0].Title,
			SimpleListDescription: mockMethodology[0].Description,
		},
		{
			ID:                    1,
			Title:                 mockMethodology[1].Title,
			Description:           mockMethodology[1].Description,
			Href:                  mockMethodology[1].URL,
			SimpleListHeading:     mockMethodology[1].Title,
			SimpleListDescription: mockMethodology[1].Description,
		},
	}

	expectedRelatedPublcation := []model.RelatedContent{
		{
			ID:                    0,
			Title:                 mockPublications[0].Title,
			Description:           mockPublications[0].Description,
			Href:                  mockPublications[0].URL,
			SimpleListHeading:     mockPublications[0].Title,
			SimpleListDescription: mockPublications[0].Description,
		},
		{
			ID:                    1,
			Title:                 mockPublications[1].Title,
			Description:           mockPublications[1].Description,
			Href:                  mockPublications[1].URL,
			SimpleListHeading:     mockPublications[1].Title,
			SimpleListDescription: mockPublications[1].Description,
		},
	}

	expectedEditVersionMetaData := model.EditVersionMetaData{
		MetaData: model.MetaData{
			Edition: mockVersion.Edition,
			Version: mockVersion.Version,
			ReleaseDate: model.ReleaseDate{
				ReleaseDate: mockVersion.ReleaseDate,
				Error:       "",
			},
			Notices:              expectedNotices,
			Dimensions:           mockVersion.Dimensions,
			UsageNotes:           expectedUsageNotes,
			LatestChanges:        expectedLatestChanges,
			Title:                mockDatasetDetails.Title,
			Summary:              mockDatasetDetails.Description,
			Keywords:             "foo, Bar, bAz",
			NationalStatistic:    mockDatasetDetails.NationalStatistic,
			License:              mockDatasetDetails.License,
			ContactName:          mockContacts[0].Name,
			ContactEmail:         mockContacts[0].Email,
			ContactTelephone:     mockContacts[0].Telephone,
			RelatedDatasets:      expectedRelatedDataset,
			RelatedPublications:  expectedRelatedPublcation,
			RelatedMethodologies: expectedRelatedMethodologies,
			ReleaseFrequency:     mockDatasetDetails.ReleaseFrequency,
			NextReleaseDate:      mockDatasetDetails.NextRelease,
			UnitOfMeassure:       mockDatasetDetails.UnitOfMeasure,
			QMI:                  mockDatasetDetails.QMI.URL,
		},
		Collection: mockVersion.CollectionID,
		InstanceID: mockVersion.ID,
		Published:  mockVersion.State == "published",
	}

	Convey("test EditDatasetVersionMetaData", t, func() {
		Convey("when working", func() {
			outcome, err := EditDatasetVersionMetaData(mockDatasetDetails, mockVersion)
			So(err, ShouldBeNil)
			So(outcome, ShouldResemble, expectedEditVersionMetaData)
		})
	})
}
