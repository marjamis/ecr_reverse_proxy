# ecs_reverse_proxy

The proxy will need a few environment variables to configure how it runs. These currently being:
* TLS_CERTIFICATE # Location to the certificate used for secure communication to the proxy.
* TLS_PRIVATE_KEY # The private key for the cert.
* ECR_REGISTRY # The ECR Registry you're proxying to.
* REGISTRY # The CNAME of the registry the proxy will accept connections on.
* PORT # The port in which the proxy will listen on.

It also requires AWS credentials, which can be passed into the application as per all the normal SDK locations though preferably instance profiles. The [ECR Credential Helper](https://github.com/awslabs/amazon-ecr-credential-helper) uses these credentials to automatically track and obtain the required credentials for ECR.
