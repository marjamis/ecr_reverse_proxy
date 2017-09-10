package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	ecr "github.com/awslabs/amazon-ecr-credential-helper/ecr-login"
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/api"

	_ "net/http/pprof"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	ech := ecr.ECRHelper{ClientFactory: api.DefaultClientFactory{}}

	log.Fatal("The application has exited. Check you have the required ENV's.\n", http.ListenAndServeTLS(":"+os.Getenv("PORT"),
		os.Getenv("TLS_CERTIFICATE"),
		os.Getenv("TLS_PRIVATE_KEY"),
		&httputil.ReverseProxy{
			Director: func(r *http.Request) {
				username, password, err := ech.Get(os.Getenv("ECR_REGISTRY"))
				if err != nil {
					log.Fatal("Unable to get ECR credentials. Error: " + err.Error() + "\n")
				}
				r.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(username+":"+password)))

				r.URL.Scheme = "https"
				r.URL.Host = os.Getenv("ECR_REGISTRY")
				r.Host = os.Getenv("ECR_REGISTRY")
			},
			//ErrorLog,
			ModifyResponse: func(resp *http.Response) (err error) {
				resp.Header.Set("Location", strings.Replace(resp.Header.Get("Location"), os.Getenv("ECR_REGISTRY"), os.Getenv("REGISTRY")+":"+os.Getenv("PORT"), 1))
				return
			},
		}))
}
