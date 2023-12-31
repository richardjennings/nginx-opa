package internal

import (
	"net/http"
	"net/url"
	"strings"
)

type (
	Input struct {
		Attributes Attributes `json:"attributes"`
	}
	Attributes struct {
		Request Request `json:"request"`
	}
	Request struct {
		Http   Http   `json:"http"`
		IpAddr string `json:"ipAddr"`
	}
	Http struct {
		Headers map[string][]string `json:"headers"`
		Method  string              `json:"method"`
		Scheme  string              `json:"scheme"`
		Host    string              `json:"host"`
		Path    string              `json:"path"`
		Query   map[string][]string `json:"query"`
	}
)

func Inputs(r *http.Request) *Input {
	inputs := &Input{
		Attributes: Attributes{
			Request: Request{
				Http: Http{
					Headers: make(http.Header),
				},
			},
		},
	}

	// lower case header keys
	for k, v := range r.Header {
		inputs.Attributes.Request.Http.Headers[strings.ToLower(k)] = v
	}

	// X-Original-Method
	if v, ok := inputs.Attributes.Request.Http.Headers["x-original-method"]; ok {
		inputs.Attributes.Request.Http.Method = strings.ToLower(v[0])
	}

	// X-Real-IP
	if v, ok := inputs.Attributes.Request.Http.Headers["x-real-ip"]; ok {
		inputs.Attributes.Request.IpAddr = v[0]
	}

	// X-Original-URL
	if v, ok := inputs.Attributes.Request.Http.Headers["x-original-url"]; ok {
		originalUrl, err := url.Parse(v[0])
		if err == nil {
			inputs.Attributes.Request.Http.Scheme = originalUrl.Scheme
			inputs.Attributes.Request.Http.Host = originalUrl.Host
			inputs.Attributes.Request.Http.Path = originalUrl.Path
			inputs.Attributes.Request.Http.Query = originalUrl.Query()
		}
	}
	return inputs
}
