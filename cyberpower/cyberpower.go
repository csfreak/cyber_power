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
	ehost, ok := os.LookupEnv("CYBERPOWER_HOST")
	if !ok {
		return nil, fmt.Errorf("unable to load CYBERPOWER_HOST")
	}

	euser, ok := os.LookupEnv("CYBERPOWER_USERNAME")
	if !ok {
		return nil, fmt.Errorf("unable to load CYBERPOWER_USERNAME")
	}

	epw, ok := os.LookupEnv("CYBERPOWER_PASSWORD")
	if !ok {
		return nil, fmt.Errorf("unable to load CYBERPOWER_PASSWORD")
	}

	c, err := NewCyberPower(ehost, euser, epw, validate)
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
	c.loginForm.Set("action", LoginAction)
	c.loginForm.Set("username", username)
	c.loginForm.Set("password", password)

	j, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		// This Error currently cannot be hit, as there is no return in cookiejar.New to return an error.
		return c, fmt.Errorf("unable to create cookiejar: %w", err)
	}

	c.client = http.Client{
		Jar: j,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	if validate {
		if err := c.login(); err != nil {
			return c, fmt.Errorf("unable to login to %s as %s; %w", c.hostpath, username, err)
		}
	}

	c.ups = &UPS{parent: c}
	c.env = &ENV{parent: c}

	cyberpowers = append(cyberpowers, c)

	return c, nil
}

func (c *CP) get(path string) (*html.Node, error) {
	if !(c._loggedIn) {
		if err := c.login(); err != nil {
			return nil, fmt.Errorf("unable to login: %w", err)
		}
	}

	resp, err := c.client.Get(c.hostpath + path)
	if err != nil {
		return nil, fmt.Errorf("unable to get %s: %w", path, err)
	}

	defer resp.Body.Close()

	node, err := html.Parse(resp.Body)
	if err != nil {
		// This should never be hit, as it only passes along read errors on resp.Body.Read()
		// Invalid HTML will not cause an error
		return nil, fmt.Errorf("unable to parse response body from %s: %w", path, err)
	}

	return node, nil
}

func (c *CP) loggedIn() bool {
	return c._loggedIn
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
		return fmt.Errorf("unable to update ups: %w", err)
	}

	err = c.env.update()
	if err != nil {
		return fmt.Errorf("unable to update env: %w", err)
	}

	err = c.logout()
	if err != nil {
		return fmt.Errorf("unable to logout: %w", err)
	}

	return nil
}
