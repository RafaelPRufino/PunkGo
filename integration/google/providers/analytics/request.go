package analytics

import "encoding/json"

// Request Request Google analytics
type Request struct {
	Id         string             `json:"viewId"`
	Metrics    []RequestMetric    `json:"metrics"`
	Dimensions []RequestDimension `json:"dimensions"`
}

type RequestMetric struct {
	Expression string `json:"expression"`
}

type RequestDimension struct {
	Name string `json:"name"`
}

func (request *Request) AddMetric(expression string) {
	request.Metrics = append(request.Metrics, RequestMetric{expression})
}

func (request *Request) AddDimension(name string) {
	request.Dimensions = append(request.Dimensions, RequestDimension{name})
}

// MarshalJSON custom parse JSON
func (request *Request) MarshalJSON() ([]byte, error) {
	type NoMethod Request
	report := struct {
		ReportRequests []NoMethod `json:"reportRequests"`
	}{}
	report.ReportRequests = []NoMethod{NoMethod(*request)}
	return json.Marshal(report)
}
