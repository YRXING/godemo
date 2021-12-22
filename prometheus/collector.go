package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"math/rand"
	"sync"
)

// customized collector, need to realize Describe and Collect interface
type Metrics struct {
	// prometheus.Desc is a descriptor of some metric.
	metrics map[string]*prometheus.Desc
	mutex 	sync.Mutex
}

// create metric descriptor
func newGlobalMetric(namespace string, metricName string, docString string, labels []string) *prometheus.Desc {
	return prometheus.NewDesc(namespace+"_"+metricName,docString,labels,nil)
}

// initial Metrics struct
func NewMetrics(namespace string) *Metrics {
	return &Metrics{
		metrics: map[string]*prometheus.Desc{
			"my_counter_metric": newGlobalMetric(namespace, "my_counter_metric", "The description of my_counter_metric", []string{"host"}),
			"my_gauge_metric": newGlobalMetric(namespace,"my_gauge_metric", "The description of my_gauge_metric",[]string{"host"}),
		},
	}
}

// used to display the describe information of metric
func (m *Metrics) Describe(ch chan<- *prometheus.Desc)  {
	for _, v := range m.metrics {
		ch <- v
	}
}

// convert data collected to the form that prometheus server recognize
func (m *Metrics) Collect(ch chan<- prometheus.Metric)  {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	mockCounterMetricData, mockGaugeMetricData := m.GenerateMockData()
	for host,currentValue := range mockCounterMetricData {
		ch <- prometheus.MustNewConstMetric(m.metrics["my_counter_metric"],prometheus.CounterValue, float64(currentValue),host)
	}

	for host,currentValue := range mockGaugeMetricData {
		ch <- prometheus.MustNewConstMetric(m.metrics["my_gauge_metric"],prometheus.GaugeValue, float64(currentValue),host)
	}
}


// generate mock data
func (m *Metrics) GenerateMockData() (mockCounterMetricData map[string]int, mockGaugeMetricData map[string]int)  {
	mockCounterMetricData = map[string]int{
		"yahoo.com": int(rand.Int31n(1000)),
		"google.com": int(rand.Int31n(1000)),
	}
	mockGaugeMetricData = map[string]int{
		"yahoo.com": int(rand.Int31n(10)),
		"google.com": int(rand.Int31n(10)),
	}
	return
}
