package mapper

import (
	"sort"

	"github.com/ONSdigital/dp-api-clients-go/dataset"
	"github.com/ONSdigital/dp-publishing-dataset-controller/model"
)

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
		return mappedDatasets[i].Title < mappedDatasets[j].Title
	})

	return mappedDatasets
}
