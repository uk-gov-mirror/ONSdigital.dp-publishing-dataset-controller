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

		So(mapped[0].Title, ShouldEqual, "1st Title")
		So(mapped[1].Title, ShouldEqual, "2nd Title")
		So(mapped[2].Title, ShouldEqual, "3rd Title")
		So(len(mapped), ShouldEqual, 3)
	})

	Convey("that datasets with an empty title are pushed to the end of the datasets slice", t, func() {
		ds := dataset.List{
			Items: []dataset.Dataset{},
		}
		ds.Items = append(ds.Items, dataset.Dataset{
			ID: "test-id-3",
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
			ID: "test-id-2",
			Next: &dataset.DatasetDetails{
				Title: "ABC",
			},
		})

		mapped := AllDatasets(ds)

		So(mapped[0].Title, ShouldEqual, "ABC")
		So(mapped[1].Title, ShouldEqual, "DFG")
		So(mapped[2].Title, ShouldEqual, "")
		So(mapped[3].Title, ShouldEqual, "")
		So(len(mapped), ShouldEqual, 4)
	})

	mockContacts := []dataset.Contact{
		{
			Name:      "foo",
			Telephone: "bar",
			Email:     "baz",
		},
		{
			Name:      "bad-foo",
			Telephone: "bad-bar",
			Email:     "bad-baz",
		},
	}
	mockMethodology := []dataset.Methodology{
		{
			Description: "foo",
			URL:         "bar",
			Title:       "baz",
		},
		{
			Description: "qux",
			URL:         "quux",
			Title:       "grault",
		},
	}
	mockPublications := []dataset.Publication{
		{
			Description: "bar",
			URL:         "baz",
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
			URL:         "foo",
			Title:       "bar",
		},
		{
			URL:         "baz",
			Title:       "qux",
		},
	}
	mockKeywords := []string{"foo", "bar", "baz"}
	mockUsageNotes := []dataset.UsageNote{
		{
			Title: "foo",
			Note: "bar",
		},
		{
			Title: "baz",
			Note: "qux",
		},
	}
	mockDatasetDetails := dataset.DatasetDetails{
		ID:                "foo",
		CollectionID:      "bar",
		Contacts:          &mockContacts,
		Description:       "baz",
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
			URL:         "bar",
			Title:       "baz",
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
	mockVersion := dataset.Version{
		Alerts:        nil,
		CollectionID:  "foo",
		Downloads:     nil,
		Edition:       "bar",
		Dimensions:    nil,
		ID:            "baz",
		InstanceID:    "qux",
		LatestChanges: nil,
		Links:         dataset.Links{},
		ReleaseDate:   "quux",
		State:         "grault",
		Temporal:      nil,
		Version:       1,
	}

	expectedEditVersionMetaData := model.EditVersionMetaData{
		MetaData: model.MetaData{
			Edition: "bar",
			Version: 1,
			ReleaseDate: model.ReleaseDate{
				ReleaseDate: "quux",
				Error:       "",
			},
			Notices:              nil,
			Dimensions:           nil,
			UsageNotes:           nil,
			LatestChanges:        nil,
			Title:                "fred",
			Summary:              "baz",
			Keywords:             "",
			NationalStatistic:    false,
			License:              "qux",
			ContactName:          "foo",
			ContactEmail:         "bar",
			ContactTelephone:     "baz",
			RelatedDatasets:      nil,
			RelatedPublications:  nil,
			RelatedMethodologies: nil,
			ReleaseFrequency:     "grault",
			NextReleaseDate:      "quux",
			UnitOfMeassure:       "plugh",
			QMI:                  "bar",
		},
		Collection: "bar",
		InstanceID: "foo",
		Published:  false,
	}

	Convey("test EditDatasetVersionMetaData", t, func() {
		Convey("when working", func() {
			outcome := EditDatasetVersionMetaData(mockDatasetDetails, mockVersion)
			So(outcome, ShouldResemble, expectedEditVersionMetaData)
		})
	})
}
