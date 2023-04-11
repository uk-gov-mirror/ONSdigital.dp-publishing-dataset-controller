package mapper

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"time"

	dataset "github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	zebedee "github.com/ONSdigital/dp-api-clients-go/v2/zebedee"
	babbageclient "github.com/ONSdigital/dp-publishing-dataset-controller/clients/topics"
	"github.com/ONSdigital/dp-publishing-dataset-controller/model"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/pkg/errors"
)

type related struct {
	publications  []model.RelatedContent
	methodologies []model.RelatedContent
	datasets      []model.RelatedContent
}

func AllDatasets(datasets dataset.List) []model.Dataset {
	var mappedDatasets []model.Dataset
	for _, ds := range datasets.Items {
		if &ds == nil || ds.Next == nil || ds.Next.Type == "nomis" {
			continue
		}
		mappedDatasets = append(mappedDatasets, model.Dataset{
			ID:    ds.ID,
			Title: ds.Next.Title,
		})
	}

	sort.Slice(mappedDatasets, func(i, j int) bool {
		return strings.ToLower(mappedDatasets[i].GetLabel()) < strings.ToLower(mappedDatasets[j].GetLabel())
	})

	return mappedDatasets
}

func AllVersions(ctx context.Context, dataset dataset.Dataset, edition dataset.Edition, versions dataset.VersionsList) model.VersionsPage {
	datasetName := dataset.Next.Title
	editionName := edition.Edition
	var mappedVersions []model.Version
	for _, v := range versions.Items {
		title := fmt.Sprintf("Version: %v", v.Version)
		if v.State == "published" {
			title += " (published)"
		}
		var timeF string
		time, err := time.Parse("2006-01-02T15:04:05Z", v.ReleaseDate)
		if err != nil {
			log.Warn(ctx, "failed to parse release date", log.FormatErrors([]error{err}))
		} else {
			timeF = time.Format("02 January 2006")
		}
		mappedVersions = append(mappedVersions, model.Version{
			ID:          v.ID,
			Title:       title,
			Version:     v.Version,
			ReleaseDate: timeF,
			State:       v.State,
		})
	}

	sort.Slice(mappedVersions, func(i, j int) bool {
		return mappedVersions[i].Version > mappedVersions[j].Version
	})

	return model.VersionsPage{
		DatasetName: datasetName,
		EditionName: editionName,
		Versions:    mappedVersions,
	}
}

func EditMetadata(d *dataset.DatasetDetails, v dataset.Version, dim []dataset.VersionDimension, c zebedee.Collection) model.EditMetadata {
	mappedMetadata := model.EditMetadata{
		Dataset:      *d,
		Version:      v,
		Dimensions:   dim,
		CollectionID: c.ID,
	}

	if len(c.Datasets) > 0 {
		for _, dataset := range c.Datasets {
			if dataset.ID == d.ID {
				mappedMetadata.CollectionState = dataset.State
				mappedMetadata.CollectionLastEditedBy = dataset.LastEditedBy
			}
		}
	}

	return mappedMetadata

}

// PutMetadata transform an EditMetadata object to the EditableMetadata as expected by dataset api
func PutMetadata(m model.EditMetadata) dataset.EditableMetadata {
	metadata := dataset.EditableMetadata{
		CanonicalTopic:    m.Dataset.CanonicalTopic,
		Description:       m.Dataset.Description,
		Dimensions:        m.Version.Dimensions,
		LatestChanges:     &m.Version.LatestChanges,
		License:           m.Dataset.License,
		NationalStatistic: &m.Dataset.NationalStatistic,
		NextRelease:       m.Dataset.NextRelease,
		QMI:               &m.Dataset.QMI,
		ReleaseDate:       m.Version.ReleaseDate,
		ReleaseFrequency:  m.Dataset.ReleaseFrequency,
		Subtopics:         m.Dataset.Subtopics,
		Survey:            m.Dataset.Survey,
		Title:             m.Dataset.Title,
		UnitOfMeasure:     m.Dataset.UnitOfMeasure,
	}

	if m.Dataset.Contacts != nil {
		metadata.Contacts = *m.Dataset.Contacts
	}
	if m.Dataset.Keywords != nil {
		metadata.Keywords = *m.Dataset.Keywords
	}
	if m.Dataset.Methodologies != nil {
		metadata.Methodologies = *m.Dataset.Methodologies
	}
	if m.Dataset.Publications != nil {
		metadata.Publications = *m.Dataset.Publications
	}
	if m.Dataset.RelatedDatasets != nil {
		metadata.RelatedDatasets = *m.Dataset.RelatedDatasets
	}
	if m.Dataset.RelatedContent != nil {
		metadata.RelatedContent = *m.Dataset.RelatedContent
	}

	if m.Version.Alerts != nil {
		metadata.Alerts = m.Version.Alerts
	}
	if m.Version.UsageNotes != nil {
		metadata.UsageNotes = m.Version.UsageNotes
	}

	return metadata
}

func EditDatasetVersionMetaData(d dataset.DatasetDetails, v dataset.Version) (model.EditVersionMetaData, error) {
	keywordsString := ""
	if d.Keywords != nil {
		keywords := *d.Keywords
		keywordsString = fmt.Sprint(strings.Join(keywords, ", "))
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
		UsageNotes:    mapUsageNotes(v.UsageNotes),
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

// Topics takes babbage topics respond and returns a slice with topic titles
func Topics(tpcs babbageclient.TopicsResult) []model.Topics {
	var topics []model.Topics
	if len(tpcs.Topics.Results) > 0 {
		for _, tpc := range tpcs.Topics.Results {
			topics = append(topics, model.Topics{
				Title: tpc.Description.Title,
			})
		}
	}
	return topics
}
