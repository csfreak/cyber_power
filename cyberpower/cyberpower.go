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

func FromENV(validate bool) (CyberPower, error) {
	c, err := NewCyberPower(os.Getenv("CYBERPOWER_HOST"), os.Getenv("CYBERPOWER_USERNAME"), os.Getenv("CYBERPOWER_PASSWORD"), validate)
	if err != nil {
		log.Print("unable to create cyberpower from environment variables")
		return c, err
	}
	return c, nil
}

func NewCyberPower(host string, username string, password string, validate bool) (CyberPower, error) {
	c := &CP{}
	c.hostpath = "http://" + host
	c.loginForm = url.Values{}
	c.loginForm.Set("action", "LOGIN")
	c.loginForm.Set("username", username)
	c.loginForm.Set("password", password)
	j, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return c, err
	}
	c.client = http.Client{
		Jar: j,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	if validate {
		if !(c.login()) {
			return c, fmt.Errorf("unable to login to %s as %s", c.hostpath, username)
		}
	}

	c.ups = &UPS{
		parent: c,
	}
	c.env = &ENV{
		parent: c,
	}
	cyberpowers = append(cyberpowers, c)
	return c, nil
}

func (c *CP) get(path string) (*html.Node, error) {
	if !(c._logged_in) {
		if !(c.login()) {
			return nil, fmt.Errorf("unable to login")
		}
	}
	resp, err := c.client.Get(c.hostpath + path)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	node, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (c *CP) logged_in() bool {
	return c._logged_in
}

func (c *CP) getHost() string {
	return c.hostpath[7:]
}

func (c *CP) getEnv() (*ENV, bool) {
	env, ok := c.env.(*ENV)
	return env, ok

}

func (c *CP) getUps() (*UPS, bool) {
	ups, ok := c.ups.(*UPS)
	return ups, ok
}

func (c *CP) update() error {
	err := c.ups.update()
	if err != nil {
		return err
	}
	err = c.env.update()
	if err != nil {
		return err
	}
	c.logout()
	return nil
}
