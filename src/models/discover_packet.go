package models

import "net"

/*
 * Builds a basic DISCOVER packet with the MAC address supplied.
 * Most things are hardcoded such as the hardware type, hardware address length
 * (could be derived from the macAddress supplied) as well as using
 * 0.0.0.0 (could specify an interface later on)
 */
func BuildDiscoverPacket(macAddress [6]byte, requestedIP *string) DHCPPacket {
	packet := DHCPPacket{Data: []byte{}}
	// Message Type
	packet.Data = append(packet.Data, BootRequest)
	// Hardware Type
	packet.Data = append(packet.Data, Ethernet)
	// Hardware address length
	packet.Data = append(packet.Data, AddressLength)
	// Hop count
	packet.Data = append(packet.Data, 0x0)
	// Transaction ID
	packet.Data = append(packet.Data, []byte{1, 2, 3, 4}...)
	// Seconds Elapsed
	packet.Data = append(packet.Data, []byte{0x0, 0x0}...)
	// BootP Flags
	packet.Data = append(packet.Data, Broadcast[:]...)
	// Client IP
	packet.Data = append(packet.Data, AllIPs[:]...)
	// Your IP
	packet.Data = append(packet.Data, AllIPs[:]...)
	// Next Server IP
	packet.Data = append(packet.Data, AllIPs[:]...)
	// Relay Agent IP
	packet.Data = append(packet.Data, AllIPs[:]...)
	// Client MAC
	packet.Data = append(packet.Data, macAddress[:]...)
	// Client MAC Padding
	packet.Data = append(packet.Data, MacPadding[:]...)
	// Server Hostname
	packet.Data = append(packet.Data, EmptyServerHostname[:]...)
	//Boot File Name
	packet.Data = append(packet.Data, EmptyBootFileName[:]...)
	// Magic Cookie
	packet.Data = append(packet.Data, MagicCookieDHCP[:]...)
	//DHCP Message Type
	packet.Data = append(packet.Data, []byte{0x35, 0x01, 0x01}[:]...)

	// Requested IP option
	if *requestedIP != "" {
		packet.Data = append(packet.Data, []byte{0x32, 0x04}[:]...)
		packet.Data = append(packet.Data, []byte(net.ParseIP(*requestedIP))[12:16]...)
	}

	// End Code
	packet.Data = append(packet.Data, EndCode)

	finishDHCPPacket(&packet)
	return packet
}
