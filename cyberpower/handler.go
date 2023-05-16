package cyberpower

import (
	"encoding/json"
	"log"
	"net/http"
)

type device struct {
	Host        string
	UPS         UPS
	Environment ENV
}

var cyberpowers []CyberPower

func RestGetHandler(res http.ResponseWriter, req *http.Request) {
	if !(req.Method == http.MethodGet) {
		res.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		var outdata []device
		for _, c := range cyberpowers {
			err := c.update()
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			e, ok := c.getEnv()
			if !ok {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			u, ok := c.getUps()
			if !ok {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			d := device{
				UPS:         *u,
				Environment: *e,
				Host:        c.getHost(),
			}

			outdata = append(outdata, d)
		}
		out, err := json.Marshal(outdata)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		_, err = res.Write(out)
		if err != nil {
			log.Printf("unable to write response: %v", err)
		}
	}
}
