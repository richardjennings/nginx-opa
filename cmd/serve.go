package cmd

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type (
	Input struct {
		Attributes Attributes `json:"attributes"`
	}
	Attributes struct {
		Request Request `json:"request"`
	}
	Request struct {
		Http Http `json:"http"`
	}
	Http struct {
		Path    string              `json:"path"`
		Method  string              `json:"method"`
		Headers map[string][]string `json:"headers"`
	}
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
		expected, ok := os.LookupEnv("EXPECTED_RESPONSE")
		if !ok {
			expected = `{"allow":true}`
		}
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		server := &http.Server{
			Addr: ":8282",
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				inputs := Input{
					Attributes: Attributes{
						Request: Request{
							Http: Http{
								Headers: make(http.Header),
							},
						},
					},
				}

				for k, v := range r.Header {
					inputs.Attributes.Request.Http.Headers[strings.ToLower(k)] = v
				}

				inputJson, err := json.Marshal(inputs)
				if err != nil {
					log.Printf("error encoding input json: %s", err)
					return
				}

				res, err := http.Post(upstream, "application/json", bytes.NewBuffer(inputJson))
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
					log.Println("response %s", buff.String())
					return
				}

				log.Println(buff.String())
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
