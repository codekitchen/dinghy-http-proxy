# Dinghy HTTP Proxy

This is the HTTP Proxy and DNS server that
[Dinghy](https://github.com/codekitchen/dinghy) uses.

The proxy is based on jwilder's excellent
[nginx-proxy](https://github.com/jwilder/nginx-proxy) project, with
modifications to make it more suitable for local development work.

A DNS resolver is also added. By default it will resolve all `*.docker` domains
to the Docker VM, but this can be changed.

As in the base nginx-proxy, you can configure a container's hostname by setting
the `VIRTUAL_HOST` environment variable in the container. In addition, this
proxy also auto-creates hostnames for docker-compose projects. The format is
`<container_name>.<compose_project_name>.<tld>`. For example, for a container
named `web` in a docker-compose project named `myapp`, you can visit
http://web.myapp.docker to be proxied to that container, without setting
`VIRTUAL_HOST`.

## Using Outside of Dinghy

Since this functionality is generally useful for local development work even
outside of Dinghy, this proxy now supports running standalone.

### OS X

You'll need the IP of your VM:

* For docker-machine, run `docker-machine ip <machine_name>` to get the IP.
* For Docker for Mac, you can use `127.0.0.1` as the IP, since it forwards docker ports to the host machine.

Then start the proxy:

    docker run -d --restart=always \
      -v /var/run/docker.sock:/tmp/docker.sock:ro \
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
      -p 80:80 -p 443:443 -p 19322:19322/udp \
      -e CONTAINER_NAME=http-proxy \
      --name http-proxy \
      codekitchen/dinghy-http-proxy

The `DNS_IP` environment variable is not necessary when Docker is running
directly on the host, as it defaults to `127.0.0.1`.

Different Linux distributions will require different steps for configuring DNS
resolution. The [Dory](https://github.com/FreedomBen/dory) project may be useful
here, it knows how to configure common distros for `dinghy-http-proxy`.
