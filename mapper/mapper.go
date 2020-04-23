package mapper

import (
	"fmt"
	"github.com/ONSdigital/dp-api-clients-go/dataset"
	"github.com/ONSdigital/dp-publishing-dataset-controller/model"
	"sort"
	"strings"
	"time"
)

type related struct {
	publications  []model.RelatedContent
	methodologies []model.RelatedContent
	datasets      []model.RelatedContent
}

func AllDatasets(datasets dataset.List) []model.Dataset {
	var mappedDatasets []model.Dataset
	for _, ds := range datasets.Items {
		if &ds == nil || ds.Next == nil {
			continue
		}
		mappedDatasets = append(mappedDatasets, model.Dataset{
			ID:    ds.ID,
			Title: ds.Next.Title,
		})
	}

	sort.Slice(mappedDatasets, func(i, j int) bool {
		if mappedDatasets[i].Title == "" {
			return false
		} else if mappedDatasets[j].Title == "" {
			return true
		}
		return mappedDatasets[i].Title < mappedDatasets[j].Title
	})

	return mappedDatasets
}

func EditDatasetVersionMetaData(d dataset.DatasetDetails, v dataset.Version) model.EditVersionMetaData {

	keywords := *d.Keywords
	keywordsString := fmt.Sprintf(strings.Join(keywords[:], ", "))
	relatedContent := mapRelatedContent(*d.RelatedDatasets, *d.Methodologies, *d.Publications)
	contacts := *d.Contacts

	releaseDate := model.ReleaseDate{
		ReleaseDate: v.ReleaseDate,
		Error:       "",
	}

	mappedMetaData := model.MetaData{
		Edition:       v.Edition,
		Version:       v.Version,
		ReleaseDate:   releaseDate,
		Notices:       mapAlerts(v),
		Dimensions:    v.Dimensions,
		UsageNotes:    mapUsageNotes(*d.UsageNotes),
		LatestChanges: mapLatestChanges(v.LatestChanges),

		Title:                d.Title,
		Summary:              d.Description,
		Keywords:             keywordsString,
		NationalStatistic:    d.NationalStatistic,
		License:              d.License,
		ContactName:          contacts[0].Name,
		ContactEmail:         contacts[0].Email,
		ContactTelephone:     contacts[0].Telephone,
		RelatedDatasets:      relatedContent.datasets,
		RelatedPublications:  relatedContent.publications,
		RelatedMethodologies: relatedContent.methodologies,
		ReleaseFrequency:     d.ReleaseFrequency,
		NextReleaseDate:      d.NextRelease,
		UnitOfMeassure:       d.UnitOfMeasure,
		QMI:                  d.QMI.URL,
	}
	var mappedCollectionValue string
	if v.CollectionID == "" {
		mappedCollectionValue = "false" // todo check it is a string false
	} else {
		mappedCollectionValue = v.CollectionID
	}
	mappedEditVersionMetaData := model.EditVersionMetaData{
		MetaData:   mappedMetaData,
		Collection: mappedCollectionValue,
		InstanceID: v.ID,
		Published:  v.State == "published",
	}

	return mappedEditVersionMetaData
}

// TODO make DRY - cast to a common type maybe
func mapRelatedContent(rd []dataset.RelatedDataset, rm []dataset.Methodology, rp []dataset.Publication) related {
	var relatedContent related

	for i, content := range rd {
		relatedContent.methodologies = append(relatedContent.methodologies, model.RelatedContent{
			ID:                    i,
			Title:                 content.Title,
			//Description:           content.Description, // TODO is it always empty?
			Href:                  content.URL,
			SimpleListHeading:     content.Title,
			//SimpleListDescription: content.Description, // TODO is it always empty?
		})
	}

	for i, content := range rm {
		relatedContent.methodologies = append(relatedContent.methodologies, model.RelatedContent{
			ID:                    i,
			Title:                 content.Title,
			Description:           content.Description,
			Href:                  content.URL,
			SimpleListHeading:     content.Title,
			SimpleListDescription: content.Description,
		})
	}

	for i, content := range rp {
		relatedContent.methodologies = append(relatedContent.methodologies, model.RelatedContent{
			ID:                    i,
			Title:                 content.Title,
			Description:           content.Description,
			Href:                  content.URL,
			SimpleListHeading:     content.Title,
			SimpleListDescription: content.Description,
		})
	}
	return relatedContent
}

func mapAlerts(v dataset.Version) []model.Notice {
	layout := "2006-01-02T15:04:05.000Z" // TODO check this is the right layout for the dates
	var notices []model.Notice

	for i, alert := range *v.Alerts {
		alertDateInDateFormat, err := time.Parse(layout, alert.Date)
		if err != nil {
			//TODO log error with the version that failed and return
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
	return notices
}

func mapUsageNotes(un []dataset.UsageNote) []model.UsageNote {
	var usageNotes []model.UsageNote
	for i, note := range un {
		usageNotes = append(usageNotes, model.UsageNote{
			ID:                    i,
			Title:                 note.Title,
			Note:                  note.Note,
			SimpleListHeading:     note.Title,
			SimpleListDescription: note.Note,
		})
	}
	return usageNotes
}

func mapLatestChanges(un []dataset.Change) []model.LatestChanges{
	var latestChanges []model.LatestChanges

	for i, change := range un {
		latestChanges = append(latestChanges, model.LatestChanges{
			ID:                    i,
			Title:                 change.Name,
			Description:           change.Description,
			SimpleListHeading:     change.Name,
			SimpleListDescription: change.Description,
		})
	}
	return latestChanges
}
