package providers

import "github.com/RafaelPRufino/PunkGo/integration/google/auth"

type ProviderTarget interface  {
	GetGoogleCredentials() auth.GoogleCredentials
	GetAccessToken() auth.OAuthCredentials
	IsAuthenticated() bool
}

