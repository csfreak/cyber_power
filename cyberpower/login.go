package cyberpower

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func (c *CP) login() error {
	resp, err := c.client.Get(c.hostpath + "/login.html")
	if err != nil {
		return fmt.Errorf("unable to get login.html: %w", err)
	}

	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read login.html: %w", err)
	}

	resp, err = c.client.PostForm(c.hostpath+"/login_pass.cgi", c.loginForm)
	if err != nil {
		return fmt.Errorf("unable to post login_pass.cgi: %w", err)
	}

	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read login_pass.cgi: %w", err)
	}

	// resp, err = c.client.Get(c.hostpath + "/login_pass.html?action=LOGIN&username=" + c.loginForm.Get("username") + "&password=" + c.loginForm.Get("password"))
	resp, err = c.client.Get(c.hostpath + "/login_pass.html")
	if err != nil {
		return fmt.Errorf("unable to get login_pass.html: %w", err)
	}

	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read login_pass.html: %w", err)
	}

	resp, err = c.client.Get(c.hostpath + "/login_counter.html?stap=0")
	if err != nil {
		return fmt.Errorf("unable to get login_counter.html: %w", err)
	}

	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read login_counter.html: %w", err)
	}

	resp, err = c.client.Get(c.hostpath + "/login_counter.html?stap=1")
	if err != nil {
		return fmt.Errorf("unable to get login_counter.html: %w", err)
	}

	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read login_counter.html: %w", err)
	}

	for {
		resp, err = c.client.Get(c.hostpath + "/login_counter.html?stap=2")
		if err != nil {
			return fmt.Errorf("unable to get login_counter.html: %w", err)
		}

		defer resp.Body.Close()

		_, err = io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("unable to read login_counter.html: %w", err)
		}

		if resp.Header.Get("auth_state") == "1" {
			break
		}
	}

	resp, err = c.client.Get(c.hostpath + "/login.cgi?action=LOGIN")
	if err != nil {
		return fmt.Errorf("unable to get login.cgi: %w", err)
	}

	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read login.cgi: %w", err)
	}

	if resp.Header.Get("Location") == c.hostpath+"/error.html" {
		return fmt.Errorf("login Failed")
	}

	log.Printf("Login Successful to %s", c.hostpath)

	c._loggedIn = true

	return nil
}

func (c *CP) logout() error {
	if !c._loggedIn {
		return nil
	}

	resp, err := c.client.Get(c.hostpath + "/logout.html")
	if err != nil {
		return fmt.Errorf("unable to get logout.html: %w", err)
	}

	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read logout.html: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to logout")
	}

	log.Printf("Logout from %s", c.hostpath)

	c._loggedIn = false

	return nil
}
