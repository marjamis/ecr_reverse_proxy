# ecr_reverse_proxy
## Configuration
### Environment Variables
The proxy will need a few environment variables to configure how it runs. These currently being:
* **TLS_CERTIFICATE** # Location to the certificate used for secure communication to the proxy.
* **TLS_PRIVATE_KEY** #  Location to the private key for the cert.
* **ECR_REGISTRY** # The ECR Registry you're proxying to.
* **REGISTRY** # The CNAME of the registry the proxy will accept connections on.
* **PORT** # The port in which the proxy will listen on.
* **REGION** # Region for the CloudWatch metrics to be pushed to.

### Additional
It also requires AWS credentials, which can be passed into the application as per all the normal SDK locations though preferably instance profiles. The [ECR Credential Helper](https://github.com/awslabs/amazon-ecr-credential-helper) uses these credentials to automatically track and obtain the required credentials for ECR.

## Running
A basic example of this via the Makefile is:
```bash
make -e REGION=<CWMetricRegion> -e ECR_REGISTRY=<BackendECR> -e REGISTRY=<ProxyURL> -e PORT=<Port> run
```
Have a read of the Makefile for some additional configurations, such as building/using this application with Docker.

### Usage
Once the application is run you can point your Docker commands to the proxy URL which will in turn use the configured ECR Registry to push/pull the specified image. I.e.:
```bash
docker push <ProxyURL & Port>/<Repository>:<Image>
docker pull <ProxyURL & Port>/<Repository>:<Image>
```

## Profiling
The package **net/http/pprof** is also in use and opens the profiling to:
```
http://localhost:6060/debug/pprof/
```
