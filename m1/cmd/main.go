package main

import (
	"fmt"
	"log"
	"microum/internal/infra"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	config := infra.Load()
	container := infra.NewContainer(config)
	defer container.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/v1/register", container.Handler.Register)
	r.Handle("/metrics", promhttp.Handler()) // Endpoint para o Prometheus coletar dados

	fmt.Printf("🚀 %s running on port %s\n", container.Config.ServerName, container.Config.ServerPort)
	log.Fatal(http.ListenAndServe(":"+container.Config.ServerPort, r))
}
