package oauth2

import (
	"github.com/RafaelPRufino/PunkGo/integration/google/auth"
	"github.com/RafaelPRufino/PunkGo/integration/google/providers"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type ClientForOAuth2 interface {
	providers.ProviderTarget
}

type GoogleOAuth2 interface {
	GetNewToken(refreshToken  auth.OAuthCredentials ) auth.OAuthCredentials
}

type oAuth2 struct {
	Client ClientForOAuth2
}

func NewOAuth2Provider(client ClientForOAuth2) GoogleOAuth2 {
	return &oAuth2{client}
}

func (receiver *oAuth2) GetNewToken(refreshToken  auth.OAuthCredentials ) auth.OAuthCredentials {
	credentials := auth.OAuthCredentials{}
	google := receiver.Client.GetGoogleCredentials()

	endpoint := google.TokenUri
	data := url.Values{}

	data.Set("client_id", google.Id)
	data.Set("client_secret", google.SecretKey)
	data.Set("refresh_token", refreshToken.RefreshToken)
	data.Set("grant_type", "refresh_token")

	client := &http.Client{}
	request, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil {
		return credentials
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	response, err := client.Do(request)
	if err != nil {
		return credentials
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return credentials
	}

	err = credentials.UnmarshalJSON(body)
	if err == nil {
		credentials.CreatedNow()
	}

	return credentials
}

