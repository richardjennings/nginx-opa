package cmd

import (
	"bytes"
	"crypto/tls"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net/http"
	url2 "net/url"
	"os"
	"strings"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "run OPA Nginx Auth Request proxy service",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		upstream, ok := os.LookupEnv("OPA_URL")
		if !ok {
			log.Fatalf("OPA_URL environment variable required")
		}
		url, err := url2.Parse(upstream)
		if err != nil {
			log.Fatalf("suppled URL %s invalid: %s", args[0], err)
		}
		expected, ok := os.LookupEnv("EXPECTED_RESPONSE")
		if !ok {
			expected = `{"allow":true}`
		}
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		server := &http.Server{
			Addr: ":8282",
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				req := &http.Request{}
				req.URL = url
				req.Header = r.Header
				req.Method = http.MethodPost
				req.Body = nil
				res, err := http.DefaultClient.Do(req)
				if err != nil {
					w.WriteHeader(500)
					_, err := w.Write([]byte(err.Error()))
					if err != nil {
						log.Printf("could not write error to client %s", err)
						return
					}
				}
				buff := bytes.NewBuffer(nil)
				defer func() { _ = res.Body.Close() }()
				_, err = io.Copy(buff, res.Body)
				if err != nil {
					log.Printf("could not copy response body: %s", err)
				}
				if strings.TrimSuffix(buff.String(), "\n") != expected {
					w.WriteHeader(401)
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
