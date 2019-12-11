module github.com/ONSdigital/dp-publishing-dataset-controller

go 1.12

replace github.com/ONSdigital/dp-api-clients-go => ../dp-api-clients-go/

require (
	github.com/ONSdigital/dp-api-clients-go v0.0.0-20191028105437-cc0f9e68e16a
	github.com/ONSdigital/go-ns v0.0.0-20191104121206-f144c4ec2e58
	github.com/ONSdigital/log.go v0.0.0-20191127134126-2a610b254f20
	github.com/davecgh/go-spew v1.1.1
	github.com/fatih/color v1.7.0 // indirect
	github.com/gorilla/mux v1.7.3
	github.com/hokaccha/go-prettyjson v0.0.0-20190818114111-108c894c2c0e // indirect
	github.com/justinas/alice v1.2.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553 // indirect
)
