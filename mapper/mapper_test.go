package mapper

import (
	"testing"

	"github.com/ONSdigital/dp-api-clients-go/dataset"
	"github.com/davecgh/go-spew/spew"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitMapper(t *testing.T) {
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

}
