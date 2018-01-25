package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/baa.v1"
)

var (
	apiResponseTime *prometheus.SummaryVec
	apiResponseSize *prometheus.SummaryVec
	apiRequestTotal *prometheus.CounterVec
	apiRequestBytes *prometheus.CounterVec
)

//Config prometheus metrics config
type Config struct {
	// metrics path
	MetricsPath string
	// prometheus namespace
	Namespace string
}

// Metrics return metrics middleware
func Metrics(config Config) baa.HandlerFunc {
	// registry metrics
	if config.Namespace == "" {
		config.Namespace = "default"
	}
	if config.MetricsPath == "" {
		config.MetricsPath = "/metrics"
	}

	initMetrics(config.Namespace)

	return func(c *baa.Context) {
		start := time.Now()

		c.Next()

		// don't keep track of the requests of  promethues collect
		if c.Req.RequestURI == config.MetricsPath {
			return
		}

		collectAPIRequestBytes(c.Req.Method, c.Req.RequestURI, strconv.Itoa(c.Resp.Status()), c.Req.RemoteAddr, float64(c.Req.ContentLength))
		collectAPIRequestTotal(c.Req.Method, c.Req.RequestURI, strconv.Itoa(c.Resp.Status()), c.Req.RemoteAddr)
		collectAPIResponseTime(c.Req.Method, c.Req.RequestURI, strconv.Itoa(c.Resp.Status()), c.Req.RemoteAddr, float64(time.Since(start).Seconds()*1000))
		collectAPIResponseSize(c.Req.Method, c.Req.RequestURI, strconv.Itoa(c.Resp.Status()), c.Req.RemoteAddr, float64(c.Resp.Size()))
	}

}

func initMetrics(namespace string) {
	apiResponseTime = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: namespace,
			Subsystem: "api",
			Name:      "response_time",
			Help:      "Reponse Time of Each Request in Microsecond.",
		},
		[]string{"method", "endpoint", "status", "host"},
	)

	apiResponseSize = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: namespace,
			Subsystem: "api",
			Name:      "response_size",
			Help:      "Reponse Bytes Size of Each Request.",
		},
		[]string{"method", "endpoint", "status", "host"},
	)

	apiRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "api",
			Name:      "request_total",
			Help:      "Total Number of Each Request.",
		},
		[]string{"method", "endpoint", "status", "host"},
	)

	apiRequestBytes = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "api",
			Name:      "reuest_bytes",
			Help:      "Total Data Bytes of Each Request.",
		},
		[]string{"method", "endpoint", "status", "host"},
	)

	prometheus.MustRegister(
		apiResponseTime,
		apiResponseSize,
		apiRequestBytes,
		apiRequestTotal,
	)
}

// collectAPIResponseTime collect api resp time
func collectAPIResponseTime(method, endpoint, status, host string, value float64) {
	apiResponseTime.WithLabelValues(method, endpoint, status, host).Observe(value)
}

// collectAPIResponseSize collect api resp size
func collectAPIResponseSize(method, endpoint, status, host string, value float64) {
	apiResponseSize.WithLabelValues(method, endpoint, status, host).Observe(value)
}

// collectAPIRequestTotal collect api req total
func collectAPIRequestTotal(method, endpoint, status, host string) {
	apiRequestTotal.WithLabelValues(method, endpoint, status, host).Inc()
}

// collectAPIRequestBytes collect api req bytes
func collectAPIRequestBytes(method, endpoint, status, host string, value float64) {
	apiRequestBytes.WithLabelValues(method, endpoint, status, host).Add(value)
}
