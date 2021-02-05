// +build real

package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"gitea-devops-shared-threefld-cicd.apps.c4.us-east-1.dev.aws.ocp.14west.io/cicd/servisbot-reportlist-interface/pkg/connectors"
	"gitea-devops-shared-threefld-cicd.apps.c4.us-east-1.dev.aws.ocp.14west.io/cicd/servisbot-reportlist-interface/pkg/handlers"
	"gitea-devops-shared-threefld-cicd.apps.c4.us-east-1.dev.aws.ocp.14west.io/cicd/servisbot-reportlist-interface/pkg/validator"
	"github.com/gorilla/mux"
	"github.com/microlib/simple"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	CONTENTTYPE     string = "Content-Type"
	APPLICATIONJSON string = "application/json"
)

var (
	logger       *simple.Logger
	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "s3bucket_manager_http_duration_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"path"})
)

// prometheusMiddleware implements mux.MiddlewareFunc.
func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(CONTENTTYPE, APPLICATIONJSON)
		// use this for cors
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept-Language")
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		next.ServeHTTP(w, r)
		timer.ObserveDuration()
	})
}

// startHttpServer - private function
func startHttpServer(con connectors.Clients) *http.Server {
	srv := &http.Server{Addr: ":" + os.Getenv("SERVER_PORT")}

	r := mux.NewRouter()

	r.Use(prometheusMiddleware)
	r.Path("/api/v2/metrics").Handler(promhttp.Handler())

	r.HandleFunc("/api/v1/list/reports/{offset}/{limit}", func(w http.ResponseWriter, req *http.Request) {
		handlers.ListHandler(w, req, con)
	}).Methods("POST", "OPTIONS")

	r.HandleFunc("/api/v1/reports/count", func(w http.ResponseWriter, req *http.Request) {
		handlers.ReportCountHandler(w, req, con)
	}).Methods("POST", "OPTIONS")

	r.HandleFunc("/api/v1/reports", func(w http.ResponseWriter, req *http.Request) {
		handlers.ReportUpdateHandler(w, req, con)
	}).Methods("POST", "OPTIONS")

	r.HandleFunc("/api/v1/stats", func(w http.ResponseWriter, req *http.Request) {
		handlers.StatsHandler(w, req, con)
	}).Methods("POST", "OPTIONS")

	r.HandleFunc("/api/v1/s3bucket/report", func(w http.ResponseWriter, req *http.Request) {
		handlers.ReportObjectHandler(w, req, con)
	}).Methods("POST", "OPTIONS")

	r.HandleFunc("/api/v2/sys/info/isalive", handlers.IsAlive).Methods("GET")

	sh := http.StripPrefix("/api/v2/api-docs/", http.FileServer(http.Dir("./swaggerui/")))
	r.PathPrefix("/api/v2/api-docs/").Handler(sh)

	http.Handle("/", r)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			con.Error("Httpserver: ListenAndServe() error: " + err.Error())
		}
	}()

	return srv
}

// main - no need for any explanation
func main() {

	if os.Getenv("LOG_LEVEL") == "" {
		logger = &simple.Logger{Level: "info"}
	} else {
		logger = &simple.Logger{Level: os.Getenv("LOG_LEVEL")}
	}

	err := validator.ValidateEnvars(logger)
	if err != nil {
		os.Exit(-1)
	}

	conn := connectors.NewClientConnections(logger)

	srv := startHttpServer(conn)
	logger.Info("Starting server on port " + srv.Addr)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	exit_chan := make(chan int)

	go func() {
		for {
			s := <-c
			switch s {
			case syscall.SIGHUP:
				exit_chan <- 0
			case syscall.SIGINT:
				exit_chan <- 0
			case syscall.SIGTERM:
				exit_chan <- 0
			case syscall.SIGQUIT:
				exit_chan <- 0
			default:
				exit_chan <- 1
			}
		}
	}()

	code := <-exit_chan

	if err := srv.Shutdown(nil); err != nil {
		panic(err)
	}
	logger.Info("Server shutdown successfully")
	os.Exit(code)
}
