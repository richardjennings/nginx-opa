package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	Invalid Result = iota
	Ok
	UnAuthorised
	UnAuthenticated
)

const defaultAuthorisedKey = "allow"
const defaultAuthorisedValue = "true"

type (
	Result int
	Config struct {
		Upstream           string
		AuthorisedKey      string
		AuthorisedValue    string
		AuthenticatedKey   string
		AuthenticatedValue string
	}
	Opa interface {
		Handle(input *Input) (Result, error)
	}
	OpaProxy struct {
		Config *Config
	}
)

func NewConfig() (*Config, error) {
	upstream, ok := os.LookupEnv("OPA_URL")
	if !ok {
		return nil, errors.New("OPA_URL environment variable required")
	}
	authenticatedKey := os.Getenv("AUTHENTICATED_KEY")
	authenticatedValue := os.Getenv("AUTHENTICATED_VALUE")
	authorisedKey, ok := os.LookupEnv("AUTHORIZED_KEY")
	if !ok {
		authorisedKey = defaultAuthorisedKey
	}
	authorisedValue, ok := os.LookupEnv("AUTHORIZED_VALUE")
	if !ok {
		authorisedValue = defaultAuthorisedValue
	}

	return &Config{
		Upstream:           upstream,
		AuthorisedKey:      authorisedKey,
		AuthorisedValue:    authorisedValue,
		AuthenticatedKey:   authenticatedKey,
		AuthenticatedValue: authenticatedValue,
	}, nil
}

func (o *OpaProxy) Handle(input *Input) (Result, error) {
	b, err := json.Marshal(input)
	if err != nil {
		return Invalid, err
	}
	body := bytes.NewBuffer(b)
	res, err := http.Post(o.Config.Upstream, "application/json", body)
	if err != nil {
		return Invalid, err
	}
	buff := bytes.NewBuffer(nil)
	_, err = io.Copy(buff, res.Body)
	defer func() { _ = res.Body.Close() }()
	if err != nil {
		return Invalid, fmt.Errorf("could not copy response body: %s", err)
	}
	var v map[string]interface{}
	if err := json.Unmarshal(buff.Bytes(), &v); err != nil {
		return Invalid, fmt.Errorf("could not unmarshal json, %s", v)
	}
	return o.Result(v)
}

func (o *OpaProxy) Result(v map[string]interface{}) (Result, error) {
	// if Authenticated Key and Value are set, handle first.
	if o.Config.AuthenticatedKey != "" && o.Config.AuthenticatedValue != "" {
		authN, ok := v[o.Config.AuthenticatedKey]
		if !ok {
			return UnAuthenticated, nil
		}
		switch authN := authN.(type) {
		case bool:
			if authN && (o.Config.AuthenticatedValue != "true") {
				return UnAuthenticated, nil
			} else if !authN && (o.Config.AuthenticatedValue != "false") {
				return UnAuthenticated, nil
			}
			return Ok, nil
		case string:
			if authN != o.Config.AuthenticatedValue {
				return UnAuthenticated, nil
			}
			return Ok, nil
		default:
			return Invalid, errors.New("unexpected type found in OPA response for AuthenticatedKey")
		}
	}
	// Authorisation
	authZ, ok := v[o.Config.AuthorisedKey]
	if !ok {
		return UnAuthorised, nil
	}
	switch authZ := authZ.(type) {
	case bool:
		if authZ && (o.Config.AuthorisedValue != "true") {
			return UnAuthorised, nil
		} else if !authZ && (o.Config.AuthorisedValue != "false") {
			return UnAuthorised, nil
		}
		return Ok, nil
	case string:
		if authZ != o.Config.AuthorisedValue {
			return UnAuthorised, nil
		}
		return Ok, nil
	default:
		return Invalid, errors.New("unexpected type found in OPA response for AuthenticatedKey")
	}
}

func NewHandler(opa Opa) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		inputs := Inputs(r)
		result, err := opa.Handle(inputs)
		if err != nil {
			log.Printf("error: %s", err)
			w.WriteHeader(500)
			return
		}
		switch result {
		default:
			w.WriteHeader(500)
			return
		case Ok:
			w.WriteHeader(200)
			return
		case UnAuthenticated:
			w.WriteHeader(401)
			return
		case UnAuthorised:
			w.WriteHeader(403)
			return
		}
	}
}
