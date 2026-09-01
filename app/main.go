package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Response struct {
	Nome    string `json:"nome"`
	Horario string `json:"horario"`
}

var (
	// Métrica para volume de requisições
	requestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total volume of HTTP requests",
		},
		[]string{"path"},
	)
	
	// Métrica para disponibilidade (up)
	serviceUp = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "service_up",
			Help: "Indicates if the service is available (1 = up)",
		},
	)
)

func main() {
	// Indica que o serviço está rodando
	serviceUp.Set(1)

	http.HandleFunc("/projeto-korp", func(w http.ResponseWriter, r *http.Request) {
		requestsTotal.WithLabelValues("/projeto-korp").Inc()
		
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
	})

	// Endpoint para o Prometheus coletar as métricas
	http.Handle("/metrics", promhttp.Handler())

	log.Println("Starting http-server-projeto-korp on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
