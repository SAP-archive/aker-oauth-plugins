package handler

import (
	"fmt"

	"github.com/SAP/aker/plugin"
)

const sessionStoreKeyLength = 32

type Config struct {
	OAuth   OAuthConfig   `yaml:"oauth"`
	Session SessionConfig `yaml:"session"`
}

type OAuthConfig struct {
	ClientID          string   `yaml:"client_id"`
	ClientSecret      string   `yaml:"client_secret"`
	SkipSSLValidation bool     `yaml:"skip_ssl_validation"`
	AuthorizationURL  string   `yaml:"authorization_url"`
	TokenURL          string   `yaml:"token_url"`
	RedirectURL       string   `yaml:"redirect_url"`
	RequiredScopes    []string `yaml:"required_scopes"`
	OptionalScopes    []string `yaml:"optional_scopes"`
}

type SessionConfig struct {
	AuthenticationKey string `yaml:"authentication_key"`
	EncryptionKey     string `yaml:"encryption_key"`
}

func ParseConfig(data []byte) (*Config, error) {
	cfg := &Config{}
	if err := plugin.UnmarshalConfig(data, cfg); err != nil {
		return nil, err
	}

	authKey := []byte(cfg.Session.AuthenticationKey)
	if len(authKey) != sessionStoreKeyLength {
		return nil, &InvalidKeyLengthError{
			key:     authKey,
			keyType: "authentication",
		}
	}
	encKey := []byte(cfg.Session.EncryptionKey)
	if len(encKey) != sessionStoreKeyLength {
		return nil, &InvalidKeyLengthError{
			key:     encKey,
			keyType: "encryption",
		}
	}
	return cfg, nil
}

type InvalidKeyLengthError struct {
	key     []byte
	keyType string
}

func (e *InvalidKeyLengthError) Error() string {
	return fmt.Sprintf("Invalid %s key of length %d, expected %d",
		e.keyType, len(e.key), sessionStoreKeyLength)
}
