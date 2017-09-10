default:

build:
	docker build -t ecr_reverse_proxy .

buildrun:
ifndef REGISTRY
	echo "REGISTRY not set"
	exit 1
endif
ifndef ECR_REGISTRY
	echo "ECR_REGISTRY not set"
	exit 1
endif
	docker run --rm -it -v $$PWD/credentials/:/credentials/:ro -e TLS_CERTIFICATE="/credentials/cert.pem" -e TLS_PRIVATE_KEY="/credentials/privkey.pem" -e REGISTRY=$(REGISTRY) -e ECR_REGISTRY=$(ECR_REGISTRY) -e PORT=5000 -v /Users/marjamis/.aws:/.aws:ro -p 6060:6060 -p 5000:5000 ecr_reverse_proxy

lbuild:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/ecr_reverse_proxy .

run:
ifndef REGISTRY
	echo "REGISTRY not set"
	exit 1
endif
ifndef ECR_REGISTRY
	echo "ECR_REGISTRY not set"
	exit 1
endif
	TLS_CERTIFICATE="./credentials/cert.pem" TLS_PRIVATE_KEY="./credentials/privkey.pem" REGISTRY=$(REGISTRY) ECR_REGISTRY=$(ECR_REGISTRY) PORT=5000 go run main.go

test:
