package analytics

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/RafaelPRufino/PunkGo/integration/google/providers"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	AnalyticsScope = "https://analyticsreporting.googleapis.com/v4/reports:batchGet"
)

type ClientForAnalytics interface {
	providers.ProviderTarget
}

type GoogleAnalytics interface {
	CreateRequest(ViewId string) Request
	Do(request Request) (Response, error)
}

type Analytics struct {
	Client ClientForAnalytics
}

func NewAnalyticsProvider(client ClientForAnalytics) GoogleAnalytics {
	var ga = Analytics{}
	ga.Client = client
	return &ga
}

func (receiver *Analytics) CreateRequest(ViewId string) Request {
	return Request{ViewId, []RequestMetric{}, []RequestDimension{}}
}

func (receiver *Analytics) Do(request Request) (Response, error) {
	response := Response{}
	if receiver.Client.IsAuthenticated() == false {
		errors.New("unauthenticated")
	}

	json, err := request.MarshalJSON()
	if err != nil {
		return response, err
	}

	googleResponse, err := receiver.DoReport(json)
	if err != nil {
		return response, err
	}

	err = response.Absorb(googleResponse)
	return response, err
}

func (receiver *Analytics) DoReport(report []byte) (GoogleResponseReport, error) {
	result := GoogleResponseReport{}

	endpoint := AnalyticsScope
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
			err = NewError("UNAUTHENTICATED")
		} else if strings.Contains(descError, "INVALID_ARGUMENT") {
			err = NewError("MalformedRequest")
		}
	}

	return result, err
}
