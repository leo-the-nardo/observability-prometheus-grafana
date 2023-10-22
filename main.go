package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var onlineUsers = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "app_online_users",
	Help: "Number of online users.",
	ConstLabels: prometheus.Labels{
		"website": "e-commerce",
	},
})
var httpRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "app_http_requests_total",
	Help: "Number of HTTP requests.",
}, []string{})
var httpDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name: "app_http_request_duration_seconds",
	Help: "HTTP duration in seconds in all http requests.",
}, []string{"handler"})

func main() {
	registry := prometheus.NewRegistry()
	registry.MustRegister(onlineUsers, httpRequestsTotal, httpDuration)
	go func() {
		for {
			onlineUsers.Set(float64(rand.Intn(2000)))
		}
	}()

	homeHandler := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		res.WriteHeader(http.StatusOK)
		res.Write([]byte("Hello World"))
	})

	durationHandler := promhttp.InstrumentHandlerDuration(
		httpDuration.MustCurryWith(prometheus.Labels{"handler": "home"}),
		promhttp.InstrumentHandlerCounter(httpRequestsTotal, homeHandler))

	http.Handle("/", durationHandler)
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(":8181", nil))
}
