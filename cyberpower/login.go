package cyberpower

import (
	"io"
	"log"
)

func (c *CyberPower) login() bool {
	resp, err := c.client.Get(c.hostpath + "/login.html")
	if err != nil {
		log.Println(err)
		return false
	}
	defer resp.Body.Close()
	io.ReadAll(resp.Body)

	resp, err = c.client.Get(c.hostpath + "/login_pass.html?action=LOGIN&username=" + c.loginForm.Get("username") + "&password=" + c.loginForm.Get("password"))
	if err != nil {
		log.Println(err)
		return false
	}
	defer resp.Body.Close()
	io.ReadAll(resp.Body)

	resp, err = c.client.Get(c.hostpath + "/login.cgi?action=LOGIN")
	if err != nil {
		log.Println(err)
		return false
	}
	defer resp.Body.Close()
	io.ReadAll(resp.Body)
	if resp.Header.Get("Location") == c.hostpath+"/error.html" {
		log.Println("Login Failed")
		return false
	}
	log.Printf("Login Successful to %s", c.hostpath)
	c.logged_in = true
	return true
}

func (c *CyberPower) Logout() {
	c.client.Get(c.hostpath + "/logout.html")
	log.Printf("Logout from %s", c.hostpath)
	c.logged_in = false
}
