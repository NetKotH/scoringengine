package main

import (
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"

	dhcp "github.com/krolaw/dhcp4"
)

type lease struct {
	nic    string    // Client's CHAddr
	expiry time.Time // When the lease expires
}

type DHCPHandler struct {
	ip            net.IP        // Server IP to use
	options       dhcp.Options  // Options to send to DHCP Clients
	start         net.IP        // Start of IP range to distribute
	leaseRange    int           // Number of IPs to distribute (starting from start)
	leaseDuration time.Duration // Lease period
	leases        map[int]lease // Map to keep track of leases
}

func (h *DHCPHandler) ServeDHCP(p dhcp.Packet, msgType dhcp.MessageType, options dhcp.Options) (d dhcp.Packet) {
	logger := log.New(os.Stderr, "[dhcp] ", log.LstdFlags)
	nic := strings.ToUpper(p.CHAddr().String())
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
		free := -1
		logger.Printf("Got a Discover from %s\n", nic)
		if ip == nil {
			for i, v := range h.leases { // Find previous lease
				if v.nic == nic {
					free = i
					break
				}
			}
			if free == -1 {
				free = h.freeLease()
			}
			if free == -1 {
				logger.Println("Early return!")
				return
			}
			logger.Printf("Free: %d, start: %s\n", free, h.start)
			newip := dhcp.IPAdd(h.start, free)
			ip = &newip
			logger.Printf("ip: %#v %s\n", newip, newip)
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

		if len(reqIP) == 4 && !reqIP.Equal(net.IPv4zero) {
			if leaseNum := dhcp.IPRange(h.start, reqIP) - 1; leaseNum >= 0 && leaseNum < h.leaseRange {
				if l, exists := h.leases[leaseNum]; !exists || l.nic == p.CHAddr().String() {
					h.leases[leaseNum] = lease{nic: p.CHAddr().String(), expiry: time.Now().Add(h.leaseDuration)}
					logger.Printf("Found an existing lease: %#v\n", reqIP)
					return dhcp.ReplyPacket(p, dhcp.ACK, h.ip, reqIP, h.leaseDuration,
						h.options.SelectOrderOrAll(options[dhcp.OptionParameterRequestList]))
				}
			}
		}
		return dhcp.ReplyPacket(p, dhcp.NAK, h.ip, nil, 0, nil)

	case dhcp.Release, dhcp.Decline:
		logger.Println("Release or Decline")
		nic := p.CHAddr().String()
		for i, v := range h.leases {
			if v.nic == nic {
				delete(h.leases, i)
				break
			}
		}
	default:
		logger.Println("default")
	}
	logger.Println("Final return!")
	return nil
}

func (h *DHCPHandler) freeLease() int {
	now := time.Now()
	b := rand.Intn(h.leaseRange) // Try random first
	for _, v := range [][]int{[]int{b, h.leaseRange}, []int{0, b}} {
		for i := v[0]; i < v[1]; i++ {
			if l, ok := h.leases[i]; !ok || l.expiry.Before(now) {
				return i
			}
		}
	}
	return -1
}
