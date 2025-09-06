# SOCKS5 WireGuard Proxy

Specifically designed to proxy WireGuard traffic over a SOCKS5 proxy in
user space.

## Usage

It can be run as a rootless container, so it can be run in Kubernetes with no
permissions at all.

Use the image
`ghcr.io/uhthomas/socks5-wireguard-proxy:latest`.

Provide the wireguard interface address with `--wireguard-address`.

The proxy looks for a config file at `/etc/socks5-wireguard-proxy/config`, but
can be changed with `--config`.

The config must be in [IPC
format](https://www.wireguard.com/xplatform/#configuration-protocol).

For example:

```text
private_key=<hex encoded private key>
public_key=<hex encoded public key>
allowed_ip=0.0.0.0/0
endpoint=192.168.1.1:51820
```

I learned the hard way that ordering of these IPC commands matters, and the
endpoint must come last or it won't work.
