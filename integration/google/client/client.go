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

type client struct {
	GoogleCredentials auth.GoogleCredentials
	RefreshToken      auth.OAuthCredentials
	AccessToken       auth.OAuthCredentials
}

func (receiver *client) GetGoogleCredentials() auth.GoogleCredentials {
	return receiver.GoogleCredentials
}

func (receiver *client) GetAccessToken() auth.OAuthCredentials {
	return receiver.AccessToken
}

func (receiver *client) IsAuthenticated() bool {
	if receiver.RefreshToken.IsExpired() != false {
		receiver.AccessToken = receiver.RefreshToken
	}

	if receiver.AccessToken.IsExpired() {
		oAuth2 := receiver.NewOAuth2Service()
		receiver.AccessToken = oAuth2.GetNewToken(receiver.RefreshToken)
	}

	return receiver.AccessToken.IsExpired() == false
}

func (receiver *client) NewOAuth2Service() oauth2.GoogleOAuth2 {
	return oauth2.NewOAuth2Provider(receiver)
}

func (receiver *client) NewAnalyticsService() analytics.GoogleAnalytics {
	return analytics.NewAnalyticsProvider(receiver)
}

func (receiver *client) SignIn(With auth.OAuthCredentials) bool {
	receiver.RefreshToken = With
	return receiver.IsAuthenticated()
}

func NewClient(credentials auth.GoogleCredentials) GoogleClient {
	var client client = client{}
	client.GoogleCredentials = credentials
	return &client
}
