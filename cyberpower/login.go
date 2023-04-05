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

	resp, err = c.client.PostForm(c.hostpath+"/login_pass.cgi", c.loginForm)
	if err != nil {
		log.Println(err)
		return false
	}
	defer resp.Body.Close()
	io.ReadAll(resp.Body)

	// resp, err = c.client.Get(c.hostpath + "/login_pass.html?action=LOGIN&username=" + c.loginForm.Get("username") + "&password=" + c.loginForm.Get("password"))
	resp, err = c.client.Get(c.hostpath + "/login_pass.html")
	if err != nil {
		log.Println(err)
		return false
	}
	defer resp.Body.Close()
	io.ReadAll(resp.Body)

	resp, err = c.client.Get(c.hostpath + "/login_counter.html?stap=0")
	if err != nil {
		log.Println(err)
		return false
	}
	defer resp.Body.Close()
	io.ReadAll(resp.Body)

	resp, err = c.client.Get(c.hostpath + "/login_counter.html?stap=1")
	if err != nil {
		log.Println(err)
		return false
	}
	defer resp.Body.Close()
	io.ReadAll(resp.Body)

	for {
		resp, err = c.client.Get(c.hostpath + "/login_counter.html?stap=2")
		if err != nil {
			log.Println(err)
			return false
		}
		defer resp.Body.Close()
		io.ReadAll(resp.Body)

		if resp.Header.Get("auth_state") == "1" {
			break
		}
	}

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
