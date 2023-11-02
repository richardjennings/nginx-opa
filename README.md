# Nginx OPA

# What
Nginx Auth Request OPA policy evaluation.

## How

### Kubernetes Ingress Example

#### Deploy [Kube-Mgmt](https://github.com/open-policy-agent/kube-mgmt) via Helm 
with Nginx OPA Proxy as an extraContainers entry:
```
opa-kube-mgmt:
  extraContainers:
  - name: ino
    image: docker.io/richardjennings/ingress-nginx-opa:0.0.1
    args:
      - serve
      - --tls-cerrt-file=/certs/tls.crt
      - --tls-private-key=/certs/tls.key
    env:
      - name: OPA_URL
        value: https://127.0.0.1:8181
      - name: AUTHENTICATED_KEY
        value: "verified"
      - name: AUTHENTICATED_VALUE
        value: "true"
      - name: AUTHORIZED_KEY
        value: "allow"
      - name: AUTHORIZED_VALUE
        value: "true"
    ports:
      - name: http
        containerPort: 8282
        protocol: TCP
    volumeMounts:
      - mountPath: /certs
        name: certs
        readOnly: true
```

#### Configure Ingress Nginx to send auth requests to the Nginx OPA Proxy, e.g. using helm chart config:
```
ingress-nginx:
  controller:
    config:
      global-auth-url: https://ino.opa.svc.cluster.local:8282
```

Nginx Auth Requests are transformed into a json input which is sent to the defined OPA REST API address
e.g. POST `http://opa-svc.opa.svc.cluster.local:8181`

The inputs provided to OPA by Nginx OPA Proxy from the Ingress Nginx Auth Request are of the form:

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

#### Deploy A Policy for OPA to evaluate against Nginx Auth Requests. 

For example a Policy to validate a Google Identity Aware Proxy JWT token signature and apply Authorization based on 
the email in the JWT token and Host header might look like:

```
kind: ConfigMap
apiVersion: v1
metadata:
  namespace: opa
  name: opa-example
  labels:
    openpolicyagent.org/policy: rego
data:
  main: |
    package system.main
    import future.keywords
    default allow = false
    default verified = false
    
    jwks_request(url) = http.send({
      "url": url,
      "method": "GET",
      "force_cache": true,
      "force_cache_duration_seconds": 60
    })
    
    jwks = jwks_request("https://www.gstatic.com/iap/verify/public_key-jwk").body
    verified = io.jwt.verify_es256(input.attributes.request.headers["x-goog-iap-jwt-assertion"][0], json.marshal(jwks))
    
    allow if {
      input.attributes.http.host == "example.domain.tld"
      verified == true
      {
        "a@domain.tld",
        "b@domain.tld",
      }[jwt.payload.gcpip.email]
    }
```


The result of policy evaluation is compared to a predefined expectation and either a `200`, `401` or `403` is
returned to Nginx. AuthenticatedKey and AuthenticatedValue may be set to match a key value in the OPA response.
If set and a match is not found, a 401 response is returned. AuthorisedKey and AuthorisedValue default to `allow` and
`true` respectively. If Authorised Key and Value do not match a key value in the OPA response a 403 is returned.

## Status

Beta