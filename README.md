dp-publishing-dataset-controller
================

Controller to coordinate all requests between frontend CMS and APIs involved in dataset upload, creation and editing. 

### Getting started

To run this service you must have [Golang](https://golang.org/) installed on a UNIX machine.

Once you have installed those dependencies and cloned this repo you need to run the following:

1. Move into the correct directory
```
cd dp-publishing-dataset-controller
```
2. Run the service
```
make debug
```


### Configuration

| Environment variable           | Default                           | Description
| ------------------------------ | --------------------------------- | -----------
| BIND_ADDR                      | :24000                            | The host and port to bind to
| DATASET_API_URL                | http://localhost:22000            | The host name for [Dataset API](https://github.com/ONSdigital/dp-dataset-api)
| ZEBEDEE_URL                    | http://localhost:8081             | The host name for [Zebedee](https://github.com/ONSdigital/zebedee)
| GRACEFUL_SHUTDOWN_TIMEOUT      | 5s                                | The graceful shutdown timeout in seconds
| HEALTHCHECK_INTERVAL           | 30s                               | Healthcheck interval in seconds
| HEALTHCHECK_CRITICAL_TIMEOUT   | 90s                               | Healthcheck timeout in seconds


### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2020, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
