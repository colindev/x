package secret

import (
	"encoding/json"
	"io/ioutil"
)

type Config interface {
	ReadFile(f string) error
}

type Service struct {
	Type         string `json:"type"`
	ProjectID    string `json:"project_id"`
	PrivateKeyID string `json:"project_key_id"`
	PrivateKey   string `json:"private_key"`
	ClientEmail  string `json:"client_email"`
	ClientID     string `json:"client_id"`
	AuthURI      string `json:"auth_uri"`
	TokenURI     string `json:"token_uri"`
}

func (ser *Service) ReadFile(f string) error {

	b, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, ser)
	if err != nil {
		return err
	}

	return nil
}

type OAuth2 struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	ProjectID    string `json:"project_id"`
	AuthURI      string `json:"auth_uri"`
	TokenURI     string `json:"token_uri"`
}

func (ser *OAuth2) ReadFile(f string) error {

	b, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &struct{ Web *OAuth2 }{ser})
	if err != nil {
		return err
	}

	return nil
}
