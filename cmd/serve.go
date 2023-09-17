package cmd

import (
	"bytes"
	"crypto/tls"
	"github.com/richardjennings/opa-nginx-auth-request/adapter"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const defaultExpected = `{"allow":true}`

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "run OPA Nginx Auth Request proxy service",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		upstream, ok := os.LookupEnv("OPA_URL")
		if !ok {
			log.Fatalf("OPA_URL environment variable required")
		}
		expected, ok := os.LookupEnv("EXPECTED_RESPONSE")
		if !ok {
			expected = defaultExpected
		}
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		server := &http.Server{
			Addr: ":8282",
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, err := adapter.RequestBody(r)
				if err != nil {

				}
				res, err := http.Post(upstream, "application/json", body)

				if err != nil {
					w.WriteHeader(500)
					_, err := w.Write([]byte(err.Error()))
					if err != nil {
						log.Printf("could not write error to client %s", err)
						return
					}
					return
				}

				buff := bytes.NewBuffer(nil)
				_, err = io.Copy(buff, res.Body)
				defer func() { _ = res.Body.Close() }()
				if err != nil {
					log.Printf("could not copy response body: %s", err)
				}

				if res.StatusCode != 200 {
					w.WriteHeader(401)
					return
				}

				if strings.TrimSuffix(buff.String(), "\n") != expected {
					w.WriteHeader(401)
					return
				}

				w.WriteHeader(200)
			}),
		}
		log.Fatal(server.ListenAndServe())
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
