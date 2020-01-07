package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	// "github.com/rs/xid"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	CONTENTTYPE     string = "Content-Type"
	APPLICATIONJSON string = "application/json"
)

// Response schema
type Response struct {
	Name       string `json:"name"`
	StatusCode string `json:"statuscode"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	Payload    string `json:"payload"`
}

type Claims struct {
	jwt.StandardClaims
}

func AnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	var response Response
	var analytics *Analytics

	addHeaders(w, r)

	body, _ := ioutil.ReadAll(r.Body)
	// we first unmarshal the payload and add needed values before writing to couchbase
	errs := json.Unmarshal(body, &analytics)
	if errs != nil {
		logger.Error(fmt.Sprintf("Could not unmarshal message data to schema %v", errs))
		response = Response{Name: os.Getenv("NAME"), StatusCode: "500", Status: "ERROR", Message: fmt.Sprintf("Could not unmarshal message data to schema %v", errs)}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		analytics.Product = "Trackmate"

		// ensure uniqueness
		// id := xid.New().String()
		id := analytics.TrackingId + analytics.From.PageName + strconv.FormatInt(analytics.Timestamp, 10)
		_, err := bucketClient.Upsert(id, analytics, 0)
		if err != nil {
			logger.Error(fmt.Sprintf("Could not insert schema into couchbase %v", err))
			response = Response{Name: os.Getenv("NAME"), StatusCode: "500", Status: "ERROR", Message: fmt.Sprintf("Could not insert schema into couchbase %v", errs)}
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			// all good :)
			logger.Debug(fmt.Sprintf("Analytics schema inserted into couchbase  %v \n", analytics))
			response = Response{Name: os.Getenv("NAME"), StatusCode: "200", Status: "OK", Message: "Data inserted succesfully", Payload: string(body)}
			w.WriteHeader(http.StatusOK)
		}
	}
	b, _ := json.MarshalIndent(response, "", "	")
	logger.Debug(fmt.Sprintf("AnatylicsHandler response : %s", string(b)))
	fmt.Fprintf(w, string(b))
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	var response Response

	token := r.Header.Get(strings.ToLower("Authorization"))
	// Remove Bearer
	tknStr := strings.Trim(token[7:], " ")
	logger.Debug(fmt.Sprintf("Token : %s", tknStr))
	addHeaders(w, r)

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	decodedSecret, _ := base64.StdEncoding.DecodeString(os.Getenv("JWT_SECRETKEY"))

	var jwtKey = []byte(decodedSecret)

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err.Error() == jwt.ErrSignatureInvalid.Error() {
			w.WriteHeader(http.StatusUnauthorized)
			response = Response{Name: os.Getenv("NAME"), StatusCode: "403", Status: "ERROR", Message: "Forbidden"}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			response = Response{Name: os.Getenv("NAME"), StatusCode: "400", Status: "ERROR", Message: "Bad Request"}
		}
	} else {
		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			response = Response{Name: os.Getenv("NAME"), StatusCode: "403", Status: "ERROR", Message: "Forbidden"}
		} else {
			response = Response{Name: os.Getenv("NAME"), StatusCode: "200", Status: "OK", Message: "Data uploaded succesfully", Payload: "Access granted"}
			w.WriteHeader(http.StatusOK)
		}
	}
	b, _ := json.MarshalIndent(response, "", "	")
	logger.Debug(fmt.Sprintf("SimpleAuthHandler response : %s", string(b)))
	fmt.Fprintf(w, string(b))
}

func IsAlive(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{ \"version\" : \"1.0.2\" , \"name\": \""+os.Getenv("NAME")+"\" }")
}

// headers (with cors) utility
func addHeaders(w http.ResponseWriter, r *http.Request) {
	var request []string
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	logger.Trace(fmt.Sprintf("Headers : %s", request))

	w.Header().Set(CONTENTTYPE, APPLICATIONJSON)
	// use this for cors
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

}
