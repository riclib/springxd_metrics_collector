// Copyright 2018 Solid Reason

package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const namespace = "springxd"
const xdHealthPath = "/management/health"
const xdMetricsPath = "/management/metrics"

var (
	addr        = flag.String("listen-address", ":9175", "The address to listen on for HTTP requests.")
	springxdURL = flag.String("springxd-url", "", "The springxd server url. (mandatory)")
)

type healthResponse struct {
	Status string
}

type metricsResponse struct {
	MemTotal                               int     `json:"mem"`
	MemFree                                int     `json:"mem.free"`
	Processors                             int     `json:"processors"`
	InstanceUptime                         int     `json:"instance.uptime"`
	Uptime                                 int     `json:"uptime"`
	SystemloadAverage                      float64 `json:"systemload.average"`
	HeapCommitted                          int     `json:"heap.committed"`
	HeapInit                               int     `json:"heap.init"`
	HeapUsed                               int     `json:"heap.used"`
	Heap                                   int     `json:"heap"`
	ThreadsPeak                            int     `json:"threads.peak"`
	ThreadsDaemon                          int     `json:"threads.daemon"`
	Threads                                int     `json:"threads"`
	Classes                                int     `json:"classes"`
	ClassesLoaded                          int     `json:"classes.loaded"`
	ClassesUnloaded                        int     `json:"classes.unloaded"`
	GcG1YoungGenerationCount               int     `json:"gc.g1_young_generation.count"`
	GcG1YoungGenerationTime                int     `json:"gc.g1_young_generation.time"`
	GcG1OldGenerationCount                 int     `json:"gc.g1_old_generation.count"`
	GcG1OldGenerationTime                  int     `json:"gc.g1_old_generation.time"`
	GaugeResponseManagementMetrics         float64 `json:"gauge.response.management.metrics"`
	GaugeResponseAdminUI                   float64 `json:"gauge.response.admin-ui"`
	GaugeResponseAuthenticate              float64 `json:"gauge.response.authenticate"`
	CounterStatus200JobsDefinitionsStar    int     `json:"counter.status.200.jobs.definitions.-star"`
	GaugeResponseStreamsDefinitionsStar    float64 `json:"gauge.response.streams.definitions.-star"`
	CounterStatus200ManagementMetrics      int     `json:"counter.status.200.management.metrics"`
	GaugeResponseManagementHealth          float64 `json:"gauge.response.management.health"`
	CounterStatus403Error                  int     `json:"counter.status.403.error"`
	CounterStatus200ManagementEnv          int     `json:"counter.status.200.management.env"`
	CounterStatus200ManagementHealth       int     `json:"counter.status.200.management.health"`
	CounterStatus200ManagementBeans        int     `json:"counter.status.200.management.beans"`
	CounterStatus200StarStarFaviconIco     int     `json:"counter.status.200.star-star.favicon.ico"`
	CounterStatus200SecurityInfo           int     `json:"counter.status.200.security.info"`
	GaugeResponseManagementBeans           float64 `json:"gauge.response.management.beans"`
	GaugeResponseStarStar                  float64 `json:"gauge.response.star-star"`
	CounterStatus200StreamsDefinitionsStar int     `json:"counter.status.200.streams.definitions.-star"`
	GaugeResponseError                     float64 `json:"gauge.response.error"`
	CounterStatus401Authenticate           int     `json:"counter.status.401.authenticate"`
	GaugeResponseSecurityInfo              float64 `json:"gauge.response.security.info"`
	GaugeResponseJobsDefinitionsStar       float64 `json:"gauge.response.jobs.definitions.-star"`
	CounterStatus200StarStar               int     `json:"counter.status.200.star-star"`
	CounterStatus302ManagementRoot         int     `json:"counter.status.302.management.root"`
	CounterStatus302AdminUI                int     `json:"counter.status.302.admin-ui"`
	GaugeResponseStarStarFaviconIco        float64 `json:"gauge.response.star-star.favicon.ico"`
	GaugeResponseManagementRoot            float64 `json:"gauge.response.management.root"`
	CounterStatus200Authenticate           int     `json:"counter.status.200.authenticate"`
	GaugeResponseManagementEnv             float64 `json:"gauge.response.management.env"`
	DatasourcePrimaryActive                int     `json:"datasource.primary.active"`
	DatasourcePrimaryUsage                 float64 `json:"datasource.primary.usage"`
}

