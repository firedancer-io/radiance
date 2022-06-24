package main

import (
	"fmt"
	"net"

	"k8s.io/klog/v2"
)

func sendUDP(addr string, txb []byte, count int) {
	// Send UDP packet to TPU
	conn, err := net.Dial("udp", addr)
	if err != nil {
		// if we fail to open a UDP socket, something has gone really wrong
		klog.Exitf("failed to dial %s: %v", addr, err)
		return
	}
	defer conn.Close()
	tn := 0
	for i := 0; i < count; i++ {
		n, err := conn.Write(txb)
		if err != nil {
			klog.Errorf("failed to write to %s: %v", addr, err)
		}
		if n != len(txb) {
			panic(fmt.Errorf("wrote %d bytes, expected %d", n, len(txb)))
		}
		tn += n
	}
	klog.V(2).Infof("sent %d bytes to %s", tn, addr)
}
