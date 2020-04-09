package models

import (
	"encoding/binary"
	"fmt"
	"net"
)

/*
 * DHCPPacket contains a byte array for storing all types of DHCP packets
 * Currently it also stores a static list of slices to easily access different
 * parts of the packet.
 * These slices could later be put into a map to allow the structure to be more
 * generic to allow all the different DHCP options available.
 */
type DHCPPacket struct {
	Data                  []byte          // Actual Data
	OpCode                []byte          // 1 Byte 	[0]
	HardwareType          []byte          // 1 Byte 	[1]
	HardwareAddressLength []byte          // 1 Byte 	[2]
	Hops                  []byte          // 1 Byte 	[3]
	TransactionID         []byte          // 4 Bytes 	[4-7]
	SecondsElapsed        []byte          // 2 Bytes 	[8-9]
	BootPFlags            []byte          // 2 Bytes 	[10-11]
	ClientIP              []byte          // 4 Bytes 	[12-15]
	YourIP                []byte          // 4 Bytes 	[16-19]
	ServerIP              []byte          // 4 Bytes 	[20-23]
	RelayAgentIP          []byte          // 4 Bytes	[24-27]
	ClientMAC             []byte          // 6 Bytes	[28-33]
	ClientMACPadding      []byte          // 10 Bytes	[34-43]
	ServerHostname        []byte          // 64 Bytes	[44-107]
	BootFileName          []byte          // 128 Bytes 	[108-235]
	MagicCookie           []byte          // 4 Bytes	[236-239]
	DHCPMessageType       DHCPMessageType // 3 Bytes	[240-243]
	Options               map[int]string  // DHCP options by code
}

/*
 * Returns the DHCP packet type of the current packet
 */
func (p DHCPPacket) Type() DHCPMessageType {
	switch p.Data[243] {
	case 0x01:
		return DISCOVER
	case 0x02:
		return OFFER
	case 0x03:
		return REQUEST
	case 0x05:
		return ACKNOWLEDGE
	case 0x06:
		return NACKNOWLEDGE
	default:
		return INVALID
	}
}

/*
 * Pretty-print a DHCP packet
 */
func (p DHCPPacket) String() string {
	return fmt.Sprintf(
		"--------------------------------------------------------------\n"+
			"Message Type: %x\n"+
			"Hardware Type: %x\n"+
			"Hardware Address Length: %x\n"+
			"Hops: %x\n"+
			"Transaction ID: %v\n"+
			"Seconds Elapsed: %v\n"+
			"Bootp Flags: %v\n"+
			"Client IP: %v\n"+
			"Your IP: %v\n"+
			"Next Server IP: %v\n"+
			"Relay Agent IP: %v\n"+
			"Client MAC: [% X]\n"+
			"Client MAC Padding: %v\n"+
			"Server Hostname: \"%s\"\n"+
			"BootFileName: ~REDACTED~\n"+
			"Magic Cookie: %X\n"+
			"DHCP Message Type: %v\n"+
			"Additional DHCP Options: \n%v\n"+
			"--------------------------------------------------------------\n",
		p.OpCode,
		p.HardwareType,
		p.HardwareAddressLength,
		p.Hops,
		p.TransactionID,
		p.SecondsElapsed,
		p.BootPFlags,
		p.ClientIP,
		p.YourIP,
		p.ServerIP,
		p.RelayAgentIP,
		p.ClientMAC,
		p.ClientMACPadding,
		string(p.ServerHostname),
		p.MagicCookie,
		p.DHCPMessageType,
		p.Options)
}

/*
 * Parses the offer packet returned from the DISCOVER request.
 */
func ParsePacket(data []byte) DHCPPacket {
	packet := DHCPPacket{}
	packet.Data = append(packet.Data, data[:]...)
	finishDHCPPacket(&packet)
	parseOptions(&packet)
	return packet
}

/*
 * Takes the DHCPPacket and populates the struct with the proper
 * packet field slices
 */
func finishDHCPPacket(packet *DHCPPacket) {
	packet.OpCode = packet.Data[0:1]
	packet.HardwareType = packet.Data[1:2]
	packet.HardwareAddressLength = packet.Data[2:3]
	packet.Hops = packet.Data[3:4]
	packet.TransactionID = packet.Data[4:8]
	packet.SecondsElapsed = packet.Data[8:10]
	packet.BootPFlags = packet.Data[10:12]
	packet.ClientIP = packet.Data[12:16]
	packet.YourIP = packet.Data[16:20]
	packet.ServerIP = packet.Data[20:24]
	packet.RelayAgentIP = packet.Data[24:28]
	packet.ClientMAC = packet.Data[28:34]
	packet.ClientMACPadding = packet.Data[34:44]
	packet.ServerHostname = packet.Data[44:108]
	packet.BootFileName = packet.Data[108:236]
	packet.MagicCookie = packet.Data[236:240]
	packet.Options = make(map[int]string)
}

/*
 * Loop through all DHCP options and set them in the packet struct
 */
func parseOptions(packet *DHCPPacket) {
	cursor := 240

	for {
		option := packet.Data[cursor]
		cursor++

		if option == 255 {
			break
		}

		length := packet.Data[cursor]
		cursor++

		var optionData []byte

		if length == 1 {
			optionData = []byte{packet.Data[cursor]}
			cursor++
		} else {
			var dataEnd = cursor + (int(length))
			optionData = packet.Data[cursor:dataEnd]
			cursor += int(length)
		}

		// Log option, length, and data
		//fmt.Printf("Option: %v | Length: %v | Data: %v\n",
		//	option, length, optionData)

		parseOption(packet, option, optionData)
	}
}

func parseOption(packet *DHCPPacket, option byte, optionData []byte) {
	switch option {
	case 53: // DHCP Message Type
		packet.DHCPMessageType = ParseDHCPMessageType(optionData[0])
		break
	case 54: // DHCP Server Address
		packet.Options[54] = fmt.Sprintf("(DHCP Server: %s)",
			net.IP(optionData))
		break
	case 51: //DHCP Lease Time
		packet.Options[51] = fmt.Sprintf("(Lease: %d sec)",
			int(binary.BigEndian.Uint32(optionData)))
		break
	case 58: // Renewal (T1) time
		packet.Options[58] = fmt.Sprintf("(Renewal: %d sec)",
			int(binary.BigEndian.Uint32(optionData)))
		break
	case 59: // Rebinding (T2) time
		packet.Options[59] = fmt.Sprintf("(Rebinding: %d sec)",
			int(binary.BigEndian.Uint32(optionData)))
		break
	case 28: // Broadcast address
		packet.Options[28] = fmt.Sprintf("(Broadcast Address: %s)",
			net.IP(optionData))
		break
	case 6: // DNS Server Address
		packet.Options[6] = fmt.Sprintf("(DNS Server: %s)",
			net.IP(optionData))
		break
	case 15: // DNS Domain Name
		packet.Options[15] = fmt.Sprintf("(DNS Domain Name: %s)",
			string(optionData))
		break
	case 1: // Subnet Mask
		packet.Options[1] = fmt.Sprintf("(Subnet Mask: %s)",
			net.IP(optionData))
		break
	case 3: // Router Address
		packet.Options[3] = fmt.Sprintf("(Router: %s)",
			net.IP(optionData))
		break
	}
}
