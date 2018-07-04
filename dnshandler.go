package main

import (
	"log"
	"os"

	"github.com/miekg/dns"
)

type dnshandler struct{}

func (this *dnshandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	logger := log.New(os.Stderr, "[dns] ", log.LstdFlags)
	msg := dns.Msg{}
	msg.SetReply(r)
	switch r.Question[0].Qtype {
	case dns.TypeA:
		msg.Authoritative = true
		// Actually, we don't care about the domain name, everything goes to the bridgeIP
		domain := msg.Question[0].Name
		logger.Printf("Request for %s\n", domain)
		msg.Answer = append(msg.Answer, &dns.A{
			Hdr: dns.RR_Header{
				Name:   domain,
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    60,
			},
			A: bridgeIP,
		})
	}
	w.WriteMsg(&msg)
}
