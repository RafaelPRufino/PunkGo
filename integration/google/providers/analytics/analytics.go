package analytics

import (
	"bytes"
	"fmt"
	"github.com/RafaelPRufino/PunkGo/integration/google/providers"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	AnalyticsBatchGET = "https://analyticsreporting.googleapis.com/v4/reports:batchGet"
)

type ClientForAnalytics interface {
	providers.ProviderTarget
}

type GoogleAnalytics interface {
	CreateRequest(ViewId string) Request
	Do(request Request) (Response, error)
}

type analytics struct {
	Client ClientForAnalytics
}

func (receiver *analytics) CreateRequest(ViewId string) Request {
	return Request{ViewId, []RequestMetric{}, []RequestDimension{}}
}

func (receiver *analytics) Do(request Request) (Response, error) {
	response := Response{}
	if receiver.Client.IsAuthenticated() == false {
		NewError("Unauthorized")
	}

	json, err := request.MarshalJSON()
	if err != nil {
		return response, err
	}

	googleResponse, err := receiver.doReport(json)
	if err != nil {
		return response, err
	}

	err = response.absorb(googleResponse)
	return response, err
}

func (receiver *analytics) doReport(report []byte) (GoogleResponseReport, error) {
	result := GoogleResponseReport{}

	endpoint := AnalyticsBatchGET
	accessToken := receiver.Client.GetAccessToken()
	client := &http.Client{}

	request, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(report)) // URL-encoded payload
	if err != nil {
		return GoogleResponseReport{}, err
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", accessToken.Authorization())

	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return GoogleResponseReport{}, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return GoogleResponseReport{}, err
	}

	err = result.UnmarshalJSON(body)
	if err != nil {
		descError := string(body)
		if strings.Contains(descError, "UNAUTHENTICATED") {
			err = NewError("Unauthenticated")
		} else if strings.Contains(descError, "INVALID_ARGUMENT") {
			err = NewError("MalformedRequest")
		}
	}

	return result, err
}

func NewAnalyticsProvider(client ClientForAnalytics) GoogleAnalytics {
	var ga = analytics{}
	ga.Client = client
	return &ga
}
