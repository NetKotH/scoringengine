package main

import (
	"log"
	"net"
	"os"
	"strings"
	"time"

	dhcp "github.com/krolaw/dhcp4"
	"github.com/mostlygeek/arp"
)

type lease struct {
	nic    string
	expiry time.Time // When the lease expires
}

type DHCPHandler struct {
	ip            net.IP           // Server IP to use
	options       dhcp.Options     // Options to send to DHCP Clients
	start         net.IP           // Start of IP range to distribute
	leaseRange    int              // Number of IPs to distribute (starting from start)
	leaseDuration time.Duration    // Lease period
	leases        map[string]lease // Map to keep track of leases
}

func (h *DHCPHandler) ServeDHCP(p dhcp.Packet, msgType dhcp.MessageType, options dhcp.Options) (d dhcp.Packet) {
	nic := strings.ToUpper(p.CHAddr().String())
	logger := log.New(os.Stderr, "[dhcp] ["+nic+"] ", log.LstdFlags)
	var vip *net.IP
	// Update victim table
	for id, VM := range victims {
		if VM.Mac == nic {
			vip = &VM.IP
			victims[id].LastSeen = time.Now()
			if victims[id].State == StateOffline {
				victims[id].State = StateRunning
			}
			break
		}
		vip = nil
	}

	switch msgType {

	case dhcp.Discover:
		if vip != nil {
			return dhcp.ReplyPacket(
				p,
				dhcp.Offer,
				h.ip,
				*vip,
				h.leaseDuration,
				h.options.SelectOrderOrAll(options[dhcp.OptionParameterRequestList]),
			)
		}
		/* FIXME */
		ip := h.findFree(nic)
		if ip == nil {
			logger.Println("Failed to find an IP")
			return nil
		}
		//logger.Println("Returning early from Discover")
		//return nil
		logger.Printf("Returning an IP: %s\n", ip)
		h.leases[ip.String()] = lease{
			nic:    nic,
			expiry: time.Now().Add(h.leaseDuration),
		}
		return dhcp.ReplyPacket(
			p,
			dhcp.Offer,
			h.ip,
			*ip,
			h.leaseDuration,
			h.options.SelectOrderOrAll(options[dhcp.OptionParameterRequestList]),
		)

	case dhcp.Request:
		logger.Printf("Got Request to %#v\n", h.ip)
		if server, ok := options[dhcp.OptionServerIdentifier]; ok && !net.IP(server).Equal(h.ip) {
			logger.Println("Message not for this server")
			return nil // Message not for this dhcp server
		}
		/*
			reqIP := net.IPv4zero
			copy(reqIP, net.IP(options[dhcp.OptionRequestedIPAddress]))
			if reqIP == nil {
				logger.Println("Using CIAddr")
				copy(reqIP, net.IP(p.CIAddr()))
			}
		*/
		reqIP := net.IP(options[dhcp.OptionRequestedIPAddress])
		if reqIP == nil {
			logger.Println("Using CIAddr")
			reqIP = net.IP(p.CIAddr())
		}
		// If This is one of our victims
		if vip != nil && reqIP.Equal(*vip) {
			logger.Println("This is a victim")
			return dhcp.ReplyPacket(p, dhcp.ACK, h.ip, reqIP, h.leaseDuration,
				h.options.SelectOrderOrAll(options[dhcp.OptionParameterRequestList]))
		}
		logger.Println("not a victim")

		if len(reqIP) != 4 {
			logger.Printf("IP isn't ipv4 %d %s\n", len(reqIP), reqIP)
			// maybe ? return dhcp.ReplyPacket(p, dhcp.NAK, h.ip, nil, 0, nil)
			return nil
		}
		if reqIP.Equal(net.IPv4zero) {
			logger.Println("IP is 0")
			// maybe ? return dhcp.ReplyPacket(p, dhcp.NAK, h.ip, nil, 0, nil)
			return nil
		}
		logger.Println("Reasonable IP request")
		leaseNum := dhcp.IPRange(h.start, reqIP) - 1
		logger.Printf("leaseNum: %#v\n", leaseNum)
		if leaseNum < 0 || leaseNum > h.leaseRange {
			return nil
		}
		if _, exists := h.leases[reqIP.String()]; !exists {
			nic := p.CHAddr().String()
			logger.Printf("New client %s = %s", nic, reqIP)
			h.leases[reqIP.String()] = lease{
				nic:    nic,
				expiry: time.Now().Add(h.leaseDuration),
			}
			logger.Println("leases")
			for nic, l := range h.leases {
				logger.Printf("%s: %#v\n", nic, l)
			}
		} else {
			logger.Printf("Found an existing lease: %#v\n", reqIP)
		}
		return dhcp.ReplyPacket(p, dhcp.ACK, h.ip, reqIP, h.leaseDuration,
			h.options.SelectOrderOrAll(options[dhcp.OptionParameterRequestList]))

	case dhcp.Release, dhcp.Decline:
		nic := p.CHAddr().String()
		logger.Printf("Release or Decline from %s", nic)
		if _, ok := h.leases[nic]; ok {
			delete(h.leases, nic)
		}
	default:
		logger.Println("default")
	}
	logger.Println("Final return!")
	return nil
}

func (h *DHCPHandler) findFree(mac string) *net.IP {
	logger := log.New(os.Stderr, "[dhcp] ", log.LstdFlags)
	ip := dhcp.IPAdd(h.start, h.leaseRange)
	startip := dhcp.IPAdd(h.start, len(victims))
	logger.Printf("ip: %s\n", ip)
	logger.Printf("startip: %s\n", startip)
	logger.Printf("leases:\n")
	now := time.Now()
	for nic, l := range h.leases {
		logger.Printf("%s: %#v\n", nic, l)
		// Expired leases
		if l.expiry.Before(now) {
			delete(h.leases, nic)
		}
	}
	for n, m := range arp.Table() {
		// If we already have an entry in the arp table for the ip, return that
		if strings.ToLower(m) == strings.ToLower(mac) {
			newip := net.ParseIP(n)
			return &newip
		}
	}
	for ; ip[3] > startip[3]; ip[3]-- {
	top:
		for i := range h.leases {
			if i == ip.String() {
				ip[3]--
				goto top
			}
		}
		return &ip
	}
	logger.Println("Out of IPs!")
	return nil
}
