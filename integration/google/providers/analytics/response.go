package analytics

import (
	"encoding/json"
	"errors"
	"github.com/RafaelPRufino/PunkGo/support"
	"github.com/RafaelPRufino/PunkGo/support/exception/scope"
	"strings"
	"time"
)

type Writer interface {
	Create() DataRow
	Add(row interface{})
	SetSchema(schema ReportSchema)
}

type DataRow interface {
	AddMetric(key string, value interface{})
	AddDimension(key string, value interface{})
	SetCreated(created int64)
}

func writerReport(reports GoogleResponseReports, writer Writer) error {
	for _, report := range reports.Reports {
		writerData(reports.Created, report.Header, report.Data, writer)
		writerSchema(report.Header, writer)
	}
	return nil
}

func writerData(created int64, header GoogleResponseReportHeader, data GoogleResponseReportData, writer Writer) {
	dimensionsHeader := header.Dimensions
	metricsHeader := header.MetricHeader.MetricHeaderEntries
	rows := data.Rows

	if writer == nil {
		writer = &DefaultReportWriter{}
	}

	for _, row := range rows {
		dimensions := row.Dimensions
		metrics := row.Metrics

		data := writer.Create()
		data.SetCreated(created)
		for i := 0; i < len(dimensionsHeader) && i < len(dimensions); i++ {
			name := strings.Replace(dimensionsHeader[i], "ga:", "", 1)
			data.AddDimension(name, dimensions[i])
		}

		for _, metric := range metrics {
			for j := 0; j < len(metricsHeader) && j < len(metric.Values); j++ {
				name := strings.Replace(metricsHeader[j].Name, "ga:", "", 1)
				data.AddMetric(name, metric.Values[j])
			}
		}
		writer.Add(data)
	}
}

// DefaultReportWriter writer report
type DefaultReportWriter struct {
	Data   []DefaultDataRow
	Schema ReportSchema `json:"schema"`
}

// Create create new DataRow for insert
func (writer *DefaultReportWriter) Create() DataRow {
	return &DefaultDataRow{}
}

// Add add new row report
func (writer *DefaultReportWriter) Add(row interface{}) {
	writer.Data = append(writer.Data, *row.(*DefaultDataRow))
}

// SetSchema set schema report
func (writer *DefaultReportWriter) SetSchema(schema ReportSchema) {
	writer.Schema = schema
}

// MarshalJSON custom parse JSON
func (writer *DefaultReportWriter) MarshalJSON() ([]byte, error) {
	data := writer.Data
	return json.Marshal(&data)
}

// DefaultDataRow data row default
type DefaultDataRow struct {
	Created    int64              `json:"created"`
	Dimensions support.Attributes `json:"dimensions"`
	Metrics    support.Attributes `json:"metrics"`
}

func (row *DefaultDataRow) AddMetric(key string, value interface{}) {
	row.Metrics.Add(key, value)
}

func (row *DefaultDataRow) AddDimension(key string, value interface{}) {
	row.Metrics.Add(key, value)
}

func (row *DefaultDataRow) SetCreated(created int64) {
	row.Created = created
}

// GoogleResponseReports Response Google analytics
type GoogleResponseReports struct {
	Created int64                  `json:"created"`
	Reports []GoogleResponseReport `json:"reports"`
}
type GoogleResponseReport struct {
	Header GoogleResponseReportHeader `json:"columnHeader"`
	Data   GoogleResponseReportData   `json:"data"`
}
type GoogleResponseReportHeader struct {
	Dimensions   []string `json:"dimensions,omitempty"`
	MetricHeader struct {
		MetricHeaderEntries []struct {
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"metricHeaderEntries"`
	} `json:"metricHeader"`
}
type GoogleResponseReportData struct {
	Rows []struct {
		Dimensions []string `json:"dimensions,omitempty"`
		Metrics    []*struct {
			Values []interface{} `json:"values,omitempty"`
		} `json:"metrics,omitempty"`
	} `json:"rows"`
}

// UnmarshalJSON JSON to GoogleResponseReports
func (response *GoogleResponseReports) UnmarshalJSON(data []byte) error {
	var err error
	scope.Observable{
		Try: func() {
			type NoMethod GoogleResponseReports
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
