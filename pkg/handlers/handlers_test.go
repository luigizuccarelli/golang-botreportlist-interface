package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/servisbot-reportlist-interface/pkg/connectors"
	"github.com/gorilla/mux"
	"github.com/microlib/simple"
)

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("Inject (force) readAll test error")
}

// TestAllHandlers - main test entry point
func TestAllHandlers(t *testing.T) {

	logger := &simple.Logger{Level: "trace"}

	t.Run("IsAlive : should pass", func(t *testing.T) {
		var STATUS int = 200
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v2/sys/info/isalive", nil)
		connectors.NewTestConnectors(STATUS, logger)
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

	t.Run("ListHandler : should pass", func(t *testing.T) {
		var STATUS int = 200
		os.Setenv("TOKEN", "1212121")
		os.Setenv("JWT_SECRETKEY", "Thr33f0ldSystems?CSsD!@%2^")

		requestPayload := `{  "email": "cduffy@tfd.ie", "jwttoken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1OTA3NTY4MjAsInN5c3RlbSI6ImNvbnRhY3QtZm9ybSIsImN1c3RvbWVyTnVtYmVyIjoiMDAwMTE5OTQ0MTYwIiwidXNlciI6ImNkdWZmeUB0ZmQuaWUifQ.fisOWBMqnbzzcNQpqO6Cmu6DEMjroaZYgTsAeEmR36A" }`
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/list/reports/0/10", bytes.NewBuffer([]byte(requestPayload)))
		conn := connectors.NewTestConnectors(STATUS, logger)
		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"lastobject": "test",
		}
		req = mux.SetURLVars(req, vars)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ListHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "ListHandler", rr.Code, STATUS))
		}
	})

	t.Run("ListHandler : should fail (force read body error)", func(t *testing.T) {
		var STATUS int = 500
		os.Setenv("TOKEN", "1212121")
		os.Setenv("JWT_SECRETKEY", "Thr33f0ldSystems?CSsD!@%2^")
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/list/reports/0/10", errReader(0))
		conn := connectors.NewTestConnectors(STATUS, logger)
		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"lastobject": "test",
		}
		req = mux.SetURLVars(req, vars)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ListHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "ListHandler", rr.Code, STATUS))
		}
	})

	t.Run("ListHandler : should fail (bad request json)", func(t *testing.T) {
		var STATUS int = 500
		os.Setenv("TOKEN", "1212121")
		os.Setenv("JWT_SECRETKEY", "Thr33f0ldSystems?CSsD!@%2^")

		requestPayload := `{  email": "cduffy@tfd.ie", "jwttoken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1OTA3NTY4MjAsInN5c3RlbSI6ImNvbnRhY3QtZm9ybSIsImN1c3RvbWVyTnVtYmVyIjoiMDAwMTE5OTQ0MTYwIiwidXNlciI6ImNkdWZmeUB0ZmQuaWUifQ.fisOWBMqnbzzcNQpqO6Cmu6DEMjroaZYgTsAeEmR36A" }`
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/list/reports/0/10", bytes.NewBuffer([]byte(requestPayload)))
		conn := connectors.NewTestConnectors(STATUS, logger)
		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"lastobject": "test",
		}
		req = mux.SetURLVars(req, vars)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ListHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "ListHandler", rr.Code, STATUS))
		}
	})

	t.Run("ListHandler : should fail (bad jwt token)", func(t *testing.T) {
		var STATUS int = 403
		os.Setenv("TOKEN", "1212121")
		os.Setenv("JWT_SECRETKEY", "Thr33f0ldSystems?CSsD!@%2^")

		requestPayload := `{  "email": "cduffy@tfd.ie", "jwttoke": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1OTA3NTY4MjAsInN5c3RlbSI6ImNvbnRhY3QtZm9ybSIsImN1c3RvbWVyTnVtYmVyIjoiMDAwMTE5OTQ0MTYwIiwidXNlciI6ImNkdWZmeUB0ZmQuaWUifQ.fisOWBMqnbzzcNQpqO6Cmu6DEMjroaZYgTsAeEmR36A" }`
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/list/reports/0/10", bytes.NewBuffer([]byte(requestPayload)))
		conn := connectors.NewTestConnectors(STATUS, logger)
		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"lastobject": "test",
		}
		req = mux.SetURLVars(req, vars)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ListHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "ListHandler", rr.Code, STATUS))
		}
	})

	t.Run("ListHandler : should fail (force db error)", func(t *testing.T) {
		var STATUS int = 500
		os.Setenv("TOKEN", "1212121")
		os.Setenv("JWT_SECRETKEY", "Thr33f0ldSystems?CSsD!@%2^")

		requestPayload := `{  "email": "cduffy@tfd.ie", "jwttoken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1OTA3NTY4MjAsInN5c3RlbSI6ImNvbnRhY3QtZm9ybSIsImN1c3RvbWVyTnVtYmVyIjoiMDAwMTE5OTQ0MTYwIiwidXNlciI6ImNkdWZmeUB0ZmQuaWUifQ.fisOWBMqnbzzcNQpqO6Cmu6DEMjroaZYgTsAeEmR36A" }`
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/list/reports/0/10", bytes.NewBuffer([]byte(requestPayload)))
		conn := connectors.NewTestConnectors(STATUS, logger)
		conn.Meta("true")
		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"lastobject": "test",
		}
		req = mux.SetURLVars(req, vars)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ListHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "ListHandler", rr.Code, STATUS))
		}
	})

	t.Run("ReportUpdateHandler : should pass", func(t *testing.T) {
		var STATUS int = 200
		os.Setenv("TOKEN", "1212121")
		os.Setenv("JWT_SECRETKEY", "Thr33f0ldSystems?CSsD!@%2^")

		requestPayload := `{ "data": {"id":"test","servisbotstats":{"emailclassification":"test","processoutcome":"test","userclassification":"test","success": false}}, "jwttoken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1OTA3NTY4MjAsInN5c3RlbSI6ImNvbnRhY3QtZm9ybSIsImN1c3RvbWVyTnVtYmVyIjoiMDAwMTE5OTQ0MTYwIiwidXNlciI6ImNkdWZmeUB0ZmQuaWUifQ.fisOWBMqnbzzcNQpqO6Cmu6DEMjroaZYgTsAeEmR36A" }`
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/reports/", bytes.NewBuffer([]byte(requestPayload)))
		conn := connectors.NewTestConnectors(STATUS, logger)
		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"lastobject": "test",
		}
		req = mux.SetURLVars(req, vars)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ReportUpdateHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "ReportUpdateHandler", rr.Code, STATUS))
		}
	})

	t.Run("ReportUpdateHandler : should fail (force read body error)", func(t *testing.T) {
		var STATUS int = 500
		os.Setenv("TOKEN", "1212121")
		os.Setenv("JWT_SECRETKEY", "Thr33f0ldSystems?CSsD!@%2^")

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/reports/", errReader(0))
		conn := connectors.NewTestConnectors(STATUS, logger)
		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"lastobject": "test",
		}
		req = mux.SetURLVars(req, vars)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ReportUpdateHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "ReportUpdateHandler", rr.Code, STATUS))
		}
	})

	t.Run("ReportUpdateHandler : should fail (bad request json)", func(t *testing.T) {
		var STATUS int = 500
		os.Setenv("TOKEN", "1212121")
		os.Setenv("JWT_SECRETKEY", "Thr33f0ldSystems?CSsD!@%2^")

		requestPayload := `{ data": {"id":"test","servisbotstats":{"emailclassification":"test","processoutcome":"test","userclassification":"test","success": false}}, "jwttoken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1OTA3NTY4MjAsInN5c3RlbSI6ImNvbnRhY3QtZm9ybSIsImN1c3RvbWVyTnVtYmVyIjoiMDAwMTE5OTQ0MTYwIiwidXNlciI6ImNkdWZmeUB0ZmQuaWUifQ.fisOWBMqnbzzcNQpqO6Cmu6DEMjroaZYgTsAeEmR36A" }`
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/reports/", bytes.NewBuffer([]byte(requestPayload)))
		conn := connectors.NewTestConnectors(STATUS, logger)
		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"lastobject": "test",
		}
		req = mux.SetURLVars(req, vars)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ReportUpdateHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "ReportUpdateHandler", rr.Code, STATUS))
		}
	})

	t.Run("ReportUpdateHandler : should fail (bad jwt token)", func(t *testing.T) {
		var STATUS int = 403
		os.Setenv("TOKEN", "1212121")
		os.Setenv("JWT_SECRETKEY", "Thr33f0ldSystems?CSsD!@%2^")

		requestPayload := `{ "data": {"id":"test","servisbotstats":{"emailclassification":"test","processoutcome":"test","userclassification":"test","success": false}}, "jwttoke": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1OTA3NTY4MjAsInN5c3RlbSI6ImNvbnRhY3QtZm9ybSIsImN1c3RvbWVyTnVtYmVyIjoiMDAwMTE5OTQ0MTYwIiwidXNlciI6ImNkdWZmeUB0ZmQuaWUifQ.fisOWBMqnbzzcNQpqO6Cmu6DEMjroaZYgTsAeEmR36A" }`
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/reports/", bytes.NewBuffer([]byte(requestPayload)))
		conn := connectors.NewTestConnectors(STATUS, logger)
		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"lastobject": "test",
		}
		req = mux.SetURLVars(req, vars)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ReportUpdateHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "ReportUpdateHandler", rr.Code, STATUS))
		}
	})

	t.Run("ReportUpdateHandler : should fail (force db error)", func(t *testing.T) {
		var STATUS int = 500
		os.Setenv("TOKEN", "1212121")
		os.Setenv("JWT_SECRETKEY", "Thr33f0ldSystems?CSsD!@%2^")

		requestPayload := `{ "data": {"id":"test","servisbotstats":{"emailclassification":"test","processoutcome":"test","userclassification":"test","success": false}}, "jwttoken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1OTA3NTY4MjAsInN5c3RlbSI6ImNvbnRhY3QtZm9ybSIsImN1c3RvbWVyTnVtYmVyIjoiMDAwMTE5OTQ0MTYwIiwidXNlciI6ImNkdWZmeUB0ZmQuaWUifQ.fisOWBMqnbzzcNQpqO6Cmu6DEMjroaZYgTsAeEmR36A" }`
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/reports/", bytes.NewBuffer([]byte(requestPayload)))
		conn := connectors.NewTestConnectors(STATUS, logger)
		conn.Meta("true")
		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"lastobject": "test",
		}
		req = mux.SetURLVars(req, vars)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ReportUpdateHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "ReportUpdateHandler", rr.Code, STATUS))
		}
	})

	t.Run("StatsHandler : should pass", func(t *testing.T) {
		var STATUS int = 200
		os.Setenv("TOKEN", "1212121")
		os.Setenv("JWT_SECRETKEY", "Thr33f0ldSystems?CSsD!@%2^")

		requestPayload := `{ "jwttoken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1OTA3NTY4MjAsInN5c3RlbSI6ImNvbnRhY3QtZm9ybSIsImN1c3RvbWVyTnVtYmVyIjoiMDAwMTE5OTQ0MTYwIiwidXNlciI6ImNkdWZmeUB0ZmQuaWUifQ.fisOWBMqnbzzcNQpqO6Cmu6DEMjroaZYgTsAeEmR36A" }`
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/stats", bytes.NewBuffer([]byte(requestPayload)))
		conn := connectors.NewTestConnectors(STATUS, logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			StatsHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "StatsHandler", rr.Code, STATUS))
		}
	})

	t.Run("StatsHandler : should fail (force read body errro)", func(t *testing.T) {
		var STATUS int = 500
		os.Setenv("TOKEN", "1212121")
		os.Setenv("JWT_SECRETKEY", "Thr33f0ldSystems?CSsD!@%2^")

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/stats", errReader(0))
		conn := connectors.NewTestConnectors(STATUS, logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			StatsHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "StatsHandler", rr.Code, STATUS))
		}
	})

	t.Run("StatsHandler : should fail (bad json request)", func(t *testing.T) {
		var STATUS int = 500
		os.Setenv("TOKEN", "1212121")
		os.Setenv("JWT_SECRETKEY", "Thr33f0ldSystems?CSsD!@%2^")

		requestPayload := `{ jwttoken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1OTA3NTY4MjAsInN5c3RlbSI6ImNvbnRhY3QtZm9ybSIsImN1c3RvbWVyTnVtYmVyIjoiMDAwMTE5OTQ0MTYwIiwidXNlciI6ImNkdWZmeUB0ZmQuaWUifQ.fisOWBMqnbzzcNQpqO6Cmu6DEMjroaZYgTsAeEmR36A" }`
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/stats", bytes.NewBuffer([]byte(requestPayload)))
		conn := connectors.NewTestConnectors(STATUS, logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			StatsHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "StatsHandler", rr.Code, STATUS))
		}
	})

	t.Run("StatsHandler : should fail (bad jwt token)", func(t *testing.T) {
		var STATUS int = 403
		os.Setenv("TOKEN", "1212121")
		os.Setenv("JWT_SECRETKEY", "Thr33f0ldSystems?CSsD!@%2^")

		requestPayload := `{ "jwttoke": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1OTA3NTY4MjAsInN5c3RlbSI6ImNvbnRhY3QtZm9ybSIsImN1c3RvbWVyTnVtYmVyIjoiMDAwMTE5OTQ0MTYwIiwidXNlciI6ImNkdWZmeUB0ZmQuaWUifQ.fisOWBMqnbzzcNQpqO6Cmu6DEMjroaZYgTsAeEmR36A" }`
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/stats", bytes.NewBuffer([]byte(requestPayload)))
		conn := connectors.NewTestConnectors(STATUS, logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			StatsHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "StatsHandler", rr.Code, STATUS))
		}
	})

	t.Run("StatsHandler : should fail (force db error)", func(t *testing.T) {
		var STATUS int = 500
		os.Setenv("TOKEN", "1212121")
		os.Setenv("JWT_SECRETKEY", "Thr33f0ldSystems?CSsD!@%2^")

		requestPayload := `{ "jwttoken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1OTA3NTY4MjAsInN5c3RlbSI6ImNvbnRhY3QtZm9ybSIsImN1c3RvbWVyTnVtYmVyIjoiMDAwMTE5OTQ0MTYwIiwidXNlciI6ImNkdWZmeUB0ZmQuaWUifQ.fisOWBMqnbzzcNQpqO6Cmu6DEMjroaZYgTsAeEmR36A" }`
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/stats", bytes.NewBuffer([]byte(requestPayload)))
		conn := connectors.NewTestConnectors(STATUS, logger)
		conn.Meta("true")
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			StatsHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "StatsHandler", rr.Code, STATUS))
		}
	})

	t.Run("ReportObjectHandler : should pass", func(t *testing.T) {
		var STATUS int = 200
		os.Setenv("TOKEN", "1212121")
		os.Setenv("JWT_SECRETKEY", "Thr33f0ldSystems?CSsD!@%2^")

		requestPayload := `{ "id":"13124","jwttoken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1OTA3NTY4MjAsInN5c3RlbSI6ImNvbnRhY3QtZm9ybSIsImN1c3RvbWVyTnVtYmVyIjoiMDAwMTE5OTQ0MTYwIiwidXNlciI6ImNkdWZmeUB0ZmQuaWUifQ.fisOWBMqnbzzcNQpqO6Cmu6DEMjroaZYgTsAeEmR36A" }`
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/s3bucket/report", bytes.NewBuffer([]byte(requestPayload)))
		conn := connectors.NewTestConnectors(STATUS, logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ReportObjectHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "StatsHandler", rr.Code, STATUS))
		}
	})

	t.Run("ReportObjectHandler : should fail (force error)", func(t *testing.T) {
		var STATUS int = 500
		os.Setenv("TOKEN", "1212121")
		os.Setenv("JWT_SECRETKEY", "Thr33f0ldSystems?CSsD!@%2^")

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/s3bucket/report", errReader(0))
		conn := connectors.NewTestConnectors(STATUS, logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ReportObjectHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "ReportObjectHandler", rr.Code, STATUS))
		}
	})

	t.Run("ReportObjectHandler : should fail (json)", func(t *testing.T) {
		var STATUS int = 500
		os.Setenv("TOKEN", "1212121")
		os.Setenv("JWT_SECRETKEY", "Thr33f0ldSystems?CSsD!@%2^")

		requestPayload := `{  email": "cduffy@tfd.ie", "jwttoken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1OTA3NTY4MjAsInN5c3RlbSI6ImNvbnRhY3QtZm9ybSIsImN1c3RvbWVyTnVtYmVyIjoiMDAwMTE5OTQ0MTYwIiwidXNlciI6ImNkdWZmeUB0ZmQuaWUifQ.fisOWBMqnbzzcNQpqO6Cmu6DEMjroaZYgTsAeEmR36A" }`
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/is3bucket/report", bytes.NewBuffer([]byte(requestPayload)))
		conn := connectors.NewTestConnectors(STATUS, logger)
		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"key": "test",
		}
		req = mux.SetURLVars(req, vars)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ReportObjectHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "ReportObjectHandler", rr.Code, STATUS))
		}
	})

	t.Run("ReportObjectHandler : should fail (jwt token)", func(t *testing.T) {
		var STATUS int = 403
		os.Setenv("TOKEN", "1212121")
		os.Setenv("JWT_SECRETKEY", "Thr")

		requestPayload := `{"id":"12345", "jwttoken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1OTA3NTY4MjAsInN5c3RlbSI6ImNvbnRhY3QtZm9ybSIsImN1c3RvbWVyTnVtYmVyIjoiMDAwMTE5OTQ0MTYwIiwidXNlciI6ImNkdWZmeUB0ZmQuaWUifQ.fisOWBMqnbzzcNQpqO6Cmu6DEMjroaZYgTsAeEmR36A" }`
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/s3bucket/report", bytes.NewBuffer([]byte(requestPayload)))
		conn := connectors.NewTestConnectors(STATUS, logger)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ReportObjectHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "ReportObjectHandler", rr.Code, STATUS))
		}
	})

	t.Run("ReportObjectHandler : should fail (force GetObject error)", func(t *testing.T) {
		var STATUS int = 500
		os.Setenv("TOKEN", "1212121")
		os.Setenv("JWT_SECRETKEY", "Thr33f0ldSystems?CSsD!@%2^")

		requestPayload := `{  "id": "123456789", "jwttoken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1OTA3NTY4MjAsInN5c3RlbSI6ImNvbnRhY3QtZm9ybSIsImN1c3RvbWVyTnVtYmVyIjoiMDAwMTE5OTQ0MTYwIiwidXNlciI6ImNkdWZmeUB0ZmQuaWUifQ.fisOWBMqnbzzcNQpqO6Cmu6DEMjroaZYgTsAeEmR36A" }`
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/s3bucket/report", bytes.NewBuffer([]byte(requestPayload)))
		conn := connectors.NewTestConnectors(STATUS, logger)
		conn.Meta("true")
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ReportObjectHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "ReportObjectHandler", rr.Code, STATUS))
		}
	})
}
