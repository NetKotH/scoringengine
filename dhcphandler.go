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
	ip     net.IP
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
	var ip *net.IP
	for id, VM := range victims {
		if VM.Mac == nic {
			ip = &VM.IP
			victims[id].LastSeen = time.Now()
			if victims[id].State == StateOffline {
				victims[id].State = StateRunning
			}
			break
		}
	}

	switch msgType {

	case dhcp.Discover:
		now := time.Now()
		if l, ok := h.leases[nic]; ok {
			if l.expiry.After(now) {
				ip = &l.ip
			}
		}
		if ip == nil {
			/* FIXME
			ip = h.findFree(nic)
			if ip == nil {
				logger.Println("Failed to find an IP")
				return nil
			}
			*/
			return nil
		}
		logger.Printf("Returning an IP: %s\n", ip)
		return dhcp.ReplyPacket(
			p,
			dhcp.Offer,
			h.ip,
			*ip,
			h.leaseDuration,
			h.options.SelectOrderOrAll(options[dhcp.OptionParameterRequestList]),
		)

	case dhcp.Request:
		logger.Printf("Got Request %#v\n", h.ip)
		if server, ok := options[dhcp.OptionServerIdentifier]; ok && !net.IP(server).Equal(h.ip) {
			return nil // Message not for this dhcp server
		}
		reqIP := net.IP(options[dhcp.OptionRequestedIPAddress])
		if reqIP == nil {
			reqIP = net.IP(p.CIAddr())
		}
		// If This is one of our victims
		if ip != nil && reqIP.Equal(*ip) {
			return dhcp.ReplyPacket(p, dhcp.ACK, h.ip, reqIP, h.leaseDuration,
				h.options.SelectOrderOrAll(options[dhcp.OptionParameterRequestList]))
		}
		return // FIXME abort if not a Victim

		if len(reqIP) == 4 && !reqIP.Equal(net.IPv4zero) {
			if leaseNum := dhcp.IPRange(h.start, reqIP) - 1; leaseNum >= 0 && leaseNum < h.leaseRange {
				if _, exists := h.leases[p.CHAddr().String()]; !exists {
					h.leases[p.CHAddr().String()] = lease{
						ip:     reqIP,
						expiry: time.Now().Add(h.leaseDuration),
					}
				} else {
					logger.Printf("Found an existing lease: %#v\n", reqIP)
					return dhcp.ReplyPacket(p, dhcp.ACK, h.ip, reqIP, h.leaseDuration,
						h.options.SelectOrderOrAll(options[dhcp.OptionParameterRequestList]))
				}
			}
		}
		//return dhcp.ReplyPacket(p, dhcp.NAK, h.ip, nil, 0, nil)

	case dhcp.Release, dhcp.Decline:
		logger.Println("Release or Decline")
		nic := p.CHAddr().String()
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
	logger.Printf("leases: %#v\n", h.leases)
	now := time.Now()
	for nic, l := range h.leases {
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
		for _, l := range h.leases {
			if l.ip.Equal(ip) {
				ip[3]--
				goto top
			}
		}
		return &ip
	}
	logger.Println("Out of IPs!")
	return nil
}
