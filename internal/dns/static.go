package dns

import (
	"net"
	"strings"
)

var staticSuffixes []string
var staticIPs [][]net.IP

func addStaticIP(name string, addrs []string) {
	var ips []net.IP
	for _, addr := range addrs {
		ips = append(ips, net.ParseIP(addr))
	}
	// use suffix point, because all DNS queries has it
	// use prefix point, because support subdomains by default
	staticSuffixes = append(staticSuffixes, "."+name+".")
	staticIPs = append(staticIPs, ips)
}

func lookupStaticIP(name string) ([]net.IP, error) {
	name = "." + name
	for i, suffix := range staticSuffixes {
		if strings.HasSuffix(name, suffix) {
			return staticIPs[i], nil
		}
	}
	return nil, nil
}
