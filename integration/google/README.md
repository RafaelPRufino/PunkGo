# Punk\Go\Integration\Google

Integração com Google API. 

```go
package main

import (
	"fmt"
	option "github.com/RafaelPRufino/PunkGo/integration/google/auth/option"
	google "github.com/RafaelPRufino/PunkGo/integration/google/client"
)

func main() {
	// New client google with Google Credentials Application
	client := google.NewClient(option.WithGoogleCredentialsFile("keys/google-keys.json"))

	//Authenticate in Google API with User Credentials
	client.SignIn(option.WithOAuthCredentialsFile("keys/access-token.json"))

	//Create new Analytics Service
	ga := client.NewAnalyticsService()
	//Create new Request Analytics Service with ViewId
	request := ga.CreateRequest("555555555")
	//  Add Dimensions Reports
	request.AddDimension("ga:browser")
	//  Add Metrics Reports
	request.AddMetric("ga:users")
	request.AddMetric("ga:newUsers")
	request.AddMetric("ga:newUsers")
	//Process report
	response, err := ga.Do(request)

	if err != nil {
		fmt.Println("Error:" +" "+ err.Error())
	} else {
		//Print Data JSON
		json, err:=response.MarshalJSON()
		fmt.Println(string(json) , err)
	}
}
```