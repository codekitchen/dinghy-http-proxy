FROM jwilder/nginx-proxy:latest
MAINTAINER Brian Palmer <brian@codekitchen.net>

RUN wget https://github.com/codekitchen/dinghy-http-proxy/releases/download/join-networks-v1/join-networks.tar.gz \
 && tar -C /app -xzvf join-networks.tar.gz \
 && rm join-networks.tar.gz

# override nginx configs
COPY *.conf /etc/nginx/conf.d/

# override nginx-proxy templating
COPY nginx.tmpl Procfile reload-nginx /app/

ENV DOMAIN_TLD docker
