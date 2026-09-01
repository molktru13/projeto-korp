package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Response é a estrutura de resposta JSON do endpoint /projeto-korp
type Response struct {
	Nome    string `json:"nome"`
	Horario string `json:"horario"`
}

// responseWriter é um wrapper para capturar status code e tamanho da resposta
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	bytes      int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.bytes += n
	return n, err
}

var (
	// Disponibilidade do serviço
	serviceUp = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "service_up",
			Help: "Indica se o serviço está disponível (1 = UP, 0 = DOWN)",
		},
	)

	// Contador total de requisições HTTP com labels
	requestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Número total de requisições HTTP processadas",
		},
		[]string{"path", "method", "status_code"},
	)

	// Histograma de duração das requisições (latência)
	requestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duração das requisições HTTP em segundos",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"path", "method", "status_code"},
	)

	// Requisições em andamento
	requestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Número de requisições HTTP sendo processadas no momento",
		},
	)

	// Histograma do tamanho das respostas
	responseSizeBytes = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "Tamanho das respostas HTTP em bytes",
			Buckets: []float64{100, 500, 1000, 5000, 10000, 50000, 100000},
		},
		[]string{"path", "method"},
	)
)

// instrumentHandler envolve um handler com métricas de observabilidade
func instrumentHandler(path string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestsInFlight.Inc()
		defer requestsInFlight.Dec()

		start := time.Now()
		rw := newResponseWriter(w)

		next(rw, r)

		duration := time.Since(start).Seconds()
		statusCode := strconv.Itoa(rw.statusCode)

		requestsTotal.WithLabelValues(path, r.Method, statusCode).Inc()
		requestDuration.WithLabelValues(path, r.Method, statusCode).Observe(duration)
		responseSizeBytes.WithLabelValues(path, r.Method).Observe(float64(rw.bytes))
	}
}

func projetoKorpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := Response{
		Nome:    "Projeto Korp",
		Horario: time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Indica que o serviço está rodando
	serviceUp.Set(1)

	// Endpoint principal com instrumentação
	http.HandleFunc("/projeto-korp", instrumentHandler("/projeto-korp", projetoKorpHandler))

	// Endpoint para o Prometheus coletar as métricas
	http.Handle("/metrics", promhttp.Handler())

	log.Println("Starting http-server-projeto-korp on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
