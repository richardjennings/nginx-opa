package cmd

import (
	"crypto/tls"
	"github.com/richardjennings/opa-nginx/internal"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"time"
)

var rootCmd = &cobra.Command{}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

const defaultAddr = ":8282"

var tlsCertFile string
var tlsPrivateKeyFile string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "run Ingress Nginx OPA Auth Proxy",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		http.DefaultTransport.(*http.Transport).ResponseHeaderTimeout = 5 * time.Second
		http.DefaultClient.Timeout = 10 * time.Second
		config, err := internal.NewConfig()
		if err != nil {
			log.Fatalln(err)
		}

		mux := http.NewServeMux()
		mux.HandleFunc("/", internal.NewHandler(&internal.OpaProxy{Config: config}))
		mux.HandleFunc("/healthz", func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(200)
		})
		server := &http.Server{
			Addr:         defaultAddr,
			Handler:      mux,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		}

		if tlsCertFile != "" && tlsPrivateKeyFile != "" {
			log.Fatalln(server.ListenAndServeTLS(tlsCertFile, tlsPrivateKeyFile))
		}
		log.Fatal(server.ListenAndServe())
	},
}

func init() {
	serveCmd.Flags().StringVar(&tlsCertFile, "tls-cert-file", "", "set path of TLS certificate file")
	serveCmd.Flags().StringVar(&tlsPrivateKeyFile, "tls-private-key-file", "", "set path of TLS private key file")

	rootCmd.AddCommand(serveCmd)
}
