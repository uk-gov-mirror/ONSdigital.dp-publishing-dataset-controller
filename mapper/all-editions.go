package mapper

import (
	"context"
	"time"

	dataset "github.com/ONSdigital/dp-api-clients-go/dataset"
	"github.com/ONSdigital/dp-publishing-dataset-controller/model"
	"github.com/ONSdigital/log.go/log"
)

// AllEditions maps dataset and editions response to editions list page model
func AllEditions(ctx context.Context, dataset dataset.Dataset, editions []dataset.Edition, latestVersions map[string]string) model.EditionsPage {
	var mappedEditions []model.Edition
	for _, e := range editions {
		var timeF string
		for k, latestVersion := range latestVersions {
			if k == e.Edition {
				time, err := time.Parse("2006-01-02T15:04:05Z", latestVersion)
				if err != nil {
					log.Event(ctx, "failed to parse release date", log.WARN, log.Error(err))
				} else {
					timeF = time.Format("02 January 2006")
				}
			}
		}
		mappedEditions = append(mappedEditions, model.Edition{
			ID:          e.Edition,
			Title:       e.Edition,
			ReleaseDate: timeF,
		})
	}

	return model.EditionsPage{
		DatasetName: dataset.Next.Title,
		Editions:    mappedEditions,
	}
}
