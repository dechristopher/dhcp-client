package models

import "net"

const (
	BootRequest   = byte(1)
	Ethernet      = byte(1)
	AddressLength = byte(6)
	EndCode       = 0xff
)

var (
	// Unicast
	_                   = make([]byte, 2) // 00
	Broadcast           = [2]byte{0x80, 00}
	AllIPs              = net.IP{0, 0, 0, 0} // 0.0.0.0
	MacPadding          = make([]byte, 10)
	EmptyServerHostname = make([]byte, 64)
	EmptyBootFileName   = make([]byte, 128)
	MagicCookieDHCP     = [4]byte{0x63, 0x82, 0x53, 0x63}
)
