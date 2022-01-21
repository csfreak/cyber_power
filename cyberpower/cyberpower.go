package cyberpower

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
)

type CyberPower struct {
	hostpath  string
	loginForm url.Values
	client    *http.Client
	ups       *UPS
	env       *ENV
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
	c.client = &http.Client{
		Jar: j,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	if !(c.login()) {
		log.Panic("unable to login")
	}
	c.ups = &UPS{
		parent: c,
	}
	c.env = &ENV{
		parent: c,
	}

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
	log.Println(c.ups)
	c.env.update()
}
