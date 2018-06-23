// Copyright 2018 Solid Reason

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const namespace = "springxd"
const xdHealthPath = "/management/health"
const xdMetricsPath = "/management/metrics"
const xdJobsPath = "/jobs/executions"

var (
	port        = flag.String("listen-port", "9175", "The port to listen on for HTTP requests.")
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

func scrapeJobExecutionMetrics(basename string, i interface{}) string {
	var res string
	return res
}

func collectJobDefs(base string, u string) string {
	var j interface{}
	var err error

	client := &http.Client{}
	req := &http.Request{Method: "GET"}
	req.URL, err = url.Parse(u)
	if err != nil {
		log.Fatalln("couldn't create request ")
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error retrieving " + u)
		return ""
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &j)
	if err != nil {
		log.Println(err)
	}
	log.Println("Retrieved job executions")
	return scrapeJobExecutionMetrics(base, j)

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
	s = s + collectJobDefs(namespace+"_jobs", *springxdURL+xdJobsPath)
	end := time.Now()
	s = s + fmt.Sprintf(namespace+"_scrape_duration_seconds %f\n", end.Sub(start).Seconds())
	w.Write([]byte(s))
}

func main() {
	flag.Parse()
	if *springxdURL == "" {
		log.Println("Test mode, retrieving metrics from json/ folder")
		*springxdURL = "http://localhost:9175"
	}
	//	start := time.Now()

	// Periodically record some sample latencies for the three services.

	// Expose the registered metrics via HTTP.
	http.HandleFunc("/metrics", handler)

	// default test endpoints
	http.HandleFunc(xdJobsPath,
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "json/jobs_executions.json")
		})
	http.HandleFunc(xdMetricsPath,
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "json/management_metrics.json")
		})
	http.HandleFunc(xdHealthPath,
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "json/management_health.json")
		})

	log.Printf("Listening for requests on localhost:" + *port)
	err := http.ListenAndServe(":"+*port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
