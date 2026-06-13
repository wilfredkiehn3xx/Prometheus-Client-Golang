package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	opsProcessed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
	opsActive = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "myapp_active_ops",
		Help: "The number of active events being processed",
	})
	opsDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "myapp_ops_duration_seconds",
		Help:    "The duration of events being processed",
		Buckets: prometheus.DefBuckets,
	})
	registerMetricsOnce sync.Once
)

func registerMetrics() {
	registerMetricsOnce.Do(func() {
		prometheus.MustRegister(opsProcessed, opsActive, opsDuration)
	})
}

type Server struct {
	Addr       string
	httpServer *http.Server
}

func NewServer() *Server {
	registerMetrics()

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		opsProcessed.Inc()
		fmt.Fprintln(w, "Hello, World!")
	})

	return &Server{
		Addr: ":8080",
		httpServer: &http.Server{
			Handler: mux,
		},
	}
}

func (s *Server) Start() error {
	s.httpServer.Addr = s.Addr
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}

func main() {
	srv := NewServer()
	fmt.Println("Starting server on :8080...")
	if err := srv.Start(); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
