FROM golang:1.9 AS build
WORKDIR /go/src/ecr_reverse_proxy/
RUN apt-get update && apt-get install unzip
RUN cd /tmp && wget -L https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && unzip dep-linux-amd64
COPY Gopkg.* /go/src/ecr_reverse_proxy/
COPY *.go /go/src/ecr_reverse_proxy/
RUN /tmp/dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ecr_reverse_proxy .

FROM alpine:latest
LABEL maintainer=marjamis
RUN apk --no-cache add ca-certificates && mkdir /.ecr/ && chown nobody:nobody /.ecr/
USER nobody
WORKDIR /app/
COPY --from=build /go/src/ecr_reverse_proxy/ecr_reverse_proxy .
ENTRYPOINT ["./ecr_reverse_proxy"]
