nginx: nginx
dockergen: docker-gen -watch -only-exposed -notify "/app/reload-nginx" /app/nginx.tmpl /etc/nginx/conf.d/default.conf
dnsmasq: dnsmasq --no-daemon --port=19322 --no-resolv --address=/.$DOMAIN_TLD/$DNS_IP
