package cyberpower

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"

	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
)

type CyberPower struct {
	hostpath  string
	loginForm url.Values
	client    http.Client
	ups       UPS
	env       ENV
}

func FromENV() *CyberPower {
	c := NewCyberPower(os.Getenv("CYBERPOWER_HOST"), os.Getenv("CYBERPOWER_USERNAME"), os.Getenv("CYBERPOWER_PASSWORD"))
	if c == nil {
		log.Print("unable to create cyberpower from environment variables")
		return nil
	}
	return c
}

func NewCyberPower(host string, username string, password string) *CyberPower {
	c := &CyberPower{}
	c.hostpath = "http://" + host
	c.loginForm = url.Values{}
	c.loginForm.Set("username", username)
	c.loginForm.Set("password", password)
	j, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}
	c.client = http.Client{
		Jar: j,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	if !(c.login()) {
		return nil
	}
	c.ups = UPS{
		parent: c,
	}
	c.env = ENV{
		parent: c,
	}
	cyberpowers = append(cyberpowers, *c)
	return c
}

func (c *CyberPower) get(path string) (*html.Node, error) {
	resp, err := c.client.Get(c.hostpath + path)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusSeeOther || resp.StatusCode == http.StatusForbidden {
		if c.login() {
			resp, err = c.client.Get(c.hostpath + path)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("unable to login")
		}
	}
	defer resp.Body.Close()
	node, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (c *CyberPower) Update() {
	c.ups.update()
	c.env.update()
}
