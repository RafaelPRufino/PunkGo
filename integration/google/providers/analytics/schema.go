package analytics

import (
	"strings"
)

type ReportSchema struct {
	Dimensions []ResponseSchemaType `json:"dimensions"`
	Metrics    []ResponseSchemaType `json:"metrics"`
}

type ResponseSchemaType struct {
	FullName string `json:"full_name"`
	Name     string `json:"name"`
	Type     string `json:"type"`
}

func caseType(nameType string) string {
	types := map[string]string{
		"CURRENCY": "FLOAT",
		"PERCENT":  "FLOAT",
	}
	taper := types[nameType]
	if taper == "" {
		taper = nameType
	}
	return taper
}

func writerSchema(header GoogleResponseReportHeader, writer Writer) {
	dimensionsHeader := header.Dimensions
	metricsHeader := header.MetricHeader.MetricHeaderEntries
	schema := ReportSchema{}

	if writer == nil {
		writer = &DefaultReportWriter{}
	}

	for i := 0; i < len(dimensionsHeader) && i < len(dimensionsHeader); i++ {
		name := strings.Replace(dimensionsHeader[i], "ga:", "", 1)
		dimension := ResponseSchemaType{"dimensions_header." + name, name, "string"}
		schema.Dimensions = append(schema.Dimensions, dimension)
	}

	for i := 0; i < len(metricsHeader) && i < len(metricsHeader); i++ {
		name := strings.Replace(metricsHeader[i].Name, "ga:", "", 1)
		metric := ResponseSchemaType{"dimensions_header." + name, name, caseType(metricsHeader[i].Type)}
		schema.Metrics = append(schema.Metrics, metric)
	}

	writer.SetSchema(schema)
}
