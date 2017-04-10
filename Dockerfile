FROM jwilder/nginx-proxy:latest
MAINTAINER Brian Palmer <brian@codekitchen.net>

RUN apt-get update \
 && apt-get install -y -q --no-install-recommends \
    dnsmasq \
 && apt-get clean \
 && rm -r /var/lib/apt/lists/*

RUN wget https://github.com/codekitchen/dinghy-http-proxy/releases/download/join-networks-v3/join-networks.tar.gz \
 && tar -C /app -xzvf join-networks.tar.gz \
 && rm join-networks.tar.gz

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
