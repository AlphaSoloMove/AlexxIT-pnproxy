**pnproxy** - Plug and Proxy is a simple home proxy for managing Internet traffic 

Features:

- work on all devices in the local network without additional settings
- proxy settings for selected sites only
- ad blocking support (like AdGuard)

Types:

- DNS proxy
- Reverse proxy for HTTP and TLS (level 4 proxy)
- HTTP anonymous proxy

## Install

- Binary - [nightly.link](https://nightly.link/AlexxIT/pnproxy/workflows/build/master)
- Docker - [alexxit/pnproxy](https://hub.docker.com/r/alexxit/pnproxy)
- Home Assistant Add-on - [alexxit/hassio-addons](https://github.com/AlexxIT/hassio-addons)

## Setup

For example, you want to block ads and also forward all Twitter traffic through external proxy server.
And want it to work on all home devices without additional configuration on each device.

1. Install pnproxy on any server in your home network (ex. IP: `192.168.1.123`).
   It is important that ports 53, 80 and 443 be free on this server.
2. Create `pnproxy.yaml`
   ```yaml
   hosts:
     adblock: doubleclick.net googlesyndication.com
     tunnel: twitter.com twimg.com t.co x.com
   
   dns:
     listen: ":53"
     rules:
       - name: adblock                         # name from hosts block
         action: static address 127.0.0.1      # block this sites
       - name: tunnel                          # name from hosts block
         action: static address 192.168.1.123  # redirect this sites to pnproxy
     default:
       action: dns server 8.8.8.8              # resolve DNS for all other sites
   
   http:
     listen: ":80"
     rules:
       - name: tunnel                          # name from hosts block
         action: redirect scheme https         # redirect this sites from HTTP to TLS module
     default:
       action: raw_pass
   
   tls:
     listen: ":443"
     rules:
       - name: tunnel                          # name from hosts block
         action: proxy_pass host 123.123.123.123 port 3128  # forward this sites to external HTTP proxy
     default:
       action: raw_pass
   
   proxy:
     listen: ":8080"                           # optionally run local HTTP proxy
   
   log:
     level: trace                              # optionally increase log level (default - info)
   ```
3. Setup DNS server for your home router to `192.168.1.123`.

Optionally, instead of step 3, you can verify that everything works by configuring an HTTP proxy to `192.168.1.123:8080` on your PC or mobile device.

## Configuration

By default, the app looks for the `pnproxy.yaml` file in the current working directory.

```shell
pnproxy -config /config/pnproxy.yaml
```

By default all modules disabled and don't listen any ports.

## Module: Hosts

Store lists of site domains for use in other modules.

- Name comparison includes all subdomains, you don't need to specify them separately!
- Names can be written with spaces or line breaks. Follow [YAML syntax](https://yaml-multiline.info/).

```yaml
hosts:
  list1: site1.com site2.com site3.net
  list2: |
    site1.com static.site1.cc
    site2.com cdnsite2.com
    site3.in site3.com site3.co.uk
```

## Module: DNS

Run DNS server and act as DNS proxy.

- Can protect from MITM DNS attack using [DNS over TLS](https://en.wikipedia.org/wiki/DNS_over_TLS) or [DNS over HTTPS](https://en.wikipedia.org/wiki/DNS_over_HTTPS) 
- Can work as AdBlock like [AdGuard](https://adguard.com/)

Enable server:

```yaml
dns:
  listen: ":53"
```

Rules action supports setting `static address` only:

- Useful for ad blocking.
- Useful for routing some sites traffic through pnproxy.

```yaml
dns:
  rules:
    - name: adblocklist
      action: static address 127.0.0.1
    - name: list1 list2 site4.com site5.net
      action: static address 192.168.1.123
```

Default action supports [DNS](https://en.wikipedia.org/wiki/Domain_Name_System), [DOT](https://en.wikipedia.org/wiki/DNS_over_TLS) and [DOH](https://en.wikipedia.org/wiki/DNS_over_HTTPS) upstream:

- Important to use server IP-address, instead of a domain name

```yaml
dns:
  default:
    # action - dns or dot or doh
    action: dns server 8.8.8.8
```

Support build-in providers - `cloudflare`, `google`, `quad9`, `opendns`, `yandex`:

- all this providers support DNS, DOH and DOT technologies.

```yaml
dns:
  default:
    action: dot provider google
```

Total config:

```yaml
dns:
  listen: ":53"
  rules:
    - name: adblocklist
      action: static address 127.0.0.1
    - name: list1 list2 site4.com site5.net
      action: static address 192.168.1.123
  default:
    action: doh provider cloudflare
```

## Module: HTTP

Run HTTP server and act as reverse proxy.

Enable server:

```yaml
http:
  listen: ":80"
```

Rules action supports setting `redirect scheme https` with optional code:

- Useful for redirect all sites traffic to TLS module.

```yaml
http:
  rules:
    - name: list1 list2 site4.com site5.net
      # code - any number (default - 307)
      action: redirect scheme https
```

Rules action supports setting `raw_pass`:

```yaml
http:
  rules:
    - name: list1 list2 site4.com site5.net
      action: raw_pass
```

Rules action supports setting `proxy_pass`:

- Useful for passing all sites traffic to additional local or remote proxy.

```yaml
http:
  rules:
    - name: list1 list2 site4.com site5.net
      # host and port - mandatory
      # username and password - optional
      # type - socks5 (default - http)
      action: proxy_pass host 123.123.123.123 port 3128 username user1 password pasw1
```

Default action support all rules actions:

```yaml
http:
  default:
    action: raw_pass
```

## Module: TLS

Run TCP server and act as Layer 4 reverse proxy.

Enable server:

```yaml
tls:
  listen: ":443"
```

Rules action supports setting `raw_pass`:

- Useful for forward HTTPS traffic to another reverse proxies with custom port.

```yaml
tls:
  rules:
    - name: list1 list2 site4.com site5.net
      # host - optional rewrite connection IP-address
      # port - optional rewrite connection port
      action: raw_pass host 123.123.123.123 port 10443
```

Rules action supports setting `proxy_pass`:

- Useful for passing all sites traffic to additional local or remote proxy.

```yaml
tls:
  rules:
    - name: list1 list2 site4.com site5.net
      # host and port - mandatory
      # username and password - optional
      # type - socks5 (default - http)
      action: proxy_pass host 123.123.123.123 port 3128 username user1 password pasw1
```

Rules action supports setting `split_pass`:

- Can try to protect from hardware MITM HTTPS attack.

```yaml
tls:
  rules:
    - name: list1 list2 site4.com site5.net
      action: split_pass
```

Default action support all rules actions:

```yaml
tls:
  default:
    action: raw_pass
```

## Module: Proxy

Run HTTP proxy server. This module does not have its own rules. It uses the HTTP and TLS module rules.
You can choose not to run DNS, HTTP, and TLS servers and use pnproxy only as HTTP proxy server.

Enable server:

```yaml
proxy:
  listen: ":8080"
```

## Tips and Tricks

**Mikrotik DNS fail over script**

- Add as System > Scheduler > Interval `00:01:00`

```
:global server "192.168.1.123"

:do {
  :resolve google.com server $server
} on-error={
  :global server "8.8.8.8"
}

:if ([/ip dns get servers] != $server) do={
  /ip dns set servers=$server
}
```

## Known bugs

In rare cases, due to [HTTP/2 connection coalescing](https://blog.cloudflare.com/connection-coalescing-experiments) technology, some site may not work properly when using a TCP/TLS Layer 4 proxy. In HTTP proxy mode everything works fine. Everything works fine in Safari browser (it doesn't support this technology). In Firefox, this feature can be disabled - `network.http.http2.coalesce-hostnames`.
