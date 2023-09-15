
.PHONY = proxy opa nginx rego

proxy:
	OPA_URL=http://localhost:8181 go run main.go serve

opa:
	docker run --platform=linux/amd64 --rm -p8181:8181 openpolicyagent/opa:0.57.0-dev run --server --addr 0.0.0.0:8181

nginx:
	docker run --rm -p8080:80 -v $$PWD/test/nginx.conf:/etc/nginx/nginx.conf --workdir /app nginx

rego:
	curl -v -T test/main.rego http://localhost:8181/v1/policies/data.system.main
	curl -v http://localhost:8181/v1/policies
