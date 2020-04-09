package models

import (
	"net"
)

const (
	BootRequest   = byte(1)
	Ethernet      = byte(1)
	AddressLength = byte(6)
	EndCode       = 0xff
)

var (
	//Unicast           = make([]byte, 2) // 00
	Broadcast           = [2]byte{0x80, 00}
	AllIPs              = net.IP{0, 0, 0, 0} // 0.0.0.0
	MacPadding          = make([]byte, 10)
	EmptyServerHostname = make([]byte, 64)
	EmptyBootFileName   = make([]byte, 128)
	MagicCookieDHCP     = [4]byte{0x63, 0x82, 0x53, 0x63}
)

type DHCPMessageType int

const (
	DISCOVER     DHCPMessageType = 1
	OFFER        DHCPMessageType = 2
	REQUEST      DHCPMessageType = 3
	ACKNOWLEDGE  DHCPMessageType = 5
	NACKNOWLEDGE DHCPMessageType = 6
	INVALID      DHCPMessageType = 0
)

/*
 * Parse the DHCP Message Type from a given byte
 */
func ParseDHCPMessageType(mType byte) DHCPMessageType {
	switch mType {
	case byte(DISCOVER):
		return DISCOVER
	case byte(OFFER):
		return OFFER
	case byte(REQUEST):
		return REQUEST
	case byte(ACKNOWLEDGE):
		return ACKNOWLEDGE
	case byte(NACKNOWLEDGE):
		return NACKNOWLEDGE
	default:
		return INVALID
	}
}

/*
 *
 */
func (t DHCPMessageType) String() string {
	switch t {
	case DISCOVER:
		return "DISCOVER"
	case OFFER:
		return "OFFER"
	case REQUEST:
		return "REQUEST"
	case ACKNOWLEDGE:
		return "ACKNOWLEDGE"
	case NACKNOWLEDGE:
		return "NACKNOWLEDGE"
	default:
		return "INVALID"
	}
}
