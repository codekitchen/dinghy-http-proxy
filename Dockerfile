FROM golang:1.16 as builder
WORKDIR /go/src/github.com/codekitchen/dinghy-http-proxy
COPY join-networks.go .
COPY go.mod .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go get -v github.com/fsouza/go-dockerclient
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o join-networks

FROM jwilder/nginx-proxy:alpine
LABEL Author="Brian Palmer <brian@codekitchen.net>"

RUN apk upgrade --no-cache \
     && apk add --no-cache --virtual=run-deps \
     su-exec \
     curl \
     dnsmasq \
     && rm -rf /tmp/* \
     /var/cache/apk/* \
     /var/tmp/*

COPY --from=builder /go/src/github.com/codekitchen/dinghy-http-proxy/join-networks /app/join-networks

COPY Procfile /app/

# override nginx configs
COPY *.conf /etc/nginx/conf.d/

# override nginx-proxy templating
COPY nginx.tmpl Procfile reload-nginx /app/

COPY htdocs /var/www/default/htdocs/

ENV DOMAIN_TLD docker
ENV DNS_IP 127.0.0.1
ENV HOSTMACHINE_IP 127.0.0.1

EXPOSE 19322
