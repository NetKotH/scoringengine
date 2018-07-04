package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func checkTCP(ip string, port int) bool {
	log.Println("Should be checking the tcp port")
	return true
}

func checkTeam(ip net.IP) (string, error) {
	url := "http://" + ip.String()
	res, err := http.Get(url)
	if err != nil {
		return "", errors.New("Error trying to retrieve " + url + ": " + err.Error())
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "", errors.New("Error trying to read body from " + url + ": " + err.Error())
	}
	teamMatch := teamRe.FindStringSubmatch(string(body))
	if len(teamMatch) != 2 {
		// FIXME maybe this shouldn't be an error
		// return "", errors.New("Team name too short")
		return "", nil
	}
	return teamMatch[1], nil
}

func checkVictims() {
	// check victims
	for i := range victims {
		logger := log.New(os.Stderr, "["+victims[i].Name+"] ", log.LstdFlags)
		if victims[i].State == StateCrashed {
			// Start the VM if it has crashed
			victims[i].VM.Start()
			continue
		}

		// If the victim is a VM that hasn't checked in in a while,
		//logger.Printf("Last seen: %s\n", victims[i].LastSeen)
		if victims[i].LastSeen.Add(scoreinterval * 3).Before(time.Now()) {
			logger.Printf("Last seen: %s\n", victims[i].LastSeen)
			// Mark it as crashed and reset it
			if victims[i].State == StateRunning {
				logger.Println("Marking as crashed")
				victims[i].State = StateCrashed
				victims[i].VM.Poweroff()
				continue
			}
		}

		team, err := checkTeam(victims[i].IP)
		if err != nil {
			victims[i].State = StateOffline
			logger.Printf("Error checking: %s\n", err)
			continue
		}
		victims[i].State = StateRunning
		if team == "" {
			victims[i].Controller = "<none>"
			logger.Printf("uncontrolled\n")
			continue
		}
		victims[i].Controller = team
		Events = append(Events, &Event{
			Team:   team,
			Type:   EventPlantFlag,
			Victim: i,
			When:   time.Now(),
		})
		logger.Printf("controlled by %s\n", team)
	}
}
