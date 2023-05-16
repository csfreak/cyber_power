package main

import (
	"log"
	"net/http"
	"os"

	"github.com/csfreak/cyber_power/cyberpower"
)

var apiListenerPort = ":8080"

// var configfilepath = "/etc/cyberpower/config.yaml"

func main() {
	if os.Getenv("CYBERPOWER_PORT") != "" {
		apiListenerPort = ":" + os.Getenv("CYBERPOWER_PORT")
	}

	var devices []cyberpower.CyberPower
	if c, err := cyberpower.FromENV(true); err != nil {
		devices = append(devices, c)
	}

	/*if os.Getenv("CYBERPOWER_CONFIG") != "" {
		configfilepath = os.Getenv("CYBERPOWER_CONFIG")
	}*/

	/*for _, conf := range cyberpower.ReadConfig(configfilepath).Cyberpower {
		devices = append(devices, cyberpower.NewCyberPower(conf.Host, conf.Username, conf.Password))
	}*/

	if len(devices) == 0 {
		log.Fatal("unable to find any devices")
	}

	http.HandleFunc("/v1/cyberpower", cyberpower.RestGetHandler)
	http.HandleFunc("/v1/cyberpower/", cyberpower.RestGetHandler)

	log.Printf("Starting HTTP Server on %s", apiListenerPort)

	// Log and Exit if http server exits
	log.Fatal(http.ListenAndServe(apiListenerPort, nil))
}
