package cyberpower

import (
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

type CyberPower interface {
	login() bool
	logout()
	get(path string) (*html.Node, error)
	getHost() string
	getEnv() (*ENV, bool)
	getUps() (*UPS, bool)
	logged_in() bool
	update() error
}

type CyberPowerModule interface {
	update() error
	getParent() CyberPower
}

type CP struct {
	hostpath   string
	loginForm  url.Values
	client     http.Client
	ups        CyberPowerModule
	env        CyberPowerModule
	_logged_in bool
}

type ENV struct {
	parent   CyberPower
	Name     string
	Location string
	Temp_f   float64
	Humidity int
}

type UPS struct {
	parent  CyberPower
	Input   INPUT_POWER
	Output  OUTPUT_POWER
	Battery BATTERY_POWER
	Temp_c  int
	Temp_f  int
	Status  string
}

type INPUT_POWER struct {
	Status    string
	Voltage   float64
	Frequency float64
}

type OUTPUT_POWER struct {
	Status      string
	Voltage     float64
	Frequency   float64
	Current     float64
	LoadWatts   int
	LoadPercent int
}

type BATTERY_POWER struct {
	Status            string
	RemainingCapacity int
	RemainingRuntime  int
}
