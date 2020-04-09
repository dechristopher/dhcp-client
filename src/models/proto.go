package models

import "fmt"

// DHCPPacket contains a byte array for storing all types of DHCP packets
// Currently it also stores a static list of slices to easily access different parts of the packet
// These slices could later be put into a map to allow the structure to be more generic to allow all the different DHCP options availiable
type DHCPPacket struct {
	Data                  []byte // Actual Data
	MessageType           []byte // 1 Byte 		[0]
	HardwareType          []byte // 1 Byte 		[1]
	HardwareAddressLength []byte // 1 Byte 		[2]
	Hops                  []byte // 1 Byte 		[3]
	TransactionID         []byte // 4 Bytes 	[4-7]
	SecondsElapsed        []byte // 2 Bytes 	[8-9]
	BootPFlags            []byte // 2 Bytes 	[10-11]
	ClientIP              []byte // 4 Bytes 	[12-15]
	YourIP                []byte // 4 Bytes 	[16-19]
	NextServerIP          []byte // 4 Bytes 	[20-23]
	RelayAgentIP          []byte // 4 Bytes		[24-27]
	ClientMAC             []byte // 6 Bytes		[28-33]
	ClientMACPadding      []byte // 10 Bytes	[34-43]
	ServerHostname        []byte // 64 Bytes	[44-107]
	BootFileName          []byte // 128 Bytes 	[108-235]
	MagicCookie           []byte // 4 Bytes		[236-239]
	DHCPMessageType       []byte // 3 Bytes		[240-252]
	ClientIdentifier      []byte // 9 Bytes		[253-261]
	RequestedIP           []byte // 4 Bytes		[262-265]
	Hostname              []byte // 13 Bytes	[266-278]
	VendorID              []byte // 10 Bytes	[279-288]
	ParameterRequestList  []byte // 16 Bytes	[289-304]
	End                   []byte // 1 Byte		[305]
	Padding               []byte // 2 Bytes		[306-307]
}

func (p DHCPPacket) String() string {
	return fmt.Sprintf("Message Type: %x\nHardware Type: %x\nHardware Address Length: %x\nHops: %x\nTransaction ID: % v\nSeconds Elapsed: % v\nBootp Flags: % v\nClient IP: % v\nYour IP: % v\nNext Server IP: % v\nRelay Agent IP: % v\nClient MAC: [% X]\nClient MAC Padding: %v \nServer Hostname:%s\nBootFileName:\n%v\nMagic Cookie: % X\nDHCP Message Type: % v\nClient Identifier: % v\nRequested IP: % v\nHostname: %v\nVendor ID: % v\nParameter Request List: % v\n", p.MessageType, p.HardwareType, p.HardwareAddressLength, p.Hops, p.TransactionID, p.SecondsElapsed, p.BootPFlags, p.ClientIP, p.YourIP, p.NextServerIP, p.RelayAgentIP, p.ClientMAC, p.ClientMACPadding, string(p.ServerHostname), p.BootFileName, p.MagicCookie, p.DHCPMessageType, p.ClientIdentifier, p.RequestedIP, p.Hostname, p.VendorID, p.ParameterRequestList)
}
