package main

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	activeConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_connections",
			Help: "Number of active connections",
		},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration, activeConnections)
}

func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		activeConnections.Inc()
		defer activeConnections.Dec()

		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()
		httpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, string(rw.statusCode)).Inc()
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.Write([]byte("Hello"))
	})

	route(mux)

	go func() {
		http.ListenAndServe(":8080", metricsMiddleware(mux))
	}()

	go func() {
		cron()
	}()

	fmt.Println("启动成功")

	select {}

}

func route(mux *http.ServeMux) {
	mux.HandleFunc("/test", test)
	mux.HandleFunc("/test2", test2)
}

func test(w http.ResponseWriter, r *http.Request) {
	time.Sleep(100 * time.Millisecond)
	w.Write([]byte("test"))
}

func test2(w http.ResponseWriter, r *http.Request) {
	sumTotal := sumByTotal()
	sumAtomicTotal := sumByAtomic()

	w.Write([]byte(fmt.Sprintf("sumTotal: %d\nsumAtomicTotal: %d\n", sumTotal, sumAtomicTotal)))
}

// 这个函数存在竞态条件，多个 goroutine 同时访问 total 变量，可能导致结果不正确
func sumByTotal() int32 {

	var wg sync.WaitGroup

	total := int32(0)
	wg.Add(100)
	for i := int32(0); i < 100; i++ {
		go func(i int32) {
			total = total + i

			wg.Done()
		}(i)
	}

	wg.Wait()

	return total
}

// 这个函数使用 atomic 包来保证 total 变量的正确性，避免了竞态条件
func sumByAtomic() int32 {
	total := int32(0)
	var wg sync.WaitGroup
	wg.Add(100)
	for i := int32(0); i < 100; i++ {
		go func(i int32) {
			atomic.AddInt32(&total, i)

			wg.Done()
		}(i)
	}

	wg.Wait()

	return total
}

func cron() {
	go cronTest1()
	go cronTest2()
}

func cronTest1() {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			TestPrintln("cronTest1")
		}
	}
}

func cronTest2() {
	ticker := time.NewTicker(time.Second * 11)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			var wg sync.WaitGroup
			wg.Add(100)
			for i := 0; i < 100; i++ {
				go func(i int) {
					TestPrintln(fmt.Sprintf("cronTest2-%2d", i))
					wg.Done()
				}(i)
			}

			wg.Wait()
		}
	}
}

// TestPrintln 输出相关内容
func TestPrintln(tip string) {
	fmt.Println(fmt.Sprintf("[%s], time:%s", tip, time.Now().Format("2006-01-02 15:04:05")))
	time.Sleep(1 * time.Second)
}
