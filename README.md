# OPA Nginx

# What
OPA integration with [Nginx Auth Requests](https://nginx.org/en/docs/http/ngx_http_auth_request_module.html).

## How

Nginx Auth Requests are transformed into a JSON structure. This structure is then sent to the defined OPA REST API 
address e.g. POST `http://opa-nginx.opa-nginx.svc.cluster.local:8282` to be evaluated.

The inputs provided to OPA by Nginx OPA Proxy are of the form:

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

A policy written to allow requests with a host header of `test.com` may look like:

```
      package system.main
      import future.keywords
      default allow := false
      
      allow if {
        input.attributes.request.http.host == "test.com"
      }
```

As a result of policy evaluation against the inputs, either a `200`, `401` or `403` status code is returned to Nginx. 

A `200` response is returned if the OPA result includes key value pairs matching the configured values for `AuthorisedKey`,
`AuthorisedValue` and optionally `AuthenticatedKey` and `AuthenticatedValue`.

If `AuthenticatedKey` is configured a  `401` response is returned if `AuthenticatedKey` is not present in the OPA result 
or does not match `AuthenticatedValue`.

A `403` response is returned if `AuthorisedKey` is not set in the OPA result or does not match `AuthorisedValue`.

AuthorisedKey and AuthorisedValue default to `allow` and `true` respectively.

## Kubernetes

### Helm Chart

### Configure Ingress Nginx to send auth requests to the Nginx OPA Proxy, e.g. using helm chart config:
```
ingress-nginx:
  controller:
    config:
      global-auth-url: https://opa-nginx.opa-nginx.svc.cluster.local:8282
```


## Status

Beta