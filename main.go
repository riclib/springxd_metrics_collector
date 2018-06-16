// Copyright 2018 Solid Reason

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const namespace = "springxd"
const xdHealthPath = "/management/health"
const xdMetricsPath = "/management/metrics"

var (
	addr        = flag.String("listen-address", ":9175", "The address to listen on for HTTP requests.")
	springxdURL = flag.String("springxd-url", "", "The springxd server url. (mandatory)")
)

func scrapeMetrics(basename string, i interface{}) string {
	var res string
	//	fmt.Println(i)
	m := i.(map[string]interface{})
	for k, v := range m {
		metricName := strings.Replace(k, ".", "_", -1)
		metricName = strings.Replace(metricName, "-", "", -1)

		switch vv := v.(type) {
		case string:
			switch v.(string) {
			case "Up", "UP":
				res = res + basename + "_" + metricName + " 1.0" + "\n"
			case "Down", "DOWN":
				res = res + basename + "_" + metricName + " 0.0" + "\n"
			default:
			}
		case float64:
			res = res + basename + "_" + metricName + " " + fmt.Sprintf("%f", v.(float64)) + "\n"
		case []interface{}:
			res = res + scrapeMetrics(basename+"_"+metricName, v)
		case map[string]interface{}:
			res = res + scrapeMetrics(basename+"_"+metricName, v)
		default:
			log.Println(metricName, "is of a type I don't know how to handle:", vv)
		}
	}
	return res
}

func collect(base string, url string) string {
	var j interface{}
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err := json.Unmarshal(body, &j)
	if err != nil {
		log.Println(err)
	}

	return scrapeMetrics(base, j)

}
func handler(w http.ResponseWriter, r *http.Request) {
	// Delegate http serving to Prometheus client library, which will call collector.Collect.
	start := time.Now()
	s := collect(namespace+"_health", *springxdURL+xdHealthPath)
	s = s + collect(namespace+"_metrics", *springxdURL+xdMetricsPath)
	end := time.Now()
	s = s + fmt.Sprintf(namespace+"_scrape_duration_seconds %f\n", end.Sub(start).Seconds())
	w.Write([]byte(s))
}

func main() {
	flag.Parse()
	if *springxdURL == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	//	start := time.Now()

	// Periodically record some sample latencies for the three services.

	// Expose the registered metrics via HTTP.
	http.HandleFunc("/metrics", handler)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
