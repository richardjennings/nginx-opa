# OPA Nginx Auth Request 

Provide a Nginx Auth Request compatible response from an OPA REST Policy evaluation response.

## How

Nginx Auth Request sends request to OPA Nginx Auth Request.

The requests header keys are lower-cased and posted to OPA as json to be evaluated as inputs.

The inputs provided to OPA are of the form:

```
"attributes": {
    "request": {
        "http": {
            "headers": {
                "example": ["value"]
            },
            //"path": "/some/path",
            //"query": {
            //    "foo": "bar"
            //},
            //"method": "GET",
            //"host": "test.com"
        }
    }
}
```

A policy can be evaluated against these inputs. 

The result of policy evaluation is compared to predefined expectation and either a `200` or `401` is 
returned to Nginx.

An example policy for the above might look like:

```
package system.main

import future.keywords.if

default allow = false

allow if {
    input.attributes.request.http.headers.example[0] == "value"
    #input.attributes.request.http.method == "GET"
    #input.attributes.request.http.host == "test.com"
    #input.attributes.request.http.query.foo == "bar"
}
```

By default, multiple aspects of the original request are not passed by Nginx in the Auth Request.

The following headers can be specified in Nginx config and are translated to OPA inputs:
```
proxy_set_header            X-Original-URL          $scheme://$http_host$request_uri;
proxy_set_header            X-Original-Method       $request_method;
proxy_set_header            X-Real-IP               $remote_addr;
proxy_set_header            X-Forwarded-For         $full_x_forwarded_for;
```


## Status

Proof of concept.