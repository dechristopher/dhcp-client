# Golang DHCP Client Demonstration

A toy golang DHCP client to explore the DHCP protocol.
Randomizes request MAC address as to not step on the same lease over and over.

## Usage
```bash
./compile.sh
./toydhcp [-ip4 address]
```

## DHCP Protocol Diagram

![Session](DHCP_session.png)
Wikipedia - Gelmo96 / CC BY-SA (https://creativecommons.org/licenses/by-sa/4.0)