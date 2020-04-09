package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/dechristopher/dhcp-client/src/client"
)

func main() {
	// Build discover packet
	discoverPacket := client.BuildDiscoverPacket([6]byte{0xE4, 0xB3, 0x18, 0xCA, 0x84, 0x83})

	// Server Address is the broadcast address on port 67 (255.255.255.255:67)
	serverAddr, _ := net.ResolveUDPAddr("udp",
		fmt.Sprintf("%s:67", net.IP{255, 255, 255, 255}))

	// Client address is 0.0.0.0:68 (all adapters) on the local machine
	clientAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:68")
	if err != nil {
		fmt.Printf("Error: %+v", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", clientAddr, serverAddr)
	// Defer UDP connection close so we can handle errors on close
	defer func() {
		if conn != nil {
			err := conn.Close()
			if err != nil {
				fmt.Printf("UDP connection close error: %+v\n", err)
			}
		}
	}()

	if err != nil {
		fmt.Printf("UDP dial error: %+v", err)
		os.Exit(1)
	}

	// Channel for responses
	responses := make(chan []byte)
	// UDP listener
	listener, err := net.ListenUDP("udp", clientAddr)
	// Defer UDP listener close so we can handle errors on close
	defer func() {
		if listener != nil {
			err := listener.Close()
			if err != nil {
				fmt.Printf("UDP listen close error: %+v\n", err)
			}
		}
	}()

	if err != nil {
		fmt.Printf("Listener error: %+v\n", err)
		os.Exit(1)
	}

	/*
	 * Goroutine to listen for DHCP OFFER
	 *
	 * This doesn't HAVE to be a goroutine but starting to wait before sending
	 * the discover packet removes the chance that we aren't ready to receive
	 * by the time it arrives.
	 *
	 * Also allows for adding easier retry handling later
	 */
	go func() {
		respBuffer := make([]byte, 2048)
		_, _, err := listener.ReadFromUDP(respBuffer)
		if err != nil {
			fmt.Printf("UDP read error  %v", err)
			os.Exit(1)
		}
		responses <- respBuffer
	}()

	fmt.Printf("Sending DISCOVER packet\n\n")

	/*
	 * If the timeout is reached, the DISCOVER is assumed to have failed and
	 * a new one is sent, on a successful OFFER exit with success.
	 */
	for {
		// Now that we are listening for offers, send out a DISCOVER
		_, err = conn.Write(discoverPacket.Data)
		if err != nil {
			fmt.Printf("Discover write error: %+v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Waiting for OFFER packet\n\n")
		// Channel waits for OFFER packet

		select {
		case response := <-responses:
			fmt.Println("OFFER received!")
			fmt.Printf("OFFER Packet: %+v\n", client.ParseOfferPacket(response))
			os.Exit(0)
		case <-time.After(2 * time.Second):
			fmt.Println("DISCOVER timeout, resending packet")
		}
	}
}
