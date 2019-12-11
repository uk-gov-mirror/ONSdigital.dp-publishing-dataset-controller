package mapper

import (
	"github.com/ONSdigital/dp-api-clients-go/dataset"
	"github.com/ONSdigital/dp-publishing-dataset-controller/model"
)

func AllDatasets(datasets dataset.ModelCollection) []model.Dataset {
	var mappedDatasets []model.Dataset
	for _, ds := range datasets.Items {
		mappedDatasets = append(mappedDatasets, model.Dataset{
			ID:           ds.ID,
			Title:        ds.Title,
			CollectionID: ds.CollectionID,
		})
	}
	return mappedDatasets
}
