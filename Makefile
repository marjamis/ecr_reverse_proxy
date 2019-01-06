export PORT=5000
export TLS_CERTIFICATE="/credentials/cert.pem"
export TLS_PRIVATE_KEY="/credentials/privkey.pem"

default:
	echo "Please run a specific target as this first one will do nothing."

checks:
ifndef REGION
	echo "REGION not set"
	exit 1
endif
ifndef REGISTRY
	echo "REGISTRY not set"
	exit 1
endif
ifndef ECR_REGISTRY
	echo "ECR_REGISTRY not set"
	exit 1
endif

playtest: checks genSelfsigned dockerBuild dockerRun

genSelfsigned:
	# Generates self-sign certificates
	openssl req -x509 -newkey rsa:4096 -keyout ./credentials/privkey.pem -out ./credentials/cert.pem -days 365 -nodes

dockerBuild:
	docker build -t ecr_reverse_proxy .

dockerRun: checks
	docker run --rm -it -v $$PWD/credentials/:/credentials/:ro -e TLS_CERTIFICATE=$(TLS_CERTIFICATE) -e REGION=$(REGION) -e TLS_PRIVATE_KEY=$(TLS_PRIVATE_KEY) -e REGISTRY=$(REGISTRY) -e ECR_REGISTRY=$(ECR_REGISTRY) -e PORT=$(PORT) -v $$HOME/.aws:/.aws:ro -p 6060:6060 -p $(PORT):$(PORT) ecr_reverse_proxy

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/ecr_reverse_proxy .

run: checks
	TLS_CERTIFICATE=$(TLS_CERTIFICATE) TLS_PRIVATE_KEY=$(TLS_PRIVATE_KEY) REGION=$(REGION) REGISTRY=$(REGISTRY) ECR_REGISTRY=$(ECR_REGISTRY) PORT=$(PORT) go run main.go

# Replaces variables in the generic k8s manifest for internal testing
k8s-manifestReplace:
	cat manifest.k8s.yml | sed -f .env > used-manifest.k8s.yml

# Generates and adds the TLS secrets for the Ingress Controller to use
k8s-addSecrets: genSelfsigned
	kubectl create secret tls ecrreverseproxy --cert ./credentials/cert.pem --key ./credentials/privkey.pem
