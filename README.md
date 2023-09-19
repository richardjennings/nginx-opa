# Ingress Nginx OPA

# What
Integration between Open Policy Agent and Ingress Nginx to allow OPA policy evaluation via Nginx Auth Request module.

## How

Ingress Nginx is configured to send Auth Requests to this proxy, e.g. via a `global-auth-url` entry in the Ingress Nginx Configmap. 

This request is transformed into a policy evaluation request sent to an OPA REST API address, 
e.g. POST `http://opa-svc.opa.svc.cluster.local:8181`

The result of policy evaluation is compared to predefined expectation and either a `200`, `401` or `403` is
returned to Nginx. AuthenticatedKey and AuthenticatedValue may be set to match a key value in the OPA response. 
If set and a match is not found, a 401 response is returned. AuthorisedKey and AuthorisedValue default to `allow` and 
`true` respectively. If Authorised Key and Value do not match a key value in the OPA response a 403 is returned.

The inputs provided to OPA are of the form:

```
{
  "attributes": {
    "request": {
      "http": {
        "headers": {
          "example": [
            "value"
          ]
        },
        "path": "/some/path",
        "query": {
          "foo": "bar"
        },
        "method": "GET",
        "host": "test.com"
      },
      "ipAddr": "10.10.10.10"
    }
  }
}
```

An example policy for the above might look like:

```
package system.main

import future.keywords.if

default allow = false

allow if {
    input.attributes.request.http.headers.example[0] == "value"
    input.attributes.request.http.method == "GET"
    input.attributes.request.http.host == "test.com"
    input.attributes.request.http.query.foo[0] == "bar"
}
```

## Status

Alpha