FROM lushdigital/docker-golang-dep as dep
COPY src /go/src/clamor
WORKDIR /go/src/clamor
RUN dep ensure

FROM golang:latest as build
COPY --from=dep /go/src/clamor /go/src/clamor
WORKDIR /go/src/clamor
RUN go build  -o gateway
CMD ["/go/src/clamor/gateway", "serve"]
