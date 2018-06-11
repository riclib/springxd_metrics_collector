# springxd_metrics_collector

## Metric collector for SpringXD
Exposes an HTTP server that converts the REST metrics from springxd as prometheus metrics

## Usage:
```
springxd_metrics_collector  -listen-address string -springxd-url string
  -listen-address string
    	The address to listen on for HTTP requests. (default ":8080")
  -springxd-url string
    	The springxd server url.
```
## Example
`springxd_metrics_collector.py -listen-address :1934 -springxd-url http://demo6819977.mockable.io`

## No extra components required to install
