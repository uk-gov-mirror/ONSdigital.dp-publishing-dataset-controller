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
| API_ROUTER_URL                 | http://localhost:23200/v1         | The URL of the [dp-api-router](https://github.com/ONSdigital/dp-api-router)
| BABBAGE_URL                    | http://localhost:8080             | The URL for [Babbage](https://github.com/ONSdigital/babbage)
| DATASET_BATCH_SIZE             | 100                               | Size of the batches, used for pagination
| DATASET_BATCH_WORKERS          | 10                                | Number of batch workers, used for pagination
| GRACEFUL_SHUTDOWN_TIMEOUT      | 5s                                | The graceful shutdown timeout in seconds
| HEALTHCHECK_INTERVAL           | 30s                               | Healthcheck interval in seconds
| HEALTHCHECK_CRITICAL_TIMEOUT   | 90s                               | Healthcheck timeout in seconds


### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2021, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
