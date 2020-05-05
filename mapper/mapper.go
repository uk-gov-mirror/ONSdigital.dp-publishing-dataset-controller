package mapper

import (
	"sort"
	"strings"

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
		return strings.ToLower(mappedDatasets[i].GetLabel()) < strings.ToLower(mappedDatasets[j].GetLabel())
	})

	return mappedDatasets
}
