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
		if mappedDatasets[i].Title == "" {
			return false
		} else if mappedDatasets[j].Title == "" {
			return true
		}
		return mappedDatasets[i].Title < mappedDatasets[j].Title
	})

	return mappedDatasets
}
