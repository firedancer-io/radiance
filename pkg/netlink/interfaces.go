package netlink

import (
	"fmt"
	"net"
)

// GetInterfaceIP returns the primary IP address associated with the given interface.
func GetInterfaceIP(iface string) (net.IP, error) {
	ifaceAddr, err := net.InterfaceByName(iface)
	if err != nil {
		return nil, err
	}
	addrs, err := ifaceAddr.Addrs()
	if err != nil {
		return nil, err
	}
	for _, addr := range addrs {
		if ip, ok := addr.(*net.IPNet); ok && ip.IP.To4() != nil {
			return ip.IP, nil
		}
	}
	return nil, fmt.Errorf("no IP address found for interface %s", iface)
}
