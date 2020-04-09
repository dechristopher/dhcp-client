package models

/*
 * Builds a basic REQUEST packet for responding to the DHCP OFFER
 */
func BuildRequestPacket(macAddress [6]byte, requestedIP []byte, server []byte) DHCPPacket {
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
	// Server IP
	packet.Data = append(packet.Data, server[:]...)
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
	//DHCP Message Type - Option 53
	packet.Data = append(packet.Data, []byte{0x35, 0x01, 0x03}...)
	// Requested IP - Option 50
	packet.Data = append(packet.Data, []byte{0x32, 0x04}...)
	packet.Data = append(packet.Data, requestedIP[:]...)
	// DHCP Server - Option 54
	packet.Data = append(packet.Data, []byte{0x36, 0x4}...)
	packet.Data = append(packet.Data, server[:]...)
	// End Code
	packet.Data = append(packet.Data, EndCode)

	finishDHCPPacket(&packet)
	parseOptions(&packet)
	return packet
}
