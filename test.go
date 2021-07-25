package main

import (
	"fmt"
	"github.com/RafaelPRufino/PunkGo/integration/google/auth/option"
	google "github.com/RafaelPRufino/PunkGo/integration/google/client"
)

func main() {
	client := google.NewClient(option.WithGoogleCredentialsFile("keys/google-keys.json"))

	client.SignIn(option.WithOAuthCredentialsFile("keys/access-token.json"))

	ga := client.NewAnalyticsService()
	request := ga.CreateRequest("230920570")
	request.AddDimension("ga:browser")
	request.AddMetric("ga:users")
	request.AddMetric("ga:newUsers")
	request.AddMetric("ga:newUsers")
	response, err := ga.Do(request)

	fmt.Println(response, err)

}
