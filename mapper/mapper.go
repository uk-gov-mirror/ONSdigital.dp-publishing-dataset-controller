package mapper

import (
	"fmt"
	"github.com/ONSdigital/dp-api-clients-go/dataset"
	"github.com/ONSdigital/dp-publishing-dataset-controller/model"
	"net/http"
	"strings"
	"time"
)

func EditDatasetVersionMetaData(req *http.Request, d dataset.DatasetDetails, v dataset.Version) model.EditVersionMetaData {
	var notices []model.Notice

	releaseDate := model.ReleaseDate{
		ReleaseDate: v.ReleaseDate,
		Error:       "",
	}

	layout := "2006-01-02T15:04:05.000Z" // TODO check this is the right layout for the dates
	for i, alert := range *v.Alerts {
		alertDateInDateFormat, err := time.Parse(layout, alert.Date)
		if err != nil {
			//TODO log error with the version that failed
		}
		noticeDate := alertDateInDateFormat.Format("01 Jan 2006")
		simpleListHeading := fmt.Sprintf(`%s (%s)`, alert.Type, noticeDate)
		notices = append(notices, model.Notice{
			ID:                    i,
			Type:                  alert.Type,
			Date:                  alert.Date,
			Description:           alert.Description,
			SimpleListHeading:     simpleListHeading,
			SimpleListDescription: alert.Description,
		})
	}

	var usageNotes []model.UsageNote
	for i, note := range *d.UsageNotes {
		usageNotes = append(usageNotes, model.UsageNote{
			ID:                    i,
			Title:                 note.Title,
			Note:                  note.Note,
			SimpleListHeading:     note.Title,
			SimpleListDescription: note.Note,
		})
	}

	var latestChanges []model.LatestChanges
	for i, change := range v.LatestChanges {
		latestChanges = append(latestChanges, model.LatestChanges{
			ID:                    i,
			Title:                 change.Name,
			Description:           change.Description,
			SimpleListHeading:     change.Name,
			SimpleListDescription: change.Description,
		})
	}

	keywordsString, err := fmt.Println(strings.Join(d.Keywords[:], ", "))
	if err != nil {
		// TODO log error with the version that failed
	}
	// TODO all below

	// TODO all above

	mappedMetaData := model.MetaData{
		Edition:       v.Edition,
		Version:       v.Version,
		ReleaseDate:   releaseDate,
		Notices:       notices,
		Dimensions:    v.Dimensions,
		UsageNotes:    usageNotes,
		LatestChanges: latestChanges,

		Title:                d.Title,
		Summary:              d.Description,
		Keywords:             keywordsString,
		NationalStatistic:    d.NationalStatistic,
		License:              d.License,
		ContactName:          d.Contacts[0].Name,
		ContactEmail:         d.Contacts[0].Email,
		ContactTelephone:     d.Contacts[0].Telephone,
		RelatedDatasets:      relatedDatasets,
		RelatedPublications:  relatedPublications,
		RelatedMethodologies: relatedMethodologies,
		ReleaseFrequency:     d.ReleaseFrequency,
		NextReleaseDate:      d.NextRelease,
		UnitOfMeassure:       d.UnitOfMeasure,
		//QMI: , todo
	}

	mappedEditVersionMetaData := model.EditVersionMetaData{
		MetaData:   mappedMetaData,
		Collection: "", // TODO
		InstanceID: v.ID,
		Published:  v.State == "published", // TODO enum
	}

	//Alerts        *[]Alert            `json:"alerts"`
	//	CollectionID  string              `json:"collection_id"`
	//	Downloads     map[string]Download `json:"downloads"`
	//	Edition       string              `json:"edition"`
	//	Dimensions    []Dimension         `json:"dimensions"`
	//	ID            string              `json:"id"`
	//	InstanceID    string              `json:"instance_id"`
	//	LatestChanges []Change            `json:"latest_changes"`
	//	Links         Links               `json:"links"`
	//	ReleaseDate   string              `json:"release_date"`
	//	State         string              `json:"state"`
	//	Temporal      []Temporal          `json:"temporal"`
	//	Version       int                 `json:"version"`
	//  			edition: version.edition,
	//                version: version.version,
	//                releaseDate: { value: version.release_date || "", error: "" },
	//                notices: version.alerts ? this.mapNoticesToState(version.alerts, version.version || version.id) : [],
	//                dimensions: version.dimensions || [],
	//                usageNotes: version.usage_notes ? this.mapUsageNotesToState(version.usage_notes, version.version || version.id) : [],
	//                latestChanges: version.latest_changes ? this.mapLatestChangesToState(version.latest_changes, version.version || version.id) : []
	return mappedEditVersionMetaData
}

func mapRelatedContent() []model.RelatedContent{
	var latestChanges []model.RelatedContent
	for i, content := range v.LatestChanges {
		latestChanges = append(latestChanges, model.RelatedContent{
			ID:                    i,
			Title:                 "",
			Description:           "",
			Href:                  "",
			SimpleListHeading:     "",
			SimpleListDescription: "",
		})
	}
}