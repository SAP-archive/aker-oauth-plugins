package token

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"golang.org/x/oauth2"
)

//go:generate counterfeiter . Decoder

type Decoder interface {
	Decode(token *oauth2.Token) (Info, error)
}

type decoder struct {
}

func NewDecoder() Decoder {
	return &decoder{}
}

func (d *decoder) Decode(token *oauth2.Token) (Info, error) {
	segments := strings.Split(token.AccessToken, ".")
	if len(segments) != 3 {
		return Info{}, fmt.Errorf("Token should have three segments.")
	}
	payload := segments[1]

	jsonPayload, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return Info{}, err
	}

	info := Info{}
	if err := json.Unmarshal([]byte(jsonPayload), &info); err != nil {
		return Info{}, err
	}
	return info, nil
}
