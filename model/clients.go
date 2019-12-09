package model

import (
	"github.com/ONSdigital/dp-api-clients-go/dataset"
	zebedee "github.com/ONSdigital/dp-api-clients-go/zebedee"
)

type Clients struct {
	Dc *dataset.Client
	Zc *zebedee.ZebedeeClient
}
