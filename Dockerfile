#get our build tool image. Thanks lush!
FROM lushdigital/docker-golang-dep as build

#load up source
WORKDIR /go/src/clamor

#get deps and cache this layer
COPY src/Gopkg.toml /go/src/clamor/Gopkg.toml
COPY src/Gopkg.lock /go/src/clamor/Gopkg.lock
RUN dep ensure --vendor-only

#copy the source as long as the deps don't change we'll skip fetching them
COPY src /go/src/clamor

#disable crosscompiling 
ENV CGO_ENABLED=0
#compile linux only
ENV GOOS=linux
#build
RUN go build -ldflags '-w -s' -a -installsuffix cgo -o gateway

#take output binary and build off a scratch image
FROM drone/ca-certs
COPY --from=build /go/src/clamor/gateway /gateway
CMD ["/gateway", "serve"]
