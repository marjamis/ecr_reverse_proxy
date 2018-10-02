# ecr_reverse_proxy
## Configuration
### Environment Variables
The proxy will need a few environment variables to configure how it runs. These are all also configurable via the Makefile environment variables, such as:
```bash
make -e NAME=value <target>
```
A full list of available environment variables:

| Name | Value |
| --- | ------ |
| TLS_CERTIFICATE | Location to the certificate used for secure communication to the proxy |
| TLS_PRIVATE_KEY |  Location to the private key for the certificate |
| ECR_REGISTRY | The ECR Registry you're proxying to. |
| REGISTRY | The CNAME of the proxy which will be used by clients|
| PORT | The port in which the proxy will listen on |
| REGION | Region for the CloudWatch metrics to be pushed to |

### Additional
It also requires AWS credentials, which can be passed into the application as per all the normal SDK locations though preferably instance profiles.

The [ECR Credential Helper](https://github.com/awslabs/amazon-ecr-credential-helper) uses these credentials to automatically track and obtain the required credentials for ECR. Meaning there is no need to worry about credential rotation for the Tokens that the docker client sends to ECR.

## Running
### Basic Testing
#### Proxy Setup
A basic example to test locally will run on port 5000 and use  generated self-signed certificates. This can be accomplished with:
```bash
make -e REGION=<CWMetricRegion> -e ECR_REGISTRY=<BackendECR> -e REGISTRY=<ProxyURL> playtest
```

#### Testing with docker
Once the application is run you can point your Docker commands to the proxy URL which will in turn use the configured ECR Registry to push/pull the specified image. I.e.:
```bash
docker push <ProxyURL>/<Repository>:<Tag>
docker pull <ProxyURL>/<Repository>:<Tag>
```

### Normal Usage
The way to run this outside of testing would be to build the image with:
```bash
make dockerBuild
```
and once the image has been created it can be run via your normal orchestration system with the required environment variables, as listed above, to configure it's operation.

## Profiling
The package **net/http/pprof** is also in use and opens the profiling on:
```
http://localhost:6060/debug/pprof/
```
