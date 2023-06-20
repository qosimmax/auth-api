package server

import (
	"fmt"
	"net/http"

	"gitlab.com/route-kz/auth-api/server/internal/handler"
	"gitlab.com/route-kz/auth-api/server/internal/middleware"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const v1API string = "/api/v1"

// setupRoutes - the root route function.
func (s *Server) setupRoutes() {
	s.Router.Handle("/metrics", promhttp.Handler()).Name("Metrics")
	s.Router.HandleFunc("/_healthz", handler.Healthz).Methods(http.MethodGet).Name("Health")

	api := s.Router.PathPrefix(v1API).Subrouter()
	api.HandleFunc("/tokens", handler.CreateToken(s.DB)).Methods(http.MethodPost).Name(fmt.Sprintf("CreateToken"))
	api.HandleFunc("/refresh-tokens", handler.RefreshToken(s.DB)).Methods(http.MethodPost).Name(fmt.Sprintf("RefreshToken"))
	api.HandleFunc("/tokens", handler.Identity(s.DB)).Methods(http.MethodGet).Name("Identity")
	api.HandleFunc("/personal-data", handler.PersonalData(s.DB)).Methods(http.MethodGet).Name("PersonalData")

	addTracingAndMetrics(api)
}

// addTracingAndMetrics - Adds tracing and metrics to a router.
func addTracingAndMetrics(r *mux.Router) {
	tm := middleware.TraceMetrics{}
	r.Use(tm.TraceMiddleware)
	r.Use(tm.MetricsMiddleware)
}
