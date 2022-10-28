package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

func gettslstatus() (tt float64) {
	result, _ := exec.Command("bash", "-c", "cat /data/inphase/tools/podcheckerror.txt|wc -l").Output()
	fmt.Printf("result1", result)
	a := string(result)
	fmt.Printf("a", a)
	b, _ := strconv.Atoi(strings.TrimSpace(a))
	fmt.Printf("b", b)
	c := float64(b)
	return c

}
func getyyhstatus() (ts float64) {
	result, _ := exec.Command("bash", "-c", "cat /data/inphase/tools/podcheck.txt|wc -l").Output()
	fmt.Printf("result1", result)
	a := string(result)
	fmt.Printf("a", a)
	b, _ := strconv.Atoi(strings.TrimSpace(a))
	fmt.Printf("b", b)
	c := float64(b)
	return c

}

type memCollect struct {
	memMetric *prometheus.Desc
}

func (collector *memCollect) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.memMetric

}

var constLabel = prometheus.Labels{"comment": "outofmem"}
var Vlabel = []string{"status"}

var Vvalue1 = "true"
var Vvalue2 = "flase"

func newmemcollector() *memCollect {
	return &memCollect{memMetric: prometheus.NewDesc("mem_outof_status", "count memoutof pod", Vlabel, constLabel)}
}

func (collector *memCollect) Collect(ch chan<- prometheus.Metric) {
	var metricvalue1 float64
	var metricvalue2 float64
	metricvalue1 = gettslstatus()
	metricvalue2 = getyyhstatus()
	ch <- prometheus.MustNewConstMetric(collector.memMetric, prometheus.CounterValue, metricvalue1, Vvalue2)
	ch <- prometheus.MustNewConstMetric(collector.memMetric, prometheus.CounterValue, metricvalue2, Vvalue1)
}

func main() {
	prometheus.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts(prometheus.ProcessCollectorOpts{})))
	prometheus.Unregister(collectors.NewBuildInfoCollector())
	prometheus.Unregister(collectors.NewGoCollector())
	mem := newmemcollector()
	prometheus.MustRegister(mem)
	http.Handle("/metrics", promhttp.Handler())
	logrus.Info("begin to server on port 45679")
	logrus.Fatal(http.ListenAndServe(":45679", nil))
}
