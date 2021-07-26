package analytics

import (
	"encoding/json"
	"errors"
	"github.com/RafaelPRufino/PunkGo/support"
	"github.com/RafaelPRufino/PunkGo/support/exception/scope"
	"strings"
	"time"
)

type Response struct {
	Data []ResponseData
}

type ResponseData struct {
	Created    int64              `json:"created"`
	Dimensions ResponseDimensions `json:"dimensions"`
	Metrics    ResponseMetrics    `json:"metrics"`
}

type ResponseDimensions struct {
	support.Attributes
}

type ResponseMetrics struct {
	support.Attributes
}

// MarshalJSON custom parse JSON
func (response *Response) MarshalJSON() ([]byte, error) {
	type NoMethod ResponseData
	data := response.Data
	return json.Marshal(&data)
}

// Add adds the value ResponseData
func (response *Response) Add(value ResponseData) {
	response.Data = append(response.Data, value)
}

// absorb adds the value to key. It appends to any existing
func (response *Response) absorb(reports GoogleResponseReport) error {
	for _, report := range reports.Reports {
		header := report.ColumnHeader
		dimHdqrs := header.Dimensions
		metricHdqrs := header.MetricHeader.MetricHeaderEntries
		rows := report.Data.Rows

		for _, row := range rows {
			dims := row.Dimensions
			metrics := row.Metrics

			data := ResponseData{reports.Created, ResponseDimensions{}, ResponseMetrics{}}
			for i := 0; i < len(dimHdqrs) && i < len(dims); i++ {
				name := strings.Replace(dimHdqrs[i], "ga:", "", 1)
				data.Dimensions.Add(name, dims[i])
			}

			data.Metrics = ResponseMetrics{}
			for _, metric := range metrics {
				for j := 0; j < len(metricHdqrs) && j < len(metric.Values); j++ {
					name := strings.Replace(metricHdqrs[j].Name, "ga:", "", 1)
					data.Metrics.Add(name, metric.Values[j])
				}
			}
			response.Add(data)
		}
	}
	return nil
}

// GoogleResponseReport Response Google analytics
type GoogleResponseReport struct {
	Created int64 `json:"created"`
	Reports []struct {
		ColumnHeader struct {
			Dimensions   []string `json:"dimensions,omitempty"`
			MetricHeader struct {
				MetricHeaderEntries []struct {
					Name string `json:"name"`
					Type string `json:"type"`
				} `json:"metricHeaderEntries"`
			} `json:"metricHeader"`
		} `json:"columnHeader"`
		Data struct {
			Rows []struct {
				Dimensions []string `json:"dimensions,omitempty"`

				Metrics []*struct {
					Values []interface{} `json:"values,omitempty"`
				} `json:"metrics,omitempty"`
			} `json:"rows"`
		} `json:"data"`
	} `json:"reports"`
}

// UnmarshalJSON JSON to GoogleResponseReport
func (response *GoogleResponseReport) UnmarshalJSON(data []byte) error {
	var err error
	scope.Observable{
		Try: func() {
			type NoMethod GoogleResponseReport
			var result struct {
				*NoMethod
			}
			err = json.Unmarshal([]byte(data), &result)
			if err == nil {
				response.Reports = result.Reports
				response.Created = time.Now().Unix()
			}
		},
		Catch: func(e scope.Exception) {
			if err == nil {
				err = errors.New("conversion failed")
			}
		},
		Finally: func() {

		},
	}.Do()
	return err
}
