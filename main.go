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
	// Desired IP Address
	requestedIP := flag.String("ip4", "", "Requested IPv4 address")
	flag.Parse()

	// Build discover packet, don't use actual interface IP here or actual lease will be returned
	discoverPacket := models.BuildDiscoverPacket([6]byte{0xA0, 0x99, 0x9B, 0x0C, 0xDE, 0xC8}, requestedIP)

	// Server Address is the broadcast address on port 67 (255.255.255.255:67)
	serverAddr, _ := net.ResolveUDPAddr("udp4",
		fmt.Sprintf("%s:67", net.IP{255, 255, 255, 255}))

	// Client address is 0.0.0.0:68 (all adapters) on the local machine
	clientAddr, err := net.ResolveUDPAddr("udp4",
		fmt.Sprintf("%s:4500", net.IP{0, 0, 0, 0}))
	if err != nil {
		fmt.Printf("Error: %+v", err)
		os.Exit(1)
	}
	fmt.Printf("Listen: %s:%d\n", clientAddr.IP, clientAddr.Port)

	conn, err := net.DialUDP("udp4", clientAddr, serverAddr)
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
	listener, err := net.ListenPacket("udp4", "0.0.0.0:68")
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
		_, _, err := listener.ReadFrom(respBuffer)
		if err != nil {
			fmt.Printf("UDP read error  %v", err)
			os.Exit(1)
		}
		responses <- respBuffer
	}()

	fmt.Printf("Sending DISCOVER\n\n")

	/*
	 * If the timeout is reached, the DISCOVER is assumed to have failed and
	 * a new one is sent, on a successful OFFER exit with success.
	 */
	for {
		// Now that we are listening for offers, send out a DISCOVER
		_, err = conn.Write(discoverPacket.Data)
		if err != nil {
			fmt.Printf("DISCOVER write error: %+v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Waiting for OFFER\n\n")
		// Channel waits for OFFER packet

		select {
		case response := <-responses:
			fmt.Println("OFFER received!")
			fmt.Printf("OFFER:\n%+v\n", models.ParsePacket(response))
			os.Exit(0)
		case <-time.After(2 * time.Second):
			fmt.Println("DISCOVER timeout, resending packet")
		}
	}
}
