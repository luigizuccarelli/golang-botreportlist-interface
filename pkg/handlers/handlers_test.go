package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/trackmate-couchbase-push/pkg/connectors"
	gocb "github.com/couchbase/gocb/v2"
	"github.com/microlib/simple"
)

type FakeCouchbase struct {
}

type Connectors struct {
	Http   *http.Client
	Bucket *FakeCouchbase
}

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewHttpTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func (r *Connectors) Close() error {
	return nil
}

func (r *Connectors) Upsert(collection string, value interface{}, opts *gocb.UpsertOptions) (*gocb.MutationResult, error) {
	if collection == "" {
		return &gocb.MutationResult{}, errors.New("Empty collection value")
	}
	//r.Bucket
	return &gocb.MutationResult{}, nil
}

func NewTestConnectors(file string, code int, logger *simple.Logger) connectors.Clients {

	// we first load the json payload to simulate a call to middleware
	// for now just ignore failures.
	data, err := ioutil.ReadFile(file)
	if err != nil {
		logger.Error(fmt.Sprintf("file data %v\n", err))
		panic(err)
	}
	httpclient := NewHttpTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: code,
			// Send response to be tested

			Body: ioutil.NopCloser(bytes.NewBufferString(string(data))),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	conns := &Connectors{Http: httpclient}
	return conns
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

func TestAllMiddleware(t *testing.T) {

	logger := &simple.Logger{Level: "info"}

	t.Run("IsAlive : should pass", func(t *testing.T) {
		var STATUS int = 200
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v2/sys/info/isalive", nil)
		NewTestConnectors("../../tests/payload.json", STATUS, logger)
		handler := http.HandlerFunc(IsAlive)
		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "IsAlive", rr.Code, STATUS))
		}
	})

	t.Run("AnalyticsHandler : should pass", func(t *testing.T) {
		var STATUS int = 200
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		data, _ := ioutil.ReadFile("../../tests/payload.json")
		req, _ := http.NewRequest("POST", "/api/v1/analytics", bytes.NewBuffer(data))
		conn := NewTestConnectors("../../tests/payload.json", STATUS, logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			AnalyticsHandler(w, r, logger, conn)
		})

		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "IsAlive", rr.Code, STATUS))
		}
	})

	t.Run("AnalyticsHandler : should fail", func(t *testing.T) {
		var STATUS int = 500
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		data := []byte("{ kaka  }")
		req, _ := http.NewRequest("POST", "/api/v1/analytics", bytes.NewBuffer(data))
		conn := NewTestConnectors("../../tests/payload.json", STATUS, logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			AnalyticsHandler(w, r, logger, conn)
		})

		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "IsAlive", rr.Code, STATUS))
		}
	})

}