var (
	// Create a summary to track fictional interservice RPC latencies for three
	// distinct services with different latency distributions. These services are
	// differentiated via a "service" label.

	up = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "springxd_up",
			Help: "SpringXD health status",
		},
	)

	memFree = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "springxd_mem_free",
			Help: "SpringXD free server memory",
		},
	)
	memTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "springxd_mem_total",
			Help: "SpringXD server memory",
		},
	)
	instanceUptime = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "springxd_instance_uptime",
			Help: "SpringXD Instance Uptime",
		},
	)
	uptime = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "springxd_uptime",
			Help: "SpringXD server uptime",
		},
	)
	heapCommitted = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "springxd_heap_committed",
			Help: "SpringXD Heap Committed",
		},
	)
	heapInit = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "springxd_heap_init",
			Help: "SpringXD Heap Init",
		},
	)
	heapUsed = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "springxd_heap_used",
			Help: "SpringXD HeapUsed",
		},
	)
	heap = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "springxd_heap",
			Help: "SpringXD Heap",
		},
	)

	threadsPeak = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "springxd_threads_peak",
			Help: "SpringXD ThreadsPeak",
		},
	)
	threadsDaemon = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "springxd_threads_daemon",
			Help: "SpringXD ThreadsDaemon",
		},
	)
	threads = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "springxd_threads",
			Help: "SpringXD server memory",
		},
	)
	gcG1YoungCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "springxd_gc_g1_young_generation_count",
			Help: "SpringXD gcG1YoungCount",
		},
	)
	gcG1YoungTime = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "springxd_gc_g1_young_generation_time",
			Help: "SpringXD gcG1YoungTime",
		},
	)
	gcG1OldCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "springxd_gc_g1_old_generation_count",
			Help: "SpringXD gcG1OldCount",
		},
	)
	gcG1OldTime = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "springxd_gc_g1_old_generation_time",
			Help: "SpringXD gcG1OldTime",
		},
	)
)

func init() {
	// Register the summary and the histogram with Prometheus's default registry.
	prometheus.MustRegister(up)
	prometheus.MustRegister(memFree)
	prometheus.MustRegister(memTotal)
	prometheus.MustRegister(instanceUptime)
	prometheus.MustRegister(uptime)
	prometheus.MustRegister(heapCommitted)
	prometheus.MustRegister(heapInit)
	prometheus.MustRegister(heapUsed)
	prometheus.MustRegister(heap)
	prometheus.MustRegister(threadsPeak)
	prometheus.MustRegister(threadsDaemon)
	prometheus.MustRegister(threads)
	prometheus.MustRegister(gcG1YoungCount)
	prometheus.MustRegister(gcG1YoungTime)
	prometheus.MustRegister(gcG1OldCount)
	prometheus.MustRegister(gcG1OldTime)

}

func scrapeHealth() {
	var r healthResponse

	resp, _ := http.Get(*springxdURL + xdHealthPath)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(body[:]))

	err := json.Unmarshal(body, &r)
	if err != nil {
		log.Println(err)
	}
	if r.Status == "UP" {
		up.Set(1)
	} else {
		up.Set(0)
	}
}
func scrapeMetrics() {
	var r metricsResponse

	resp, _ := http.Get(*springxdURL + xdMetricsPath)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(body[:]))

	err := json.Unmarshal(body, &r)
	if err != nil {
		log.Println(err)
	}
	memFree.Set(float64(r.MemFree))
	memTotal.Set(float64(r.MemTotal))
	instanceUptime.Set(float64(r.InstanceUptime))
	uptime.Set(float64(r.Uptime))
	heapCommitted.Set(float64(r.HeapCommitted))
	heapInit.Set(float64(r.HeapInit))
	heapUsed.Set(float64(r.HeapUsed))
	heap.Set(float64(r.Heap))
	threadsPeak.Set(float64(r.ThreadsPeak))
	threadsDaemon.Set(float64(r.ThreadsDaemon))
	threads.Set(float64(r.Threads))
	gcG1YoungCount.Set(float64(r.GcG1YoungGenerationCount))
	gcG1YoungTime.Set(float64(r.GcG1YoungGenerationTime))
	gcG1OldCount.Set(float64(r.GcG1OldGenerationCount))
	gcG1OldTime.Set(float64(r.GcG1OldGenerationTime))
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Delegate http serving to Prometheus client library, which will call collector.Collect.
	scrapeHealth()
	scrapeMetrics()
	h := promhttp.Handler()
	h.ServeHTTP(w, r)
}

func main() {
	flag.Parse()
	if *springxdURL == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	//	start := time.Now()

	memFree.Set(1)
	// Periodically record some sample latencies for the three services.

	// Expose the registered metrics via HTTP.
	http.HandleFunc("/metrics", handler)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
