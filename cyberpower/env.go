package cyberpower

import (
	"log"
)

type ENV struct {
	parent   *CyberPower
	name     string
	location string
	temp_f   float64
	humidity int64
	contact  []string
}

var env_path = "/env_status_update.html"

func (e *ENV) update() {
	root, err := e.parent.get(env_path)
	if err != nil {
		log.Printf("Unable to update ENV on %s", e.parent.hostpath)
	}
	log.Printf("type: %d; data: %s", root.Type, root.Data)
}
