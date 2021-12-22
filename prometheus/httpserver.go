package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)


// use customized collector
func CollectorServer() {
	collector := NewMetrics("system")
	registry := prometheus.NewRegistry()
	registry.MustRegister(collector)

	http.Handle("/metrics", promhttp.HandlerFor(registry,promhttp.HandlerOpts{}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>A Prometheus Exporter</title></head>
			<body>
			<h1>A Prometheus Exporter</h1>
			<p><a href='/metrics'>Metrics</a></p>
			</body>
			</html>`))
	})
	http.ListenAndServe(":9001",nil)
}

func DirectServer()  {
	// either a collector or a Metric, the value can be set itself
	cpuTemp := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "cpu_temperature_celsius",
			Help: "Current temperature of the CPU",
		})

	cpuTemp2 := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cpu_temperature_celsius_ppp",
			Help: "Current temperature of the CPU.",
		},
		[]string{
			"endpoint",
			"component",
		})

	hdFailures := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hd_errors_total",
			Help: "Number of hard-disk errors.",
		},
		[]string{"device"})

	// register to default registry
	prometheus.MustRegister(cpuTemp,cpuTemp2)
	prometheus.MustRegister(hdFailures)

	cpuTemp.Set(65.3)
	cpuTemp2.WithLabelValues("root","mmm").Add(3)
	cpuTemp2.WithLabelValues("user","aaa").Add(4)
	hdFailures.With(prometheus.Labels{"device":"/dev/sda"}).Inc()

	http.Handle("/metrics",promhttp.Handler())
	http.ListenAndServe(":9002",nil)

}
