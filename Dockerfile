FROM jwilder/nginx-proxy:0.3.6
MAINTAINER Brian Palmer <brian@codekitchen.net>

# XXX: this can be removed once upstream nginx-proxy is updated
ENV DOCKER_GEN_VERSION 0.6-groupByLabel
RUN wget https://github.com/codekitchen/docker-gen/releases/download/v0.6-groupByLabel/docker-gen-linux-amd64-$DOCKER_GEN_VERSION.tar.gz \
 && tar -C /usr/local/bin -xvzf docker-gen-linux-amd64-$DOCKER_GEN_VERSION.tar.gz \
 && rm docker-gen-linux-amd64-$DOCKER_GEN_VERSION.tar.gz
# /XXX

# override nginx configs
COPY *.conf /etc/nginx/conf.d/

# override nginx-proxy templating
COPY nginx.tmpl Procfile reload-nginx /app/
