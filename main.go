package main

// Unused for now
//go:generate go get -v github.com/omeid/go-resources/cmd/resources
//go:generate resources -output=embed.go -var=Assets -tag=embed -trim=static/ static/*

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	vboxmanage "github.com/brimstone/go-vboxmanage"
	dhcp "github.com/krolaw/dhcp4"
	"github.com/krolaw/dhcp4/conn"
	apachelog "github.com/lestrrat-go/apache-logformat"
	"github.com/miekg/dns"
	"github.com/spf13/viper"
)

type StateType string

const (
	StateOffline = "offline"
	StateRunning = "online"
	StateCrashed = "crashed"
)

type victim struct {
	Controller  string
	IP          net.IP
	LastSeen    time.Time
	LastSeenRel time.Duration
	Mac         string
	Name        string
	Type        string
	VM          vboxmanage.VM
	State       StateType
	visited     bool
}

var victims map[string]*victim

type templateVars struct {
	Teams   []TeamScore
	Victims []victim
}
type TeamScore struct {
	Name  string
	Score int
}

type EventType int

const (
	//EventBeacon    = 1
	//EventFoundFlag = 2
	EventPlantFlag = 1
)

type Event struct {
	Team   string
	Type   EventType
	Victim string
	When   time.Time
}

var Events []*Event

// unused for now
// var Assets http.FileSystem

var teamRe *regexp.Regexp

var bridgeIP net.IP
var scoreinterval time.Duration

func main() {
	logger := log.New(os.Stderr, "[main] ", log.LstdFlags)

	// Handle signals
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		s := <-sig
		logger.Printf("Got signal: %#v\n", s)
		os.Exit(0)
	}()

	//vboxmanage.Loglevel = 0

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetDefault("bridge", "eth0")
	viper.SetDefault("scoreinterval", "1m")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Println("Unable to load config file, using defaults")
	}

	teamRe = regexp.MustCompile("<team>([A-Za-z0-9]+)</team>")

	victims = make(map[string]*victim)

	l, _ := apachelog.New(`%h %l %u %t "%r" %>s %b "%{Referer}i" "%{User-agent}i" "%v"`)

	// Setup a custom mux for timing
	muxs := http.NewServeMux()
	// Setup routes
	muxs.HandleFunc("/", handleRoot)
	// Setup an actual https server
	srvs := &http.Server{
		Addr:           ":443",
		Handler:        l.Wrap(muxs, os.Stderr),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		TLSConfig:      setupTLS(),
		TLSNextProto:   make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	// Figure out the ip of the bridge interface
	iface, err := net.InterfaceByName(viper.GetString("bridge"))
	addrs, err := iface.Addrs()
	// handle err
	if err != nil {
		logger.Print("Perhaps the bridge address doesn't exist?")
		panic(err)
	}
	for _, addr := range addrs {
		switch v := addr.(type) {
		case *net.IPNet:
			if len(v.Mask) == 16 {
				continue
			}
			bridgeIP = v.IP
		case *net.IPAddr:
			bridgeIP = v.IP
		}
		if bridgeIP.To4() == nil {
			continue
		}
	}

	bridgeIP = bridgeIP.To4()
	// If the bridge IP is unset
	if bridgeIP.Equal(net.IPv4zero) || bridgeIP == nil {
		logger.Fatalf("No IP found on bridge interface %s\n", iface)
	}
	logger.Printf("bridgeIP: %s\n", bridgeIP)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRootRedirect)
	srv := &http.Server{
		Addr:           ":80",
		Handler:        l.Wrap(mux, os.Stderr),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	scoreinterval, err = time.ParseDuration(viper.GetString("scoreinterval"))
	dhcphandler := &DHCPHandler{
		ip:            bridgeIP,
		leaseDuration: scoreinterval,
		start:         bridgeIP,
		leaseRange:    50,
		leases:        make(map[int]lease, 10),
		options: dhcp.Options{
			dhcp.OptionSubnetMask:       []byte{255, 255, 240, 0},
			dhcp.OptionDomainNameServer: []byte(bridgeIP), // Presuming Server is also your DNS server
		},
	}

	// Setup dhcp server
	go func() {
		logger.Printf("Starting DHCP server on %s\n", viper.GetString("bridge"))
		connection, err := conn.NewUDP4BoundListener(viper.GetString("bridge"), ":67")
		if err != nil {
			// echo out setcap command
			logger.Printf("Error starting DHCP Server: %s\n", err)
			logger.Println("Does this have the right capabilities?")
			// TODO make sure this wins the race against the other servers
			logger.Fatalf("setcap cap_net_bind_service,cap_net_raw+ep %s\n", os.Args[0])
		}
		err = dhcp.Serve(connection, dhcphandler)
		if err != nil {
			logger.Fatalf("Error with DHCP Server: %s\n", err)
		}
	}()

	// Setup http server
	go func() {
		logger.Println("Starting http service")
		// Start the server
		err = srv.ListenAndServe()
		if err != nil {
			logger.Fatalf("Error starting server: %s\n", err)
		}
		os.Exit(1)
	}()

	// Setup https server
	go func() {
		logger.Println("Starting https service")
		// Start the server
		err = srvs.ListenAndServeTLS("", "")
		if err != nil {
			logger.Fatalf("Error starting server: %s\n", err)
		}
		os.Exit(1)
	}()

	// Setup dns server
	go func() {
		logger.Println("Starting dns service")
		srv := &dns.Server{Addr: ":53", Net: "udp"}
		srv.Handler = &dnshandler{}
		if err := srv.ListenAndServe(); err != nil {
			logger.Fatalf("Failed to set udp listener %s\n", err.Error())
		}
	}()

	// Give services a chance to start
	logger.Println("Waiting for services to start before starting VMs")
	time.Sleep(time.Second * 2)
	logger.Println("Starting victim checks")
	// Start target system checks in the foreground
	for {
		refreshVMs()
		refreshContainers()
		checkVictims()
		// TODO make this a ticker instead
		time.Sleep(scoreinterval)
	}
}
