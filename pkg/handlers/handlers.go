package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/trackmate-couchbase-push/pkg/connectors"
	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/trackmate-couchbase-push/pkg/schema"
	gocb "github.com/couchbase/gocb/v2"
	"github.com/microlib/simple"
)

const (
	CONTENTTYPE     string = "Content-Type"
	APPLICATIONJSON string = "application/json"
)

func AnalyticsHandler(w http.ResponseWriter, r *http.Request, logger *simple.Logger, con connectors.Clients) {
	var response *schema.Response
	var analytics *schema.SegmentIO

	addHeaders(w, r)

	body, _ := ioutil.ReadAll(r.Body)
	// we first unmarshal the payload and add needed values before writing to couchbase
	errs := json.Unmarshal(body, &analytics)
	if errs != nil {
		logger.Error(fmt.Sprintf("Could not unmarshal message data to schema %v", errs))
		response = &schema.Response{Name: os.Getenv("NAME"), StatusCode: "500", Status: "ERROR", Message: fmt.Sprintf("Could not unmarshal message data to schema %v", errs)}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		// ensure uniqueness
		analytics.Id = analytics.MessageID
		// get a collection reference
		upsertResult, err := con.Upsert(analytics.Id, analytics, &gocb.UpsertOptions{})

		if err != nil {
			logger.Error(fmt.Sprintf("Could not insert schema into couchbase %v", err))
			response = &schema.Response{Name: os.Getenv("NAME"), StatusCode: "500", Status: "ERROR", Message: fmt.Sprintf("Could not insert schema into couchbase %v", errs)}
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			// all good :)
			logger.Debug(fmt.Sprintf("Analytics schema inserted into couchbase  %v \n", analytics))
			response = &schema.Response{Name: os.Getenv("NAME"), StatusCode: "200", Status: "OK", Message: "Data inserted succesfully", Payload: upsertResult}
			w.WriteHeader(http.StatusOK)
		}
	}
	b, _ := json.MarshalIndent(response, "", "	")
	logger.Debug(fmt.Sprintf("AnatylicsHandler response : %s", string(b)))
	fmt.Fprintf(w, string(b))
}

func IsAlive(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{ \"version\" : \""+os.Getenv("VERSION")+"\" , \"name\": \""+os.Getenv("NAME")+"\" }")
}

// headers (with cors) utility
func addHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(CONTENTTYPE, APPLICATIONJSON)
	// use this for cors
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
