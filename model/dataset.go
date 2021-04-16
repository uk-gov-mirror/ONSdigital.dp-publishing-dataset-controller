package model

import (
	"github.com/ONSdigital/dp-api-clients-go/dataset"
	datasetclient "github.com/ONSdigital/dp-api-clients-go/dataset"
)

type Dataset struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type Version struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Version     int    `json:"version"`
	ReleaseDate string `json:"release_date"`
}

type EditMetadata struct {
	Dataset                datasetclient.DatasetDetails     `json:"dataset"`
	Version                datasetclient.Version            `json:"version"`
	Dimensions             []datasetclient.VersionDimension `json:"dimensions"`
	CollectionID           string                           `json:"collection_id"`
	CollectionState        string                           `json:"collection_state"`
	CollectionLastEditedBy string                           `json:"collection_last_edited_by"`
}

type EditVersionMetaData struct {
	MetaData   MetaData `json:"meta_data"`
	Collection string   `json:"collection"`
	InstanceID string   `json:"instance_id"`
	Published  bool     `json:"published"`
}

type MetaData struct {
	Edition       string                     `json:"edition"`
	Version       int                        `json:"version"`
	ReleaseDate   ReleaseDate                `json:"release-date"`
	Notices       []Notice                   `json:"notices"`
	Dimensions    []dataset.VersionDimension `json:"dimensions"`
	UsageNotes    []UsageNote                `json:"usage_notes"`
	LatestChanges []LatestChanges            `json:"latest_changes"`

	Title                string           `json:"title"`
	Summary              string           `json:"summary"`
	Keywords             string           `json:"keywords"`
	NationalStatistic    bool             `json:"national_statistic"`
	License              string           `json:"license"`
	ContactName          string           `json:"contact_name"`
	ContactEmail         string           `json:"contact_email"`
	ContactTelephone     string           `json:"contact_telephone"`
	RelatedDatasets      []RelatedContent `json:"related_datasets"`
	RelatedPublications  []RelatedContent `json:"related_publications"`
	RelatedMethodologies []RelatedContent `json:"related_methodologies"`
	ReleaseFrequency     string           `json:"release_frequency"`
	NextReleaseDate      string           `json:"next_release_date"`
	UnitOfMeassure       string           `json:"unit_of_meassure"`
	QMI                  string           `json:"qmi"`
}

type ReleaseDate struct {
	ReleaseDate string `json:"release_date"`
	Error       string `json:"error"`
}

type Notice struct {
	ID                    int    `json:"id"`
	Type                  string `json:"type"`
	Date                  string `json:"date"`
	Description           string `json:"description"`
	SimpleListHeading     string `json:"simple_list_heading"`
	SimpleListDescription string `json:"simple_list_description"`
}

type UsageNote struct {
	ID                    int    `json:"id"`
	Title                 string `json:"title"`
	Note                  string `json:"note"`
	SimpleListHeading     string `json:"simple_list_heading"`
	SimpleListDescription string `json:"simple_list_description"`
}

type LatestChanges struct {
	ID                    int    `json:"id"`
	Title                 string `json:"title"`
	Description           string `json:"description"`
	SimpleListHeading     string `json:"simple_list_heading"`
	SimpleListDescription string `json:"simple_list_description"`
}

type RelatedContent struct {
	ID                    int    `json:"id"`
	Title                 string `json:"title"`
	Description           string `json:"description"`
	Href                  string `json:"href"`
	SimpleListHeading     string `json:"simple_list_heading"`
	SimpleListDescription string `json:"simple_list_description"`
}

// GetLabel will return the dataset's name. If the dataset does not have a title, it will instead return the ID
func (d Dataset) GetLabel() string {
	if d.Title == "" {
		return d.ID
	}
	return d.Title
}

type Topics struct {
	Title string `json:"title"`
}
