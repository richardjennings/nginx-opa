package integration

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func Test_Integration(t *testing.T) {
	// @todo currently should be run from parent directory
	ctx := context.Background()
	nginxConfFile, err := filepath.Abs("integration/nginx.conf")
	if err != nil {
		t.Fatal(err)
	}
	net, err := network.New(ctx)
	if err != nil {
		t.Fatal(err)
	}
	testPolicy := `package system.main

import future.keywords.if

default allow = false
default verified = false

allow if {
    input.attributes.request.http.headers["test"][0] == "ok"
    input.attributes.request.http.method == "get"
}

verified if {
	input.attributes.request.http.headers["jwt"][0] == "valid"
}
`

	requests := testcontainers.ParallelContainerRequest{
		{
			ContainerRequest: testcontainers.ContainerRequest{
				FromDockerfile: testcontainers.FromDockerfile{
					Context: ".",
				},
				Env: map[string]string{
					"OPA_URL":             "http://opa:8181",
					"AUTHENTICATED_KEY":   "verified",
					"AUTHENTICATED_VALUE": "true",
					"AUTHORIZED_KEY":      "allow",
					"AUTHORIZED_VALUE":    "true",
				},
				Hostname:     "opa-nginx",
				ExposedPorts: []string{"8282:8282/tcp"},
				WaitingFor:   wait.ForListeningPort("8282/tcp"),
				Networks:     []string{net.Name},
			},
			ProviderType: testcontainers.ProviderDocker,
			Started:      true,
		},
		{
			ContainerRequest: testcontainers.ContainerRequest{
				Image: "nginx",

				ExposedPorts: []string{"8080:80/tcp"},
				WaitingFor:   wait.ForHTTP("/healthz"),
				Files: []testcontainers.ContainerFile{
					{
						//Reader:            r,
						HostFilePath:      nginxConfFile, // will be discarded internally
						ContainerFilePath: "/etc/nginx/nginx.conf",
						FileMode:          0o600,
					},
				},
				//Mounts: testcontainers.ContainerMounts{
				//	testcontainers.ContainerMount{
				//		Source: testcontainers.GenericVolumeMountSource{}
				//		//Source: testcontainers.GenericBindMountSource{
				//		//	HostPath: nginxConfFile,
				//		//},
				//		Target: "/etc/nginx/nginx.conf",
				//	},
				//},
				Hostname: "nginx",
				Networks: []string{net.Name},
			},
			ProviderType: testcontainers.ProviderDocker,
			Started:      true,
		},
		{
			ContainerRequest: testcontainers.ContainerRequest{
				Image:        "openpolicyagent/opa:0.57.0-dev",
				ExposedPorts: []string{"8181:8181/tcp"},
				WaitingFor:   wait.ForListeningPort("8181/tcp"),
				Cmd:          []string{"run", "--server", "--set", "decision_logs.console=true", "--addr", "0.0.0.0:8181"},
				Hostname:     "opa",
				Networks:     []string{net.Name},
			},
			ProviderType: testcontainers.ProviderDocker,
			Started:      true,
		},
	}
	containers, err := testcontainers.ParallelContainers(ctx, requests, testcontainers.ParallelContainersOptions{})
	if err != nil {
		t.Fatal(err)
	}
	log.Println("container started")

	// Expect a 401 when no Policy is in place
	res, err := http.DefaultClient.Get("http://localhost:8080/ok")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	log.Println("load policy")

	// load rego policy
	r, err := http.NewRequest(
		"PUT",
		"http://localhost:8181/v1/policies/data.system.main",
		bytes.NewBufferString(testPolicy),
	)
	if err != nil {
		t.Fatal(err)
	}
	res, err = http.DefaultClient.Do(r)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// Policy update pause
	time.Sleep(time.Second * 5)

	log.Println("test policy passes")

	// check for a 200 with verb get headre test[0] ok
	r, _ = http.NewRequest("GET", "http://localhost:8080", nil)
	r.Header.Add("test", "ok")
	r.Header.Add("jwt", "valid")
	res, err = http.DefaultClient.Do(r)
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	// a 404 is considered a 200 without the content configured to be served by nginx
	assert.Equal(t, http.StatusNotFound, res.StatusCode)

	log.Println("test policy fail")

	// check for a Forbidden 403
	r, _ = http.NewRequest("GET", "http://localhost:8080", nil)
	r.Header.Add("test", "notok")
	r.Header.Add("jwt", "valid")
	res, err = http.DefaultClient.Do(r)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusForbidden, res.StatusCode)

	// check for unauthorized 401
	r, _ = http.NewRequest("GET", "http://localhost:8080", nil)
	r.Header.Add("test", "ok")
	r.Header.Add("jwt", "notvalid")
	res, err = http.DefaultClient.Do(r)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	log.Println("cleanup")

	// defer cleanup
	for _, c := range containers {
		c := c
		defer func() {
			if err := c.Terminate(ctx); err != nil {
				t.Fatalf("failed to terminate container: %s", c)
			}
		}()
	}

	t.Cleanup(func() {
		require.NoError(t, net.Remove(ctx))
	})

}
