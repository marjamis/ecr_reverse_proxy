package main

import (
	"encoding/base64"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/api"

	ecr "github.com/awslabs/amazon-ecr-credential-helper/ecr-login"

	log "github.com/sirupsen/logrus"
	defaultLog "log"
	_ "net/http/pprof"
)

var (
	CWSVC                 *cloudwatch.CloudWatch
	CWNamespace           = "Custom/ecr_reverse_proxy"
	CWMetricName          = "ECRImagePulls"
	CWDimensionRepository = "Repository"
	CWDimensionImage      = "Image"
	CWUnit                = "Count"
	CWValue               = float64(1)

	ECR_REGISTRY string
)

func init() {
	// Set all Environment Variables to internal variables where it makes sense
	ECR_REGISTRY = os.Getenv("ECR_REGISTRY")

	// Set up some basic things based on the environment variables
	CWSVC = cloudwatch.New(session.New(), aws.NewConfig().WithRegion(os.Getenv("REGION")))

	// Setting default level to debug
	log.SetLevel(log.DebugLevel)
}

func main() {
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	ecrCredHelper := ecr.ECRHelper{ClientFactory: api.DefaultClientFactory{}}

	// Configures the log for the below ErrorLog requirements
	logger := log.New()
	w := logger.Writer()
	defaultLog.SetOutput(logger.Writer())
	defer w.Close()

	log.WithFields(log.Fields{
		"Function": "main",
	}).Fatal("The application has exited. Check you have the required ENV's.\n", http.ListenAndServeTLS(":"+os.Getenv("PORT"),
		os.Getenv("TLS_CERTIFICATE"),
		os.Getenv("TLS_PRIVATE_KEY"),
		&httputil.ReverseProxy{
			Director: func(r *http.Request) {
				username, password, err := ecrCredHelper.Get(ECR_REGISTRY)
				if err != nil {
					log.WithFields(log.Fields{
						"Function": "main/ReverseProxy",
					}).Fatal("Unable to get ECR credentials. Error: " + err.Error() + "\n")
				}
				r.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(username+":"+password)))

				r.URL.Scheme = "https"
				r.URL.Host = ECR_REGISTRY
				r.Host = ECR_REGISTRY

				log.Info("Pulling Image/Manifest")
				log.WithFields(log.Fields{
					"r.Host": r.URL.Host,
					"r.Path": r.URL.Path,
				}).Debug("Request details")
				if strings.Contains(r.URL.Path, "/manifests/") {
					go metrics(r.URL.Path)
				}
			},
			ErrorLog: defaultLog.New(w, "", 0),
			ModifyResponse: func(resp *http.Response) (err error) {
				resp.Header.Set("Location", strings.Replace(resp.Header.Get("Location"), ECR_REGISTRY, os.Getenv("REGISTRY")+":"+os.Getenv("PORT"), 1))
				return
			},
		}))
}

func metrics(path string) {
	split := strings.Split(path, "/")
	repo := ECR_REGISTRY + "/" + split[2]
	image := split[4]
	time := time.Now()
	contextLogger := log.WithFields(log.Fields{
		"Function":                "metrics",
		"CW:Namespace":            CWNamespace,
		"CW:Metric":               CWMetricName,
		"CW:Dimension:Repository": repo,
		"CW:Dimension:Image":      image,
	})

	contextLogger.Debug("+1 PutMetric Data")

	dimensions := make([]*cloudwatch.Dimension, 2)
	dimensions[0] = &cloudwatch.Dimension{
		Name:  &CWDimensionRepository,
		Value: &repo,
	}
	dimensions[1] = &cloudwatch.Dimension{
		Name:  &CWDimensionImage,
		Value: &image,
	}
	metricData := make([]*cloudwatch.MetricDatum, 1)
	metricData[0] = &cloudwatch.MetricDatum{
		Dimensions: dimensions,
		MetricName: &CWMetricName,
		Timestamp:  &time,
		Unit:       &CWUnit,
		Value:      &CWValue,
	}
	_, err := CWSVC.PutMetricData(&cloudwatch.PutMetricDataInput{
		MetricData: metricData,
		Namespace:  &CWNamespace,
	})
	if err != nil {
		contextLogger.Error("+1 PutMetric Data failed with: %s", err)
	}
}
