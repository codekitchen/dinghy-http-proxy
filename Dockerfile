FROM jwilder/nginx-proxy
MAINTAINER Brian Palmer <brian@codekitchen.net>

ADD *.conf /etc/nginx/conf.d/
