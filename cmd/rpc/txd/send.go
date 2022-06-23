package main

import (
	"fmt"
	"net"

	"k8s.io/klog/v2"
)

func sendUDP(addr string, txb []byte) {
	// Send UDP packet to TPU
	conn, err := net.Dial("udp", addr)
	if err != nil {
		klog.Errorf("failed to dial %s: %v", addr, err)
	}
	defer conn.Close()
	n, err := conn.Write(txb)
	if err != nil {
		klog.Errorf("failed to write to %s: %v", addr, err)
	}
	if n != len(txb) {
		panic(fmt.Errorf("wrote %d bytes, expected %d", n, len(txb)))
	}
	klog.V(2).Infof("sent %d bytes to %s", n, addr)
}
