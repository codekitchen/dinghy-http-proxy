FROM jwilder/nginx-proxy:alpine
MAINTAINER Brian Palmer <brian@codekitchen.net>

### Install Application
RUN apk upgrade --no-cache && \
    apk add --no-cache --virtual=run-deps \
      su-exec \
      curl \
      dnsmasq && \
    rm -rf /tmp/* \
           /var/cache/apk/*  \
           /var/tmp/*

COPY join-networks /app/
RUN chmod +x /app/join-networks

COPY Procfile /app/

# override nginx configs
COPY *.conf /etc/nginx/conf.d/

# override nginx-proxy templating
COPY nginx.tmpl Procfile reload-nginx /app/

COPY htdocs /var/www/default/htdocs/

ENV DOMAIN_TLD docker
ENV DNS_IP 127.0.0.1
ENV HOSTMACHINE_IP 127.0.0.1

HEALTHCHECK CMD curl --fail http://localhost/ || exit 1

EXPOSE 19322
