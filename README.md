# Dinghy HTTP Proxy

This is the HTTP Proxy and DNS server that
[Dinghy](https://github.com/codekitchen/dinghy) uses.

The proxy is based on jwilder's excellent
[nginx-proxy](https://github.com/jwilder/nginx-proxy) project, with
modifications to make it more suitable for local development work.

A DNS resolver is also added. By default it will resolve all `*.docker` domains
to the Docker VM, but this can be changed.

## Configuration

### Exposed Ports

The proxy will by default use the first port exposed by your container as the
HTTP port to proxy to. This can be overridden by setting the VIRTUAL_PORT
environment variable on the container to the desired HTTP port.

### Docker Compose Projects

The proxy will auto-generate a hostname based on the docker tags that
docker-compose adds to each container. This hostname is of the form
`<service>.<project>.<tld>`. For instance, assuming the default `*.docker` TLD,
a "web" service in a "myapp" docker-compose project will be automatically made
available at http://web.myapp.docker/.

### Explicitly Setting a Hostname

As in the base nginx-proxy, you can configure a container's hostname by setting
the `VIRTUAL_HOST` environment variable in the container.

You can set the `VIRTUAL_HOST`
environment variable either with the `-e` option to docker or
the environment hash in docker-compose. For instance setting
`VIRTUAL_HOST=myrailsapp.docker` will make the container's exposed port
available at http://myrailsapp.docker/.

This will work even if dinghy auto-generates a hostname based on the
docker-compose tags.

#### Multiple Hosts

If you need to support multiple virtual hosts for a container, you can separate each entry with commas.  For example, `foo.bar.com,baz.bar.com,bar.com` and each host will be setup the same.

#### Wildcard Hosts

You can also use wildcards at the beginning and the end of host name, like `*.bar.com` or `foo.bar.*`. Or even a regular expression, which can be very useful in conjunction with a wildcard DNS service like [xip.io](http://xip.io), using `~^foo\.bar\..*\.xip\.io` will match `foo.bar.127.0.0.1.xip.io`, `foo.bar.10.0.2.2.xip.io` and all other given IPs. More information about this topic can be found in the nginx documentation about [`server_names`](http://nginx.org/en/docs/http/server_names.html).

### Enabling CORS

You can set the `CORS_ENABLED`
environment variable either with the `-e` option to docker or
the environment hash in docker-compose. For instance setting
`CORS_ENABLED=true` will allow the container's web proxy to accept cross domain
requests.

### Subdomain Support

If you want your container to also be available at all subdomains to the given
domain, prefix a dot `.` to the provided hostname. For instance setting
`VIRTUAL_HOST=.myrailsapp.docker` will also make your app avaiable at
`*.myrailsapp.docker`.

This happens automatically for the auto-generated docker-compose hostnames.

### SSL Support

SSL is supported using single host certificates using naming conventions.

To enable SSL, just put your certificates and privates keys in the ```HOME/.dinghy/certs``` directory
for any virtual hosts in use.  The certificate and keys should be named after the virtual host with a `.crt` and
`.key` extension.  For example, a container with `VIRTUAL_HOST=foo.bar.com.docker` should have a
`foo.bar.com.docker.crt` and `foo.bar.com.docker.key` file in the certs directory.

#### How SSL Support Works

The SSL cipher configuration is based on [mozilla nginx intermediate profile](https://wiki.mozilla.org/Security/Server_Side_TLS#Nginx) which
should provide compatibility with clients back to Firefox 1, Chrome 1, IE 7, Opera 5, Safari 1,
Windows XP IE8, Android 2.3, Java 7.  The configuration also enables HSTS, and SSL
session caches.

The default behavior for the proxy when port 80 and 443 are exposed is as follows:

* If a container has a usable cert, port 80 will redirect to 443 for that container so that HTTPS
is always preferred when available.
* If the container does not have a usable cert, port 80 will be used.

To serve traffic in both SSL and non-SSL modes without redirecting to SSL, you can include the
environment variable `HTTPS_METHOD=noredirect` (the default is `HTTPS_METHOD=redirect`).  You can also
disable the non-SSL site entirely with `HTTPS_METHOD=nohttp`.

#### How to quickly generate self-signed certificates

You can generate self-signed certificates using ```openssl```.

```bash
openssl req -x509 -newkey rsa:2048 -keyout foo.bar.com.docker.key \
-out foo.bar.com.docker.crt -days 365 -nodes \
-subj "/C=US/ST=Oregon/L=Portland/O=Company Name/OU=Org/CN=foo.bar.com.docker" \
-config <(cat /etc/ssl/openssl.cnf <(printf "[SAN]\nsubjectAltName=DNS:foo.bar.com.docker")) \
-reqexts SAN -extensions SAN
```

To prevent your browser to emit warning regarding self-signed certificates, you can install them on your system as trusted certificates.

## Using Outside of Dinghy

Since this functionality is generally useful for local development work even
outside of Dinghy, this proxy now supports running standalone.

#### Environment variables

We include a few environment variables to customize the proxy / dns server:

- `DOMAIN_TLD` default: `docker` - The DNS server will only respond to `*.docker` by default. You can change this to `dev` if it suits your workflow
- `DNS_IP` default: `127.0.0.1` - Setting this variable is explained below

### OS X

You'll need the IP of your VM:

* For docker-machine, run `docker-machine ip <machine_name>` to get the IP.
* For Docker for Mac, you can use `127.0.0.1` as the IP, since it forwards docker ports to the host machine.

Then start the proxy:

    docker run -d --restart=always \
      -v /var/run/docker.sock:/tmp/docker.sock:ro \
      -v ~/.dinghy/certs:/etc/nginx/certs \
      -p 80:80 -p 443:443 -p 19322:19322/udp \
      -e DNS_IP=<vm_ip> -e CONTAINER_NAME=http-proxy \
      --name http-proxy \
      codekitchen/dinghy-http-proxy

You will also need to configure OS X to use the DNS resolver. To do this, create
a file `/etc/resolver/docker` (creating the `/etc/resolver` directory if it does
not exist) with these contents:

```
nameserver <vm_ip>
port 19322
```

You only need to do this step once, or when the VM's IP changes.

### Linux

For running Docker directly on a Linux host machine, the proxy can still be
useful for easy access to your development environments. Similar to OS X, start
the proxy:

    docker run -d --restart=always \
      -v /var/run/docker.sock:/tmp/docker.sock:ro \
      -v ~/.dinghy/certs:/etc/nginx/certs \
      -p 80:80 -p 443:443 -p 19322:19322/udp \
      -e CONTAINER_NAME=http-proxy \
      --name http-proxy \
      codekitchen/dinghy-http-proxy

The `DNS_IP` environment variable is not necessary when Docker is running
directly on the host, as it defaults to `127.0.0.1`.

Different Linux distributions will require different steps for configuring DNS
resolution. The [Dory](https://github.com/FreedomBen/dory) project may be useful
here, it knows how to configure common distros for `dinghy-http-proxy`.
