package mapper

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/ONSdigital/dp-api-clients-go/dataset"
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

}
