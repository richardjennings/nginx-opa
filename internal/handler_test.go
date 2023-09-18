package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOpaProxy_Result_Authorised_bool(t *testing.T) {
	o := OpaProxy{
		Config: &Config{
			AuthorisedKey:   "allow",
			AuthorisedValue: "true",
		},
	}
	tc := map[string]interface{}{
		"allow": true,
	}
	r, err := o.Result(tc)
	assert.Nil(t, err)
	assert.Equal(t, Ok, r)
}

func TestOpaProxy_Result_Authorised_string(t *testing.T) {
	o := OpaProxy{
		Config: &Config{
			AuthorisedKey:   "allow",
			AuthorisedValue: "true",
		},
	}
	tc := map[string]interface{}{
		"allow": "true",
	}
	r, err := o.Result(tc)
	assert.Nil(t, err)
	assert.Equal(t, Ok, r)
}

func TestOpaProxy_Result_Unauthorised_bool(t *testing.T) {
	o := OpaProxy{
		Config: &Config{
			AuthorisedKey:   "allow",
			AuthorisedValue: "true",
		},
	}
	tc := map[string]interface{}{
		"allow": false,
	}
	r, err := o.Result(tc)
	assert.Nil(t, err)
	assert.Equal(t, UnAuthorised, r)
}

func TestOpaProxy_Result_Unauthorised_string(t *testing.T) {
	o := OpaProxy{
		Config: &Config{
			AuthorisedKey:   "allow",
			AuthorisedValue: "true",
		},
	}
	tc := map[string]interface{}{
		"allow": "false",
	}
	r, err := o.Result(tc)
	assert.Nil(t, err)
	assert.Equal(t, UnAuthorised, r)
}

func TestOpaProxy_Result_Authenticated_bool(t *testing.T) {
	o := OpaProxy{
		Config: &Config{
			AuthenticatedKey:   "verify",
			AuthenticatedValue: "true",
			AuthorisedKey:      "allow",
			AuthorisedValue:    "true",
		},
	}
	tc := map[string]interface{}{
		"allow":  true,
		"verify": true,
	}
	r, err := o.Result(tc)
	assert.Nil(t, err)
	assert.Equal(t, Ok, r)
}

func TestOpaProxy_Result_Authenticated_string(t *testing.T) {
	o := OpaProxy{
		Config: &Config{
			AuthenticatedKey:   "verify",
			AuthenticatedValue: "true",
			AuthorisedKey:      "allow",
			AuthorisedValue:    "true",
		},
	}
	tc := map[string]interface{}{
		"allow":  "true",
		"verify": "true",
	}
	r, err := o.Result(tc)
	assert.Nil(t, err)
	assert.Equal(t, Ok, r)
}

func TestOpaProxy_Result_Unauthenticated_bool(t *testing.T) {
	o := OpaProxy{
		Config: &Config{
			AuthenticatedKey:   "verify",
			AuthenticatedValue: "true",
		},
	}
	tc := map[string]interface{}{
		"verify": false,
	}
	r, err := o.Result(tc)
	assert.Nil(t, err)
	assert.Equal(t, UnAuthenticated, r)
}

func TestOpaProxy_Result_Unauthenticated_string(t *testing.T) {
	o := OpaProxy{
		Config: &Config{
			AuthenticatedKey:   "verify",
			AuthenticatedValue: "true",
		},
	}
	tc := map[string]interface{}{
		"verify": "false",
	}
	r, err := o.Result(tc)
	assert.Nil(t, err)
	assert.Equal(t, UnAuthenticated, r)
}

func TestOpaProxy_Result_Missing_AuthenticatedKey(t *testing.T) {
	o := OpaProxy{
		Config: &Config{
			AuthenticatedKey:   "verify",
			AuthenticatedValue: "true",
		},
	}
	r, err := o.Result(map[string]interface{}{})
	assert.Nil(t, err)
	assert.Equal(t, UnAuthenticated, r)
}

func TestOpaProxy_Result_Missing_AuthorizedKey(t *testing.T) {
	o := OpaProxy{
		Config: &Config{
			AuthorisedKey:   "verify",
			AuthorisedValue: "true",
		},
	}
	r, err := o.Result(map[string]interface{}{})
	assert.Nil(t, err)
	assert.Equal(t, UnAuthorised, r)
}
