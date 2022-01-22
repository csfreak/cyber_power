package cyberpower

import (
	"encoding/json"
	"net/http"
)

type device struct {
	Host        string
	UPS         *UPS
	Environment *ENV
}

var (
	cyberpowers []CyberPower
)

func RestGetHandler(res http.ResponseWriter, req *http.Request) {
	if !(req.Method == http.MethodGet) {
		res.WriteHeader(http.StatusMethodNotAllowed)
	} else {
		var outdata []device
		for _, c := range cyberpowers {
			d := device{
				UPS:         c.ups,
				Environment: c.env,
				Host:        c.hostpath[7 : len(c.hostpath)-1],
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
		res.Write(out)
	}
}
