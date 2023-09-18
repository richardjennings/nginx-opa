package cmd

import (
	"crypto/tls"
	"github.com/richardjennings/opa-nginx-auth-request/internal"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

var rootCmd = &cobra.Command{}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "run Ingress Nginx OPA Auth Proxy",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		config, err := internal.NewConfig()
		if err != nil {
			log.Fatalln(err)
		}
		server := &http.Server{
			Addr:    ":8282",
			Handler: internal.NewHandler(&internal.OpaProxy{Config: config}),
		}
		log.Fatal(server.ListenAndServe())
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
