package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/netip"
	"os"

	"github.com/things-go/go-socks5"
	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun/netstack"
)

type tnetResolver struct {
	tnet *netstack.Net
}

func (r *tnetResolver) Resolve(ctx context.Context, name string) (context.Context, net.IP, error) {
	addrs, err := r.tnet.LookupContextHost(ctx, name)
	if err != nil {
		return nil, nil, err
	}
	for _, addr := range addrs {
		ip, err := netip.ParseAddr(addr)
		if err != nil {
			continue
		}
		return ctx, ip.AsSlice(), nil
	}
	return ctx, nil, nil
}

func setDeviceConfig(wg *device.Device, name string) error {
	f, err := os.Open(name)
	if err != nil {
		return fmt.Errorf("open config: %w", err)
	}
	defer f.Close()
	return wg.IpcSetOperation(f)
}

func run() error {
	configFilename := flag.String("config", "/etc/socks5-wireguard-proxy/config", "Path to config file")
	wireguardAddress := flag.String("wireguard-address", "", "WireGuard interface address")
	dnsServer := flag.String("dns-server", "1.1.1.1", "DNS server address")
	listenAddress := flag.String("listen-address", "0.0.0.0:1080", "Address to listen on")
	flag.Parse()

	tunDev, tnet, err := netstack.CreateNetTUN(
		[]netip.Addr{netip.MustParseAddr(*wireguardAddress)},
		[]netip.Addr{netip.MustParseAddr(*dnsServer)},
		1420,
	)
	if err != nil {
		return fmt.Errorf("create net tun: %w", err)
	}

	wg := device.NewDevice(tunDev, conn.NewDefaultBind(), device.NewLogger(device.LogLevelError, "wg"))

	if err := setDeviceConfig(wg, *configFilename); err != nil {
		return fmt.Errorf("set config: %w", err)
	}

	if err := wg.Up(); err != nil {
		return fmt.Errorf("up: %w", err)
	}

	s := socks5.NewServer(
		socks5.WithLogger(socks5.NewLogger(log.New(os.Stdout, "socks5: ", log.LstdFlags))),
		socks5.WithResolver(&tnetResolver{tnet: tnet}),
		socks5.WithDial(tnet.DialContext),
	)

	log.Printf("listening on %s", *listenAddress)

	if err := s.ListenAndServe("tcp", *listenAddress); err != nil {
		return fmt.Errorf("serve: %w", err)
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
