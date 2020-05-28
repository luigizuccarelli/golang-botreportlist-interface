package schema

// Response schema
type Response struct {
	Name       string `json:"name"`
	StatusCode string `json:"statuscode"`
	Status     string `json:"status"`
	Message    string `json:"message"`
}

type AudienceSchema struct {
	Data struct {
		Segments    []string `json:"segments"`
		SegmentsAll []string `json:"segments_all"`
	} `json:"data"`
	Message string `json:"message"`
	Meta    struct {
		ByFields  []string    `json:"by_fields"`
		Conflicts interface{} `json:"conflicts"`
		Format    string      `json:"format"`
		Name      string      `json:"name"`
	} `json:"meta"`
	Status int `json:"status"`
}
