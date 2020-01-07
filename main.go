package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/couchbase/gocb.v1"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

var (
	logger       Logger
	bucketClient *gocb.Bucket
)

func startHttpServer(port string) *http.Server {
	srv := &http.Server{Addr: ":" + port}

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/auth", AuthHandler).Methods("GET", "POST")
	r.HandleFunc("/api/v1/analytics", AnalyticsHandler).Methods("POST")
	r.HandleFunc("/api/v2/sys/info/isalive", IsAlive).Methods("GET")
	http.Handle("/", r)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error("Httpserver: ListenAndServe() error: " + err.Error())
		}
	}()

	return srv
}

func main() {

	ValidateEnvars()

	logger.Level = "info"
	if os.Getenv("LOG_LEVEL") != "" {
		logger.Level = os.Getenv("LOG_LEVEL")
	}

	var port string = "9000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	cluster, err := gocb.Connect(os.Getenv("COUCHBASE_HOST"))
	if err != nil {
		logger.Error(fmt.Sprintf("Could not connect to couchbase %v\n", err))
		panic(err)
	}
	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: os.Getenv("COUCHBASE_USER"),
		Password: os.Getenv("COUCHBASE_PASSWORD"),
	})

	// we have a bucket for all analytics
	bucketClient, err = cluster.OpenBucket(os.Getenv("COUCHBASE_BUCKET"), "")
	if err != nil {
		logger.Error(fmt.Sprintf("Could not open bucket %v\n", err))
		panic(err)
	}

	defer bucketClient.Close()
	bucketClient.Manager("", "").CreatePrimaryIndex("", true, false)

	srv := startHttpServer(port)
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

func checkEnvar(item string) {
	name := strings.Split(item, ",")[0]
	required, _ := strconv.ParseBool(strings.Split(item, ",")[1])
	if os.Getenv(name) == "" {
		if required {
			logger.Error(fmt.Sprintf("%s envar is mandatory please set it", name))
			os.Exit(-1)
		} else {
			logger.Error(fmt.Sprintf("%s envar is empty please set it", name))
		}
	}
}

// ValidateEnvars : public call that groups all envar validations
// These envars are set via the openshift template
func ValidateEnvars() {
	items := []string{
		"LOG_LEVEL,false",
		"SERVER_PORT,true",
		"COUCHBASE_HOST,true",
		"COUCHBASE_DATABASE,true",
		"COUCHBASE_USER,true",
		"COUCHBASE_PASSWORD,true",
		"VERSION,true",
		"COUCHBASE_BUCKET,true",
	}
	for x, _ := range items {
		checkEnvar(items[x])
	}
}
