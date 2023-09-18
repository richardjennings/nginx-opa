package internal

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestBody_default(t *testing.T) {
	r := &http.Request{
		Header: map[string][]string{
			"X-Original-Method": {"GET"},
			"X-Original-URL":    {"https://test.com/sale?sort=desc"},
			"X-Real-IP":         {"192.168.0.1"},
		},
	}
	input := Inputs(r)
	expected := &Input{
		Attributes: Attributes{
			Request: Request{
				Http: Http{
					Headers: map[string][]string{
						"x-original-method": {"GET"},
						"x-original-url":    {"https://test.com/sale?sort=desc"},
						"x-real-ip":         {"192.168.0.1"},
					},
					Method: "get",
					Scheme: "https",
					Host:   "test.com",
					Path:   "/sale",
					Query: map[string][]string{
						"sort": {"desc"},
					},
				},
				IpAddr: "192.168.0.1",
			},
		},
	}
	assert.Equal(t, expected, input)
}
