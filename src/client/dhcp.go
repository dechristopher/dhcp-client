package client

import "github.com/dechristopher/dhcp-client/src/models"

/*
 * Takes the DHCPPacket and adds in the slices to it
 */
func finishDHCPPacket(packet *models.DHCPPacket) {
	packet.MessageType = packet.Data[0:1]
	packet.HardwareType = packet.Data[1:2]
	packet.HardwareAddressLength = packet.Data[2:3]
	packet.Hops = packet.Data[3:4]
	packet.TransactionID = packet.Data[4:8]
	packet.SecondsElapsed = packet.Data[8:10]
	packet.BootPFlags = packet.Data[10:12]
	packet.ClientIP = packet.Data[12:16]
	packet.YourIP = packet.Data[16:20]
	packet.NextServerIP = packet.Data[20:24]
	packet.RelayAgentIP = packet.Data[24:28]
	packet.ClientMAC = packet.Data[28:34]
	packet.ClientMACPadding = packet.Data[34:44]
	packet.ServerHostname = packet.Data[44:108]
	packet.BootFileName = packet.Data[108:236]
	packet.MagicCookie = packet.Data[236:240]
}

/*
 * Builds a basic DISCOVER packet with the MAC address supplied.
 * Most things are hardcoded such as the hardware type, hardware address length
 * (could be derived from the macAddress supplied) as well as using
 * 0.0.0.0 (could specify an interface later on)
 */
func BuildDiscoverPacket(macAddress [6]byte) models.DHCPPacket {
	packet := models.DHCPPacket{Data: []byte{}}
	// Message Type, Hardware Type, Hardware Address Length, Hops
	packet.Data = append(packet.Data, models.BootRequest, models.Ethernet, models.AddressLength, 0x0)
	// Transaction ID
	packet.Data = append(packet.Data, []byte{1, 2, 3, 4}...)
	// Seconds Elapsed
	packet.Data = append(packet.Data, []byte{0x0, 0x0}...)
	// BootP Flags
	packet.Data = append(packet.Data, models.Broadcast[:]...)
	// Client IP
	packet.Data = append(packet.Data, models.AllIPs[:]...)
	// Your IP
	packet.Data = append(packet.Data, models.AllIPs[:]...)
	// Next Server IP
	packet.Data = append(packet.Data, models.AllIPs[:]...)
	// Relay Agent IP
	packet.Data = append(packet.Data, models.AllIPs[:]...)
	// Client MAC
	packet.Data = append(packet.Data, macAddress[:]...)
	// Client MAC Padding
	packet.Data = append(packet.Data, models.MacPadding[:]...)
	// Server Hostname
	packet.Data = append(packet.Data, models.EmptyServerHostname[:]...)
	//Boot File Name
	packet.Data = append(packet.Data, models.EmptyBootFileName[:]...)
	// Magic Cookie
	packet.Data = append(packet.Data, models.MagicCookieDHCP[:]...)
	//DHCP Message Type
	packet.Data = append(packet.Data, []byte{0x35, 0x01, 0x01}[:]...)
	// Client Identifier
	packet.Data = append(packet.Data, make([]byte, 9)...)
	// Requested IP
	packet.Data = append(packet.Data, make([]byte, 4)...)
	// Hostname
	packet.Data = append(packet.Data, make([]byte, 13)...)
	// Vendor ID
	packet.Data = append(packet.Data, make([]byte, 10)...)
	// Parameter Request List
	packet.Data = append(packet.Data, make([]byte, 16)...)
	// End Code
	packet.Data = append(packet.Data, models.EndCode)
	// Padding
	packet.Data = append(packet.Data, make([]byte, 2)...)

	finishDHCPPacket(&packet)
	return packet
}

/*
 * Parses the offer packet returned from the DISCOVER request.
 */
func ParseOfferPacket(data []byte) models.DHCPPacket {
	packet := models.DHCPPacket{}
	packet.Data = append(packet.Data, data[:]...)
	finishDHCPPacket(&packet)
	return packet
}
