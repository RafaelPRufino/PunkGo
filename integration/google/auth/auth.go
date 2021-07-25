package auth

import (
	"encoding/json"
	"time"
)

type OAuthCredentials struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	Created      int64  `json:"created"`
}

// IsExpired Notify if expired
func (receiver *OAuthCredentials) IsExpired() bool {
	expiresIn := receiver.Created + (receiver.ExpiresIn - 30)
	return expiresIn < time.Now().Unix()
}

// CreatedNow Receiver new Credentials
func (receiver *OAuthCredentials) CreatedNow() {
	receiver.Created = time.Now().Unix()
}
// UnmarshalJSON Load OAuth Credentials from JSON
func (receiver *OAuthCredentials) UnmarshalJSON(data []byte) error {
	type NoMethod OAuthCredentials
	var result struct {
		*NoMethod
	}

	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		return err
	}

	receiver.AccessToken = result.AccessToken
	receiver.ExpiresIn = result.ExpiresIn
	receiver.Scope = result.Scope
	receiver.TokenType = result.TokenType
	receiver.IdToken = result.IdToken
	receiver.RefreshToken = result.RefreshToken
	receiver.Created = result.Created
	return nil
}
// Authorization Get Authorization
func (receiver *OAuthCredentials) Authorization() string {
	return receiver.TokenType + " " + receiver.AccessToken
}

// GoogleCredentials Google Credentials App
type GoogleCredentials struct {
	Id                      string `json:"client_id"`
	ProjectId               string `json:"project_id"`
	AuthUri                 string `json:"auth_uri"`
	TokenUri                string `json:"token_uri"`
	AuthProviders509century string `json:"auth_provider_x509_cert_url"`
	SecretKey               string `json:"client_secret"`
}

// UnmarshalJSON Load Google Credentials from JSON
func (receiver *GoogleCredentials) UnmarshalJSON(data []byte) error {
	type NoMethod GoogleCredentials
	var result struct {
		Web *NoMethod `json:"web"`
	}
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		return err
	}

	receiver.Id = result.Web.Id
	receiver.ProjectId = result.Web.ProjectId
	receiver.AuthUri = result.Web.AuthUri
	receiver.AuthProviders509century = result.Web.AuthProviders509century
	receiver.SecretKey = result.Web.SecretKey
	receiver.TokenUri = result.Web.TokenUri
	return nil
}
