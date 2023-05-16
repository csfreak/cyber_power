package cyberpower

import (
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/net/html"

	"github.com/h2non/gock"
)

type mock struct {
	_calls int
}

func (m *mock) calls() int {
	return m._calls
}

func (m *mock) calledOnce() bool {
	return m._calls == 1
}

type mockCpModule struct {
	mock
	parent    CyberPower
	updateErr error
}

func (m *mockCpModule) update() error {
	m._calls++
	return m.updateErr
}

func (m *mockCpModule) getParent() CyberPower {
	m._calls++
	return m.parent
}

type mockCP struct {
	mock
	ups        CPModule
	env        CPModule
	updateData *html.Node
	updateErr  bool
}

func (m *mockCP) login() error  { m._calls++; return errors.New("not implemented") }
func (m *mockCP) logout() error { m._calls++; return nil }
func (m *mockCP) get(_ string) (*html.Node, error) {
	m._calls++
	if m.updateErr {
		return nil, errors.New("update error")
	}

	return m.updateData, nil
}
func (m *mockCP) getHost() string { return "" }
func (m *mockCP) getEnv() (*ENV, bool) {
	e, ok := m.env.(*ENV)
	return e, ok
}

func (m *mockCP) getUps() (*UPS, bool) {
	u, ok := m.ups.(*UPS)
	return u, ok
}
func (m *mockCP) loggedIn() bool { return false }
func (m *mockCP) update() error {
	m._calls++

	err := m.ups.update()
	if err != nil {
		return fmt.Errorf("unable to update ups: %v", err)
	}

	err = m.env.update()
	if err != nil {
		return fmt.Errorf("unable to update env: %v", err)
	}

	err = m.logout()
	if err != nil {
		return fmt.Errorf("unable to logout: %v", err)
	}

	return nil
}

// HTTP Mocks

var (
	reqLoginHTML    = func() *gock.Request { return gock.NewRequest().URL(testhostpath).Get("/login.html") }
	reqLoginPassCgi = func() *gock.Request {
		return gock.NewRequest().URL(testhostpath).Post("/login_pass.cgi").MatchType("url").BodyString(testloginform.Encode())
	}
	reqLoginPassHTML  = func() *gock.Request { return gock.NewRequest().URL(testhostpath).Get("/login_pass.html") }
	reqLoginCounterS0 = func() *gock.Request {
		return gock.NewRequest().URL(testhostpath).Get("/login_counter.html").MatchParam("stap", "0")
	}
	reqLoginCounterS1 = func() *gock.Request {
		return gock.NewRequest().URL(testhostpath).Get("/login_counter.html").MatchParam("stap", "1")
	}
	reqLoginCounterS2 = func() *gock.Request {
		return gock.NewRequest().URL(testhostpath).Get("/login_counter.html").MatchParam("stap", "2")
	}
	reqLoginCgi = func() *gock.Request {
		return gock.NewRequest().URL(testhostpath).Get("/login.cgi").MatchParam("action", LoginAction)
	}
	reqLogoutHTML       = func() *gock.Request { return gock.NewRequest().URL(testhostpath).Get("/logout.html") }
	respGenericOk       = func() *gock.Response { return gock.NewResponse().Status(http.StatusOK).BodyString("") }
	respGenericErr      = func() *gock.Response { return gock.NewResponse().SetError(fmt.Errorf("generic error")) }
	respLoginAuthState0 = func() *gock.Response {
		return gock.NewResponse().Status(http.StatusOK).AddHeader("auth_state", "0").BodyString("")
	}
	respLoginAuthState1 = func() *gock.Response {
		return gock.NewResponse().Status(http.StatusOK).AddHeader("auth_state", "1").BodyString("")
	}
	respLoginRedirSuccess = func() *gock.Response {
		return gock.NewResponse().Status(http.StatusFound).AddHeader("Location", fmt.Sprintf("%s/status.html", testhostpath)).BodyString("")
	}
	respLoginRedirFail = func() *gock.Response {
		return gock.NewResponse().Status(http.StatusFound).AddHeader("Location", fmt.Sprintf("%s/error.html", testhostpath)).BodyString("")
	}
	mckrLoginHTMLOk          = func() *gock.Mocker { return gock.NewMock(reqLoginHTML(), respGenericOk()) }
	mckrLoginHTMLErr         = func() *gock.Mocker { return gock.NewMock(reqLoginHTML(), respGenericErr()) }
	mckrLoginPassCgiOk       = func() *gock.Mocker { return gock.NewMock(reqLoginPassCgi(), respGenericOk()) }
	mckrLoginPassCgiErr      = func() *gock.Mocker { return gock.NewMock(reqLoginPassCgi(), respGenericErr()) }
	mckrLoginPassHTMLOk      = func() *gock.Mocker { return gock.NewMock(reqLoginPassHTML(), respGenericOk()) }
	mckrLoginPassHTMLErr     = func() *gock.Mocker { return gock.NewMock(reqLoginPassHTML(), respGenericErr()) }
	mckrLoginCounterS0Ok     = func() *gock.Mocker { return gock.NewMock(reqLoginCounterS0(), respGenericOk()) }
	mckrLoginCounterS0Err    = func() *gock.Mocker { return gock.NewMock(reqLoginCounterS0(), respGenericErr()) }
	mckrLoginCounterS1Ok     = func() *gock.Mocker { return gock.NewMock(reqLoginCounterS1(), respGenericOk()) }
	mckrLoginCounterS1Err    = func() *gock.Mocker { return gock.NewMock(reqLoginCounterS1(), respGenericErr()) }
	mckrLoginCounterS2Auth0  = func() *gock.Mocker { return gock.NewMock(reqLoginCounterS2(), respLoginAuthState0()) }
	mckrLoginCounterS2Auth1  = func() *gock.Mocker { return gock.NewMock(reqLoginCounterS2(), respLoginAuthState1()) }
	mckrLoginCounterS2Err    = func() *gock.Mocker { return gock.NewMock(reqLoginCounterS2(), respGenericErr()) }
	mckrLoginCgiRedirSuccess = func() *gock.Mocker { return gock.NewMock(reqLoginCgi(), respLoginRedirSuccess()) }
	mckrLoginCgiRedirFail    = func() *gock.Mocker { return gock.NewMock(reqLoginCgi(), respLoginRedirFail()) }
	mckrLoginCgiErr          = func() *gock.Mocker { return gock.NewMock(reqLoginCgi(), respGenericErr()) }
	mckrLogoutHTMLOk         = func() *gock.Mocker { return gock.NewMock(reqLogoutHTML(), respGenericOk()) }
)
