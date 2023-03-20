module github.com/ONSdigital/dp-publishing-dataset-controller

go 1.19

// fixes vulnerabilities in github.com/hashicorp/consul/api@v1.1.0 and github.com/hashicorp/consul/sdk@v0.1.1 dependencies
replace github.com/spf13/cobra => github.com/spf13/cobra v1.4.0

// fixes vulnerabilities in golang.org/x/text@v0.3.7 dependency
replace golang.org/x/text => golang.org/x/text v0.3.8

require (
	github.com/ONSdigital/dp-api-clients-go/v2 v2.234.1-0.20230309160605-57ab06778da7
	github.com/ONSdigital/dp-healthcheck v1.5.0
	github.com/ONSdigital/dp-net v1.5.0
	github.com/ONSdigital/log.go/v2 v2.3.0
	github.com/golang/mock v1.6.0
	github.com/gorilla/mux v1.8.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/smartystreets/goconvey v1.7.2
)

require (
	github.com/ONSdigital/dp-api-clients-go v1.43.0 // indirect
	github.com/ONSdigital/dp-net/v2 v2.8.0 // indirect
	github.com/aws/aws-sdk-go v1.44.204 // indirect
	github.com/fatih/color v1.14.1 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/hokaccha/go-prettyjson v0.0.0-20211117102719-0474bc63780f // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/justinas/alice v1.2.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/smartystreets/assertions v1.13.0 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
)
