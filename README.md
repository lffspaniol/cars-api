# Car API

## run

you can run with VScode task build or type `docker compose up`.

## Development guide

### Lint

The project utilizes the [golangci-lint](https://golangci-lint.run/) as a lint runner [file](.golangci.yaml)

#### Recommended

for better CI acceptance is recommended to run the [Integrations](https://golangci-lint.run/usage/integrations/)


The demo exposes the following backends:

- Jaeger at http://0.0.0.0:16686
- Zipkin at http://0.0.0.0:9411
- Prometheus at http://0.0.0.0:9090 
