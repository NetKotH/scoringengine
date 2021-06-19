package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/brimstone/go-vboxmanage"
	"github.com/mostlygeek/arp"
	"github.com/spf13/viper"
)

func checkVMConfig(vm victim) {
	logger := log.New(os.Stderr, "["+vm.Name+"] ", log.LstdFlags)
	var err error
	var ok bool
	bridgeAdapter := viper.GetString("bridge")
	dirty := false
	// make Fresh snapshots of any VMs missing it
	var vmMode string
	if vmMode, ok = vm.VM.Meta["netkoth.mode"]; !ok {
		vmMode = "bridged"

	}
	if vmMode != "ignore" && vm.VM.Nic != vmMode {
		logger.Println("Modifying networking to bridged")
		// Restore to snapshot fresh
		vm.VM.RestoreSnapshot("Fresh")
		err = vm.VM.Modify("nic1", "bridged") // TODO set this to vmNic
		if err != nil {
			logger.Println("Unable to set network to bridged. Shutdown the VM if it's on.")
			return
		}
		dirty = true
	}
	if vmMode != "ignore" && vm.VM.Bridge != bridgeAdapter {
		logger.Println("Modifying bridged adapter to", bridgeAdapter)
		// Restore to snapshot fresh
		vm.VM.RestoreSnapshot("Fresh")
		vm.VM.Poweroff()
		err = vm.VM.Modify("bridgeadapter1", bridgeAdapter)
		if err != nil {
			panic(err)
		}
		dirty = true
	}
	if dirty {
		// Delete snapshot Fresh
		vm.VM.DeleteSnapshot("Fresh")
	}
	if vm.VM.Snapshots == nil {
		logger.Println("Making snapshot")
		vm.VM.MakeSnapshot("Fresh", "Fresh VM")
	}
	if vm.VM.Power == "off" {
		logger.Println("Starting VM")
		vm.VM.RestoreSnapshot("Fresh")
		vm.VM.Start()
	}
}

func refreshVMs() {
	var err error
	var vmsRaw []vboxmanage.VM
	// TODO toggle all VMs as unvisited
	// Start VM refresh loop
	vmsRaw, err = vboxmanage.ListVMs()
	if err != nil {
		log.Println("Error refreshing VMs:", err)
		return
	}
	needGroup := false
	for _, vmRaw := range vmsRaw {
		if vmRaw.Group == "NetKotH" {
			needGroup = true
			break
		}
	}
	freeIP := make(net.IP, len(bridgeIP))
	copy(freeIP, bridgeIP)
	freeIP[3]++
	for _, vmRaw := range vmsRaw {
		var nicID int64 = 1
		if nic, ok := vmRaw.Meta["netkoth.nic"]; ok {
			var err error
			nicID, err = strconv.ParseInt(nic, 10, 10)
			if err != nil {
				log.Println("Error parsing netkoth.nic value")
				continue
			}
			if nicID == 0 {
				nicID = 1
			}
		}
		if len(vmRaw.MAC) == 0 {
			continue
		}
		mac := vmRaw.MAC[nicID-1]

		if needGroup && vmRaw.Group != "NetKotH" {
			continue
		}
		if victims[mac] == nil {
			log.Printf("Adding victim: %#v\n", vmRaw)
			if staticIP, ok := vmRaw.Meta["ip"]; ok {
				freeIP = net.ParseIP(staticIP)
			} else {
				// Pick a usable ip
				for flag := false; flag == false; {
					log.Printf("Now trying ip: %d\n", freeIP[3])
					flag = true
					// First check the ARP table.
					for i, m := range arp.Table() {
						if i == freeIP.String() && strings.ToLower(m) != strings.ToLower(mac) {
							fmt.Println("Found ip in arp table")
							freeIP[3]++
							flag = true
							break
						}
					}
					// Next, check any other known victims
					for _, v := range victims {
						if v.IP.Equal(freeIP) {
							fmt.Println("Found ip in known victim lists")
							freeIP[3]++
							flag = true
							break
						}
					}
				}
			}
			vmIP := make(net.IP, len(freeIP))
			copy(vmIP, freeIP)
			victims[mac] = &victim{
				Controller: "<none>",
				IP:         vmIP,
				LastSeen:   time.Now(),
				Mac:        mac,
				Name:       vmRaw.Name,
				State:      StateOffline,
				Type:       "VM",
			}
		}
		victims[mac].VM = vmRaw
		/*
			if vmRaw.Power == "on" {
				victims[vmRaw.MAC].State = StateRunning
			}
		*/
		checkVMConfig(*victims[mac])
	}
	// TODO remove VMs that aren't visited
}
