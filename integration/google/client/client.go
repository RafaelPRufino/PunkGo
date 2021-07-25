package client

import (
	"github.com/RafaelPRufino/PunkGo/integration/google/auth"
	"github.com/RafaelPRufino/PunkGo/integration/google/providers/analytics"
	"github.com/RafaelPRufino/PunkGo/integration/google/providers/oauth2"
)

type GoogleClient interface {
	SignIn(With auth.OAuthCredentials) bool
	IsAuthenticated() bool
	NewAnalyticsService() analytics.GoogleAnalytics
}

type Client struct {
	GoogleCredentials auth.GoogleCredentials
	RefreshToken      auth.OAuthCredentials
	AccessToken       auth.OAuthCredentials
}

func (receiver *Client) GetGoogleCredentials() auth.GoogleCredentials {
	return receiver.GoogleCredentials
}

func (receiver *Client) GetAccessToken() auth.OAuthCredentials {
	return receiver.AccessToken
}

func (receiver *Client) IsAuthenticated() bool {
	if receiver.RefreshToken.IsExpired() != false {
		receiver.AccessToken = receiver.RefreshToken
	}

	if receiver.AccessToken.IsExpired() {
		oAuth2 := receiver.CreateOAuth2Provider()
		receiver.AccessToken = oAuth2.GetNewToken(receiver.RefreshToken)
	}

	return receiver.AccessToken.IsExpired() == false
}

func (receiver *Client) CreateOAuth2Provider() oauth2.OAuth2 {
	return oauth2.NewOAuth2Provider(receiver)
}

func (receiver *Client) NewAnalyticsService() analytics.GoogleAnalytics {
	return analytics.NewAnalyticsProvider(receiver)
}

func (receiver *Client) SignIn(With auth.OAuthCredentials) bool {
	receiver.RefreshToken = With
	return receiver.IsAuthenticated()
}

func NewClient(credentials auth.GoogleCredentials) GoogleClient {
	var client Client = Client{}
	client.GoogleCredentials = credentials
	return &client
}
