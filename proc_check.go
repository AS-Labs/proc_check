package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/process"
)

// Define metric descriptors
var (
	processExists = prometheus.NewDesc(
		"process_exists",
		"Whether the process exists (1 = exists, 0 = does not)",
		[]string{"pid", "name"},
		nil,
	)
	processCPU = prometheus.NewDesc(
		"process_cpu_usage",
		"CPU usage percentage of the process",
		[]string{"pid", "name"},
		nil,
	)
	processMemory = prometheus.NewDesc(
		"process_memory_usage_bytes",
		"Memory usage of the process in bytes",
		[]string{"pid", "name"},
		nil,
	)
	processArg = prometheus.NewDesc(
		"process_arg",
		"Command-line arguments of the process",
		[]string{"pid", "name", "index", "value"},
		nil,
	)
)

// ProcessCollector implements the Prometheus Collector interface
type ProcessCollector struct {
	processName string
}

// Describe sends metric descriptors to the channel
func (c *ProcessCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- processExists
	ch <- processCPU
	ch <- processMemory
	ch <- processArg
}

// Collect gathers metrics for processes matching the given name
func (c *ProcessCollector) Collect(ch chan<- prometheus.Metric) {
	processes, err := process.Processes()
	if err != nil {
		log.Printf("Error retrieving processes: %v", err)
		return
	}

	processFound := false
	for _, p := range processes {
		cmdline, err := p.Cmdline()
		if err != nil || !strings.Contains(cmdline, c.processName) {
			continue
		}
		processFound = true
		pid := fmt.Sprintf("%d", p.Pid)
		name, _ := p.Name()

		// Metric: Process Existence
		ch <- prometheus.MustNewConstMetric(
			processExists,
			prometheus.GaugeValue,
			1, // 1 indicates the process exists
			pid, name,
		)

		// Metric: CPU Usage
		cpuPercent, err := p.CPUPercent()
		if err == nil {
			ch <- prometheus.MustNewConstMetric(
				processCPU,
				prometheus.GaugeValue,
				cpuPercent,
				pid, name,
			)
		}

		// Metric: Memory Usage
		memInfo, err := p.MemoryInfo()
		if err == nil {
			ch <- prometheus.MustNewConstMetric(
				processMemory,
				prometheus.GaugeValue,
				float64(memInfo.RSS),
				pid, name,
			)
		}

		// Metric: Command-Line Arguments
		args, err := p.CmdlineSlice()
		if err == nil {
			for i, arg := range args {
				ch <- prometheus.MustNewConstMetric(
					processArg,
					prometheus.GaugeValue,
					1, // Value is 1, actual argument is in the "value" label
					pid, name, fmt.Sprintf("%d", i), arg,
				)
			}
		}
	}

	if !processFound {
		ch <- prometheus.MustNewConstMetric(
			processExists,
			prometheus.GaugeValue,
			0, // 0 indicates no process found
			"", "",
		)
	}
}

func main() {
	// Parse command-line arguments
	processName := flag.String("process", "", "Name of the process to monitor")
	flag.Parse()

	if *processName == "" {
		fmt.Println("Usage: ./exporter -process <process_name>")
		os.Exit(1)
	}

	// Create and register the collector
	collector := &ProcessCollector{processName: *processName}
	prometheus.MustRegister(collector)

	// Set up HTTP handler for Prometheus metrics
	http.Handle("/metrics", promhttp.Handler())
	fmt.Printf("Starting exporter on :8081/metrics for process '%s'\n", *processName)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
