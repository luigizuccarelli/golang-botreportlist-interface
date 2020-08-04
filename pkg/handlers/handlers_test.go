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

	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/servisbot-s3bucket-manager/pkg/connectors"
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
		connectors.NewTestConnectors("../../tests/payload.json", STATUS, logger)
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

	t.Run("ListBucketHandler : should pass", func(t *testing.T) {
		var STATUS int = 200
		os.Setenv("TOKEN", "1212121")
		os.Setenv("JWT_SECRETKEY", "Thr33f0ldSystems?CSsD!@%2^")

		requestPayload := `{  "email": "cduffy@tfd.ie", "jwttoken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1OTA3NTY4MjAsInN5c3RlbSI6ImNvbnRhY3QtZm9ybSIsImN1c3RvbWVyTnVtYmVyIjoiMDAwMTE5OTQ0MTYwIiwidXNlciI6ImNkdWZmeUB0ZmQuaWUifQ.fisOWBMqnbzzcNQpqO6Cmu6DEMjroaZYgTsAeEmR36A" }`
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/list/reports/234324", bytes.NewBuffer([]byte(requestPayload)))
		conn := connectors.NewTestConnectors("../../tests/payload.json", STATUS, logger)
		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"lastobject": "test",
		}
		req = mux.SetURLVars(req, vars)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ListBucketHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "ListBucketHandler", rr.Code, STATUS))
		}
	})

	t.Run("ListBucketHandler : should fail (force error)", func(t *testing.T) {
		var STATUS int = 500
		os.Setenv("TOKEN", "1212121")
		os.Setenv("JWT_SECRETKEY", "Thr33f0ldSystems?CSsD!@%2^")

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/list/reports/234324", errReader(0))
		conn := connectors.NewTestConnectors("../../tests/payload.json", STATUS, logger)
		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"lastobject": "test",
		}
		req = mux.SetURLVars(req, vars)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ListBucketHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "ListBucketHandler", rr.Code, STATUS))
		}
	})

	t.Run("ListBucketHandler : should fail (json data)", func(t *testing.T) {
		var STATUS int = 500
		os.Setenv("TOKEN", "1212121")
		os.Setenv("JWT_SECRETKEY", "Thr33f0ldSystems?CSsD!@%2^")

		requestPayload := `{  jwttoken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1OTA3NTY4MjAsInN5c3RlbSI6ImNvbnRhY3QtZm9ybSIsImN1c3RvbWVyTnVtYmVyIjoiMDAwMTE5OTQ0MTYwIiwidXNlciI6ImNkdWZmeUB0ZmQuaWUifQ.fisOWBMqnbzzcNQpqO6Cmu6DEMjroaZYgTsAeEmR36A" }`
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/list/reports/234324", bytes.NewBuffer([]byte(requestPayload)))
		conn := connectors.NewTestConnectors("../../tests/payload.json", STATUS, logger)
		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"lastobject": "test",
		}
		req = mux.SetURLVars(req, vars)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ListBucketHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "ListBucketHandler", rr.Code, STATUS))
		}
	})

	t.Run("ListBucketHandler : should fail (JWT)", func(t *testing.T) {
		var STATUS int = 403
		os.Setenv("TOKEN", "1212121")
		os.Setenv("JWT_SECRETKEY", "Thr")

		requestPayload := `{ "data":"test", "jwttoken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1OTA3NTY4MjAsInN5c3RlbSI6ImNvbnRhY3QtZm9ybSIsImN1c3RvbWVyTnVtYmVyIjoiMDAwMTE5OTQ0MTYwIiwidXNlciI6ImNkdWZmeUB0ZmQuaWUifQ.fisOWBMqnbzzcNQpqO6Cmu6DEMjroaZYgTsAeEmR36A" }`
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/list/reports/234324", bytes.NewBuffer([]byte(requestPayload)))
		conn := connectors.NewTestConnectors("../../tests/payload.json", STATUS, logger)
		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"lastobject": "test",
		}
		req = mux.SetURLVars(req, vars)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ListBucketHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "ListBucketHandler", rr.Code, STATUS))
		}
	})

	t.Run("EmailObjectHandler : should pass", func(t *testing.T) {
		var STATUS int = 200
		os.Setenv("TOKEN", "1212121")
		os.Setenv("JWT_SECRETKEY", "Thr33f0ldSystems?CSsD!@%2^")

		requestPayload := `{  "email": "cduffy@tfd.ie", "jwttoken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1OTA3NTY4MjAsInN5c3RlbSI6ImNvbnRhY3QtZm9ybSIsImN1c3RvbWVyTnVtYmVyIjoiMDAwMTE5OTQ0MTYwIiwidXNlciI6ImNkdWZmeUB0ZmQuaWUifQ.fisOWBMqnbzzcNQpqO6Cmu6DEMjroaZYgTsAeEmR36A" }`
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/emails/234324", bytes.NewBuffer([]byte(requestPayload)))
		conn := connectors.NewTestConnectors("../../tests/payload.json", STATUS, logger)
		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"key": "test",
		}
		req = mux.SetURLVars(req, vars)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			EmailObjectHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "EmailObjectHandler", rr.Code, STATUS))
		}
	})

	t.Run("EmailObjectHandler : should fail (force error)", func(t *testing.T) {
		var STATUS int = 500
		os.Setenv("TOKEN", "1212121")
		os.Setenv("JWT_SECRETKEY", "Thr33f0ldSystems?CSsD!@%2^")

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/emails/4324324324", errReader(0))
		conn := connectors.NewTestConnectors("../../tests/payload.json", STATUS, logger)
		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"key": "test",
		}
		req = mux.SetURLVars(req, vars)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			EmailObjectHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "EmailObjectHandler", rr.Code, STATUS))
		}
	})

	t.Run("EmailObjectHandler : should fail (json)", func(t *testing.T) {
		var STATUS int = 500
		os.Setenv("TOKEN", "1212121")
		os.Setenv("JWT_SECRETKEY", "Thr33f0ldSystems?CSsD!@%2^")

		requestPayload := `{  email": "cduffy@tfd.ie", "jwttoken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1OTA3NTY4MjAsInN5c3RlbSI6ImNvbnRhY3QtZm9ybSIsImN1c3RvbWVyTnVtYmVyIjoiMDAwMTE5OTQ0MTYwIiwidXNlciI6ImNkdWZmeUB0ZmQuaWUifQ.fisOWBMqnbzzcNQpqO6Cmu6DEMjroaZYgTsAeEmR36A" }`
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/emails/234324", bytes.NewBuffer([]byte(requestPayload)))
		conn := connectors.NewTestConnectors("../../tests/payload.json", STATUS, logger)
		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"key": "test",
		}
		req = mux.SetURLVars(req, vars)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			EmailObjectHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "EmailObjectHandler", rr.Code, STATUS))
		}
	})

	t.Run("EmailObjectHandler : should fail (JWT)", func(t *testing.T) {
		var STATUS int = 403
		os.Setenv("TOKEN", "1212121")
		os.Setenv("JWT_SECRETKEY", "Thr")

		requestPayload := `{  "email": "cduffy@tfd.ie", "jwttoken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1OTA3NTY4MjAsInN5c3RlbSI6ImNvbnRhY3QtZm9ybSIsImN1c3RvbWVyTnVtYmVyIjoiMDAwMTE5OTQ0MTYwIiwidXNlciI6ImNkdWZmeUB0ZmQuaWUifQ.fisOWBMqnbzzcNQpqO6Cmu6DEMjroaZYgTsAeEmR36A" }`
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/emails/234324", bytes.NewBuffer([]byte(requestPayload)))
		conn := connectors.NewTestConnectors("../../tests/payload.json", STATUS, logger)
		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"key": "test",
		}
		req = mux.SetURLVars(req, vars)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			EmailObjectHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "EmailObjectHandler", rr.Code, STATUS))
		}
	})

	t.Run("EmailObjectHandler : should fail (force GetObject error)", func(t *testing.T) {
		var STATUS int = 500
		os.Setenv("TOKEN", "1212121")
		os.Setenv("JWT_SECRETKEY", "Thr33f0ldSystems?CSsD!@%2^")

		requestPayload := `{  "email": "cduffy@tfd.ie", "jwttoken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1OTA3NTY4MjAsInN5c3RlbSI6ImNvbnRhY3QtZm9ybSIsImN1c3RvbWVyTnVtYmVyIjoiMDAwMTE5OTQ0MTYwIiwidXNlciI6ImNkdWZmeUB0ZmQuaWUifQ.fisOWBMqnbzzcNQpqO6Cmu6DEMjroaZYgTsAeEmR36A" }`
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/emails/234324", bytes.NewBuffer([]byte(requestPayload)))
		conn := connectors.NewTestConnectors("../../tests/payload.json", STATUS, logger)
		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"key": "test",
		}
		conn.Meta("true")
		req = mux.SetURLVars(req, vars)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			EmailObjectHandler(w, r, conn)
		})
		handler.ServeHTTP(rr, req)
		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "EmailObjectHandler", rr.Code, STATUS))
		}
	})

	t.Run("ReportObjectHandler : should pass (push)", func(t *testing.T) {
		var STATUS int = 200
		os.Setenv("TOKEN", "1212121")
		os.Setenv("JWT_SECRETKEY", "Thr33f0ldSystems?CSsD!@%2^")

		requestPayload := `{  "data": "{\"field\":\"value-cduffy@tfd.ie\"}", "jwttoken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1OTA3NTY4MjAsInN5c3RlbSI6ImNvbnRhY3QtZm9ybSIsImN1c3RvbWVyTnVtYmVyIjoiMDAwMTE5OTQ0MTYwIiwidXNlciI6ImNkdWZmeUB0ZmQuaWUifQ.fisOWBMqnbzzcNQpqO6Cmu6DEMjroaZYgTsAeEmR36A" }`
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/reports/4324324324", bytes.NewBuffer([]byte(requestPayload)))
		conn := connectors.NewTestConnectors("../../tests/payload.json", STATUS, logger)
		conn.SetMode("push")
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

	t.Run("ReportObjectHandler : should pass (pull)", func(t *testing.T) {
		var STATUS int = 200
		os.Setenv("TOKEN", "1212121")
		os.Setenv("JWT_SECRETKEY", "Thr33f0ldSystems?CSsD!@%2^")

		requestPayload := `{  "data": "{\"field\":\"value-cduffy@tfd.ie\"}", "jwttoken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1OTA3NTY4MjAsInN5c3RlbSI6ImNvbnRhY3QtZm9ybSIsImN1c3RvbWVyTnVtYmVyIjoiMDAwMTE5OTQ0MTYwIiwidXNlciI6ImNkdWZmeUB0ZmQuaWUifQ.fisOWBMqnbzzcNQpqO6Cmu6DEMjroaZYgTsAeEmR36A" }`
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/reports/4324324324", bytes.NewBuffer([]byte(requestPayload)))
		conn := connectors.NewTestConnectors("../../tests/payload.json", STATUS, logger)
		conn.SetMode("pull")
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

	t.Run("ReportObjectHandler : should fail (force error readall)", func(t *testing.T) {
		var STATUS int = 500
		os.Setenv("TOKEN", "1212121")
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/reports/4324324324", errReader(0))
		conn := connectors.NewTestConnectors("../../tests/payload.json", STATUS, logger)
		conn.SetMode("push")
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

	t.Run("ReportObjectHandler : should fail (json data)", func(t *testing.T) {
		var STATUS int = 500
		os.Setenv("TOKEN", "1212121")
		os.Setenv("JWT_SECRETKEY", "Thr33f0ldSystems?CSsD!@%2^")

		requestPayload := `{  data": "{\"field\":\"value-cduffy@tfd.ie\"}", "jwttoken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1OTA3NTY4MjAsInN5c3RlbSI6ImNvbnRhY3QtZm9ybSIsImN1c3RvbWVyTnVtYmVyIjoiMDAwMTE5OTQ0MTYwIiwidXNlciI6ImNkdWZmeUB0ZmQuaWUifQ.fisOWBMqnbzzcNQpqO6Cmu6DEMjroaZYgTsAeEmR36A" }`
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/reports/4324324324", bytes.NewBuffer([]byte(requestPayload)))
		conn := connectors.NewTestConnectors("../../tests/payload.json", STATUS, logger)
		conn.SetMode("push")
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

	t.Run("ReportObjectHandler : should fail (JWT)", func(t *testing.T) {
		var STATUS int = 403
		os.Setenv("TOKEN", "1212121")
		os.Setenv("JWT_SECRETKEY", "Thr")

		requestPayload := `{  "data": "{\"field\":\"value-cduffy@tfd.ie\"}", "jwttoken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1OTA3NTY4MjAsInN5c3RlbSI6ImNvbnRhY3QtZm9ybSIsImN1c3RvbWVyTnVtYmVyIjoiMDAwMTE5OTQ0MTYwIiwidXNlciI6ImNkdWZmeUB0ZmQuaWUifQ.fisOWBMqnbzzcNQpqO6Cmu6DEMjroaZYgTsAeEmR36A" }`
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/reports/4324324324", bytes.NewBuffer([]byte(requestPayload)))
		conn := connectors.NewTestConnectors("../../tests/payload.json", STATUS, logger)
		conn.SetMode("push")
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

	t.Run("ReportObjectHandler : should fail (force ReportObject error)", func(t *testing.T) {
		var STATUS int = 500
		os.Setenv("TOKEN", "1212121")
		os.Setenv("JWT_SECRETKEY", "Thr33f0ldSystems?CSsD!@%2^")

		requestPayload := `{  "data": "{\"field\":\"value-cduffy@tfd.ie\"}", "jwttoken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1OTA3NTY4MjAsInN5c3RlbSI6ImNvbnRhY3QtZm9ybSIsImN1c3RvbWVyTnVtYmVyIjoiMDAwMTE5OTQ0MTYwIiwidXNlciI6ImNkdWZmeUB0ZmQuaWUifQ.fisOWBMqnbzzcNQpqO6Cmu6DEMjroaZYgTsAeEmR36A" }`
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/reports/4324324324", bytes.NewBuffer([]byte(requestPayload)))
		conn := connectors.NewTestConnectors("../../tests/payload.json", STATUS, logger)
		conn.SetMode("push")
		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"key": "test",
		}
		conn.Meta("true")
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

}
