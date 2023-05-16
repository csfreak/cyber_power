package cyberpower

import (
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

type CyberPower interface {
	login() error
	logout() error
	get(path string) (*html.Node, error)
	getHost() string
	getEnv() (*ENV, bool)
	getUps() (*UPS, bool)
	loggedIn() bool
	update() error
}

type CPModule interface {
	update() error
	getParent() CyberPower
}

type CP struct {
	hostpath  string
	loginForm url.Values
	client    http.Client
	ups       CPModule
	env       CPModule
	_loggedIn bool
}

type ENV struct {
	parent   CyberPower
	Name     string
	Location string
	TempF    float64
	Humidity int
}

type UPS struct {
	parent  CyberPower
	Input   InputPower
	Output  OutputPower
	Battery BatteryPower
	TempC   int
	TempF   int
	Status  string
}

type InputPower struct {
	Status    string
	Voltage   float64
	Frequency float64
}

type OutputPower struct {
	Status      string
	Voltage     float64
	Frequency   float64
	Current     float64
	LoadWatts   int
	LoadPercent int
}

type BatteryPower struct {
	Status            string
	RemainingCapacity int
	RemainingRuntime  int
}
