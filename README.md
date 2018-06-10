# springxd_metrics_collector

## Metric collector for SpringXD
Exposes an HTTP server that converts the REST metrics from springxd as prometheus metrics/

## Usage:
`python springxd_metrics_collector.py <port> <server_url>`

## Example
`python springxd_metrics_collector.py 1934 http://demo6819977.mockable.io`

## Requires installing the following to run
`pip install prometheus_client requests`
