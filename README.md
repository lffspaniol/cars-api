# Car API

## Doc

The open api spec is on [file](./openapi.yaml)

if you prefer postman collection is on [file](cars-api.postman_collection.json)

## run

you can run with VScode task build or type `docker compose up`.

### Runing test

`
    go test ./... -race
`

### Integration Testt

install the [newman](https://learning.postman.com/docs/collections/using-newman-cli/installing-running-newman/)

`docker compose up -d `

`newman run cars-api.postman_collection.json`

warning you node need to be LTS to run integration test
you can run it too on github action

## Development guide

### Lint

The project utilizes the [golangci-lint](https://golangci-lint.run/) as a lint runner [file](.golangci.yaml)

#### Recommended

for better CI acceptance is recommended to run the [Integrations](https://golangci-lint.run/usage/integrations/)


The demo exposes the following backends:

- Jaeger at http://0.0.0.0:16686
- Zipkin at http://0.0.0.0:9411
- Prometheus at http://0.0.0.0:9090 
