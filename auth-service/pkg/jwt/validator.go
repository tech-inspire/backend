package jwt

import (
	"encoding/json"
	"fmt"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

type Validator struct {
	keyFunc jwt.Keyfunc
}

func NewValidatorFromKeyFunc(keyFunc jwt.Keyfunc) *Validator {
	return &Validator{
		keyFunc: keyFunc,
	}
}

func NewValidatorFromURL(jwksUrl string) (*Validator, error) {
	var (
		kf  keyfunc.Keyfunc
		err error
	)

	if json.Valid([]byte(jwksUrl)) {
		kf, err = keyfunc.NewJWKSetJSON([]byte(jwksUrl))
		if err != nil {
			return nil, fmt.Errorf("create keyfunc from json: %w", err)
		}
	} else {
		kf, err = keyfunc.NewDefault([]string{jwksUrl})
		if err != nil {
			return nil, fmt.Errorf("create jwt keyfunc: %w", err)
		}
	}

	return &Validator{
		keyFunc: kf.Keyfunc,
	}, nil
}
