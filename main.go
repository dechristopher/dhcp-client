package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/dechristopher/dhcp-client/src/models"
)

func main() {
	// Desired IP Address command line argument
	requestedIP := flag.String("ip4", "", "Requested IPv4 address")
	flag.Parse()

	// Random MAC Address
	sampleMac := [6]byte{0xA0, 0x99, 0x9B, 0x0C, 0xDE, 0xC8}

	// Build discover packet, don't use actual interface MAC here or actual
	// computer lease will be returned from DHCP server
	discoverPacket := models.BuildDiscoverPacket(sampleMac, requestedIP)

	// Server Address is the broadcast address on port 67 (255.255.255.255:67)
	serverAddr, _ := net.ResolveUDPAddr("udp4",
		fmt.Sprintf("%s:67", net.IP{255, 255, 255, 255}))

	// Client address should be 0.0.0.0:68 (all adapters) on the local machine
	clientAddr, _ := net.ResolveUDPAddr("udp4",
		fmt.Sprintf("%s:68", net.IP{0, 0, 0, 0}))

	// Open UDP socket to DHCP server
	conn, err := net.ListenUDP("udp4", clientAddr)
	// Defer UDP connection close so we can handle errors on close
	defer func() {
		if conn != nil {
			err := conn.Close()
			if err != nil {
				fmt.Printf("UDP connection close error: %+v\n", err)
			}
		}
	}()
	// Make sure UDP dial doesn't fail
	if err != nil {
		fmt.Printf("UDP dial error: %+v", err)
		os.Exit(1)
	}

	// Channel for responses
	responses := make(chan []byte)

	/*
	 * Goroutine to listen for DHCP packets
	 *
	 * This doesn't HAVE to be a goroutine but starting to wait before sending
	 * the discover packet removes the chance that we aren't ready to receive
	 * by the time it arrives.
	 *
	 * Also allows for adding easier retry handling later
	 */
	go func() {
		for {
			respBuffer := make([]byte, 2048)
			b, _, err := conn.ReadFrom(respBuffer)
			fmt.Printf("READ: %d bytes\n\n", b)
			//_, _, err := listener.ReadFrom(respBuffer)
			if err != nil {
				fmt.Printf("UDP read error  %v", err)
				os.Exit(1)
			}
			responses <- respBuffer
		}
	}()

	// Timeout and re-send DISCOVER after 5 seconds
	timeout := time.NewTicker(time.Second * 5)

	/*
	 * If the timeout is reached, the DISCOVER is assumed to have failed and
	 * a new one is sent, on a successful OFFER exit with success.
	 */
	for {
		// Now that we are listening for offers, send out a DISCOVER
		_, err = conn.WriteTo(discoverPacket.Data, serverAddr)
		if err != nil {
			fmt.Printf("DISCOVER write error: %+v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Sent DISCOVER\n\n")

		fmt.Printf("Waiting for OFFER..\n\n")
		// Channel waits for OFFER packet

		// Wait until we have a response
		select {
		case response := <-responses:
			// Pause timeout since we're in a flow
			timeout.Stop()

			// Parse offer packet
			offer := models.ParsePacket(response)
			fmt.Println("OFFER received!")
			fmt.Printf("OFFER:\n%+v\n", offer)

			// Print what server offered
			fmt.Printf("Server offered: %s", net.IP(offer.YourIP))
			if net.IP(offer.YourIP).String() == *requestedIP {
				fmt.Printf(" (our requested IP!)\n")
			} else {
				fmt.Println()
			}

			// Build the REQUEST packet given the server's response
			request := models.BuildRequestPacket(sampleMac, offer.YourIP, offer.ServerIP)

			// Send REQUEST
			_, err := conn.WriteTo(request.Data, serverAddr)
			if err != nil {
				fmt.Printf("REQUEST write error: %+v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Sent REQUEST for IP: %s\n\n", net.IP(offer.YourIP))

			// Pull ACK MESSAGE
			response = <-responses

			// Parse ACK packet
			ack := models.ParsePacket(response)

			// Make sure we have a positive acknowledge
			if ack.DHCPMessageType == models.ACKNOWLEDGE {
				fmt.Println("ACK received!")
				fmt.Printf("ACK: \n%+v\n", ack)
				fmt.Printf("Leased IP: %s", net.IP(ack.YourIP))
			} else {
				fmt.Println("Uh-oh! NACK received!")
				fmt.Printf("NACK: \n%+v\n", ack)
			}

			os.Exit(0)
		case <-timeout.C:
			fmt.Println("DISCOVER timeout, resending DISCOVER")
		}
	}
}
