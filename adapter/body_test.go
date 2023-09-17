package adapter

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
	body, err := RequestBody(r)
	assert.Nil(t, err)
	expected := `{"attributes":{"request":{"http":{"headers":{"x-original-method":["GET"],"x-original-url":["https://test.com/sale?sort=desc"],"x-real-ip":["192.168.0.1"]},"method":"get","scheme":"https","host":"test.com","path":"/sale","query":{"sort":["desc"]}},"ipAddr":"192.168.0.1"}}}`
	assert.Equal(t, expected, body.String())
}
