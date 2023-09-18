package system.main

import future.keywords.if

default allow = false

allow if {
    input.attributes.request.http.method == "get"
    input.attributes.request.http.headers.accept[0] == "*/*"
    input.attributes.request.http.headers["some-header"][0] != ""
}
