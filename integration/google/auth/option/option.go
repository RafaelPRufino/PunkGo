package option

import (
	"errors"
	"github.com/RafaelPRufino/PunkGo/integration/google/auth"
	"io/ioutil"
	"os"
)

func WithFile(filename string) ([]byte, error) {
	bytes, err := os.Open(filename)

	if err != nil {
		return []byte{}, errors.New("Error opening file {" + filename + "}")
	}

	values, err := ioutil.ReadAll(bytes)
	if err != nil {
		return []byte{}, errors.New("Error reading file{" + filename + "}")
	}
	return values, nil
}

// WithGoogleCredentialsFile Google Credentials file JSON
func WithGoogleCredentialsFile(filename string) auth.GoogleCredentials {
	gooCred := auth.GoogleCredentials{}
	bytes, err := WithFile(filename)

	if err == nil {
		err = gooCred.UnmarshalJSON(bytes)
	}

	return gooCred
}

// WithOAuthCredentialsFile OAuth Credentials file JSON
func WithOAuthCredentialsFile(filename string) auth.OAuthCredentials {
	authCred := auth.OAuthCredentials{}
	bytes, err := WithFile(filename)

	if err == nil {
		err = authCred.UnmarshalJSON(bytes)
	}

	return authCred
}
