package main

import (
	"log"
	"net/http"
	"os"

	"github.com/csfreak/cyber_power/cyberpower"
)

var apiListenerPort = ":8080"
var configfilepath = "/etc/cyberpower/config.yaml"

func main() {
	var devices []*cyberpower.CyberPower
	devices = append(devices, cyberpower.FromENV())

	if os.Getenv("CYBERPOWER_CONFIG") != "" {
		configfilepath = os.Getenv("CYBERPOWER_CONFIG")
	}

	if os.Getenv("CYBERPOWER_PORT") != "" {
		apiListenerPort = ":" + os.Getenv("CYBERPOWER_PORT")
	}

	for _, conf := range read_config(configfilepath).cyberpower {
		devices = append(devices, cyberpower.NewCyberPower(conf.host, conf.username, conf.password))
	}

	for i, c := range devices {
		if c == nil && len(devices) > 1 {
			devices = append(devices[:i], devices[i+1:]...)
		} else if c == nil && len(devices) <= 1 {
			devices = make([]*cyberpower.CyberPower, 0)
		}
	}
	if len(devices) == 0 {
		log.Fatal("unable to find any devices")
	}

	http.HandleFunc("/v1/cyberpower", cyberpower.RestGetHandler)
	http.HandleFunc("/v1/cyberpower/", cyberpower.RestGetHandler)

	log.Printf("Starting HTTP Server on %s", apiListenerPort)

	//Log and Exit if http server exits
	log.Fatal(http.ListenAndServe(apiListenerPort, nil))
}
