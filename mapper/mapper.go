package mapper

import (
	"fmt"
	"github.com/ONSdigital/dp-api-clients-go/dataset"
	"github.com/ONSdigital/dp-publishing-dataset-controller/model"
	"github.com/pkg/errors"
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
		return mappedDatasets[i].GetTitle() < mappedDatasets[j].GetTitle()
	})

	return mappedDatasets
}

func EditDatasetVersionMetaData(d dataset.DatasetDetails, v dataset.Version) (model.EditVersionMetaData, error) {
	keywordsString := ""
	if d.Keywords != nil {
		keywords := *d.Keywords
		keywordsString = fmt.Sprintf(strings.Join(keywords[:], ", "))
	}

	relatedContent := mapRelatedContent(d.RelatedDatasets, d.Methodologies, d.Publications)

	var contacts = []dataset.Contact{
		{
			Name:      "",
			Telephone: "",
			Email:     "",
		},
	}
	if d.Contacts != nil {
		contacts = *d.Contacts
	}

	releaseDate := model.ReleaseDate{
		ReleaseDate: v.ReleaseDate,
		Error:       "",
	}

	notices, err := mapAlerts(v)
	if err != nil {
		return model.EditVersionMetaData{}, errors.Wrap(err, "error whilst parsing alerts")
	}

		mappedMetaData := model.MetaData{
		Edition:       v.Edition,
		Version:       v.Version,
		ReleaseDate:   releaseDate,
		Notices:       notices,
		Dimensions:    v.Dimensions,
		UsageNotes:    mapUsageNotes(d.UsageNotes),
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
		mappedCollectionValue = "false"
	} else {
		mappedCollectionValue = v.CollectionID
	}
	mappedEditVersionMetaData := model.EditVersionMetaData{
		MetaData:   mappedMetaData,
		Collection: mappedCollectionValue,
		InstanceID: v.ID,
		Published:  v.State == "published",
	}

	return mappedEditVersionMetaData, nil
}

func mapRelatedContent(rd *[]dataset.RelatedDataset, rm *[]dataset.Methodology, rp *[]dataset.Publication) related {
	var relatedContent related
	if rd != nil {
		for i, content := range *rd {
			relatedContent.datasets = append(relatedContent.datasets, model.RelatedContent{
				ID:                i,
				Title:             content.Title,
				Href:              content.URL,
				SimpleListHeading: content.Title,
			})
		}

	}

	if rm != nil {
		for i, content := range *rm {
			relatedContent.methodologies = append(relatedContent.methodologies, model.RelatedContent{
				ID:                    i,
				Title:                 content.Title,
				Description:           content.Description,
				Href:                  content.URL,
				SimpleListHeading:     content.Title,
				SimpleListDescription: content.Description,
			})
		}
	}

	if rp != nil {
		for i, content := range *rp {
			relatedContent.publications = append(relatedContent.publications, model.RelatedContent{
				ID:                    i,
				Title:                 content.Title,
				Description:           content.Description,
				Href:                  content.URL,
				SimpleListHeading:     content.Title,
				SimpleListDescription: content.Description,
			})
		}
	}
	return relatedContent
}

func mapAlerts(v dataset.Version) ([]model.Notice, error) {
	//layout := "2006-01-02T15:04:05.000Z"
	var notices []model.Notice

	if v.Alerts == nil {
		return notices, nil
	}
	for i, alert := range *v.Alerts {
		alertDateInDateFormat, err := time.Parse(time.RFC3339Nano, alert.Date)
		if err != nil {
			return nil, errors.Wrap(err, "error whilst parsing time from alert date")
		}

		noticeDate := alertDateInDateFormat.Format("02 Jan 2006")
		simpleListHeading := fmt.Sprintf(`%s (%s)`, alert.Type, noticeDate)
		notices = append(notices, model.Notice{
			ID:                    i,
			Type:                  alert.Type,
			Date:                  noticeDate,
			Description:           alert.Description,
			SimpleListHeading:     simpleListHeading,
			SimpleListDescription: alert.Description,
		})
	}

	return notices, nil
}

func mapUsageNotes(un *[]dataset.UsageNote) []model.UsageNote {
	var usageNotes []model.UsageNote
	if un == nil {
		return usageNotes
	}

	for i, note := range *un {
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

func mapLatestChanges(un []dataset.Change) []model.LatestChanges {
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
