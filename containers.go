package main

import (
	"log"
	"net"

	"github.com/fsouza/go-dockerclient"
)

func refreshContainers() {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		//log.Printf("Unable to refresh containers: %s\n", err)
		return
	}
	victimContainers, err := client.ListContainers(docker.ListContainersOptions{
		All:     true,
		Filters: map[string][]string{"label": {"victim"}},
	})
	if err != nil {
		log.Printf("Unable to refresh containers: %s\n", err)
		return
	}
	for _, c := range victimContainers {
		log.Printf("Found container: %s\n", c.Names[0])
		if victims[c.ID] == nil {
			victims[c.ID] = &victim{
				Controller: "<none>",
				IP:         net.ParseIP(c.Networks.Networks["bridge"].IPAddress).To4(),
				Mac:        c.ID,
				Name:       c.Names[0],
				State:      StateOffline,
				Type:       "container",
			}
		}
		if c.State == "running" {
			victims[c.ID].State = StateRunning
		}
	}

}
