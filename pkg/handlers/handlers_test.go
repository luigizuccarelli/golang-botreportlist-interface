package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/trackmate-lytics-interface/pkg/connectors"
	"github.com/gorilla/mux"
	"github.com/microlib/simple"
)

func TestAllMiddleware(t *testing.T) {

	logger := &simple.Logger{Level: "info"}

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

	t.Run("ProfileHandler : should pass", func(t *testing.T) {
		var STATUS int = 200
		os.Setenv("TOKEN", "[{ \"id\": 1, \"name\": \"BH-01\", \"token\": \"1212121\"}]")
		os.Setenv("TESTING", "true")
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/BH-01/profile/test@test.com", nil)
		conn := connectors.NewTestConnectors("../../tests/payload.json", STATUS, logger)

		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"affiliateid": "BH-01",
			"email":       "test@test.com",
		}

		req = mux.SetURLVars(req, vars)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ProfileHandler(w, r, conn)
		})

		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "ProfileHandler", rr.Code, STATUS))
		}
	})

	t.Run("ProfileHandler : should pass", func(t *testing.T) {
		var STATUS int = 200
		os.Setenv("TOKEN", "[{ \"id\": 1, \"name\": \"BH-01\", \"token\": \"1212121\"}]")
		os.Setenv("TESTING", "true")
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/BH-01/profile/test@test.com", nil)
		conn := connectors.NewTestConnectors("../../tests/payload.json", STATUS, logger)

		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"affiliateid": "BH-01",
			"email":       "test@test.com",
		}

		req = mux.SetURLVars(req, vars)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ProfileHandler(w, r, conn)
		})

		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "ProfileHandler", rr.Code, STATUS))
		}
	})

	t.Run("ProfileHandler : should fail (invalid token)", func(t *testing.T) {
		var STATUS int = 500
		os.Setenv("TOKEN", "[{ \"id\": 1, \"name\": \"BX-01\", \"token\": \"1212121\"}]")
		os.Setenv("TESTING", "false")
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/BH-01/profile/test@test.com", nil)
		conn := connectors.NewTestConnectors("../../tests/payload.json", STATUS, logger)

		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"affiliateid": "BH-01",
			"email":       "test@test.com",
		}

		req = mux.SetURLVars(req, vars)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ProfileHandler(w, r, conn)
		})

		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "ProfileHandler", rr.Code, STATUS))
		}
	})

	t.Run("ProfileHandler : should fail email (json unmarshal)", func(t *testing.T) {
		var STATUS int = 500
		os.Setenv("TOKEN", "[{ \"id\": 1, \"name\": \"BH-01\", \"token\": \"1212121\"}]")
		os.Setenv("TESTING", "false")
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/BH-01/profile/test@test.com", nil)
		conn := connectors.NewTestConnectors("../../tests/payload.json", STATUS, logger)

		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"affiliateid": "BH-01",
			"email":       "test@test.com",
		}

		req = mux.SetURLVars(req, vars)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ProfileHandler(w, r, conn)
		})

		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "ProfileHandler", rr.Code, STATUS))
		}
	})

	t.Run("ProfileHandler : should fail profile (json unmarshal)", func(t *testing.T) {
		var STATUS int = 500
		os.Setenv("TOKEN", "[{ \"id\": 1, \"name\": \"BH-01\", \"token\": \"1212121\"}]")
		os.Setenv("TESTING", "false")
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/BH-01/profile/test@test.com", nil)
		conn := connectors.NewTestConnectors("../../tests/payload.json", STATUS, logger)

		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"affiliateid": "BH-01",
			"email":       "test@test.com",
		}

		req = mux.SetURLVars(req, vars)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ProfileHandler(w, r, conn)
		})

		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "ProfileHandler", rr.Code, STATUS))
		}
	})

	t.Run("ProfileHandler : should fail (forced request error)", func(t *testing.T) {
		var STATUS int = 500
		os.Setenv("TOKEN", "[{ \"id\": 1, \"name\": \"BH-01\", \"token\": \"1212121\"}]")
		os.Setenv("TESTING", "false")
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/BH-01/profile/test@test.com", nil)
		conn := connectors.NewTestConnectors("../../tests/payload.json", STATUS, logger)

		//Hack to try to fake gorilla/mux vars
		vars := map[string]string{
			"affiliateid": "BH-01",
			"email":       "test@test.com",
		}

		req = mux.SetURLVars(req, vars)
		conn.Meta("true")
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ProfileHandler(w, r, conn)
		})

		handler.ServeHTTP(rr, req)

		body, e := ioutil.ReadAll(rr.Body)
		if e != nil {
			t.Fatalf("Should not fail : found error %v", e)
		}
		logger.Trace(fmt.Sprintf("Response %s", string(body)))
		// ignore errors here
		if rr.Code != STATUS {
			t.Errorf(fmt.Sprintf("Handler %s returned with incorrect status code - got (%d) wanted (%d)", "ProfileHandler", rr.Code, STATUS))
		}
	})

}
