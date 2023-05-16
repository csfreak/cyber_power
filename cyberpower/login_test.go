package cyberpower

import (
	"testing"

	"github.com/h2non/gock"
)

func TestCP_login(t *testing.T) {
	c, _ := NewCyberPower(testhostname, testuser, testpw, false)
	cp := c.(*CP)

	gock.Intercept()
	defer gock.Off()

	tests := []struct {
		name    string
		mocks   []*gock.Mocker
		wantErr bool
	}{
		{
			name:    "success - logged in",
			mocks:   []*gock.Mocker{mckrLoginHTMLOk(), mckrLoginPassCgiOk(), mckrLoginPassHTMLOk(), mckrLoginCounterS0Ok(), mckrLoginCounterS1Ok(), mckrLoginCounterS2Auth1(), mckrLoginCgiRedirSuccess()},
			wantErr: false,
		},
		{
			name:    "success - logged in after auth-state 0",
			mocks:   []*gock.Mocker{mckrLoginHTMLOk(), mckrLoginPassCgiOk(), mckrLoginPassHTMLOk(), mckrLoginCounterS0Ok(), mckrLoginCounterS1Ok(), mckrLoginCounterS2Auth0(), mckrLoginCounterS2Auth1(), mckrLoginCgiRedirSuccess()},
			wantErr: false,
		},
		{
			name:    "fail - login.html",
			mocks:   []*gock.Mocker{mckrLoginHTMLErr()},
			wantErr: true,
		},
		{
			name:    "fail - login_pass.cgi",
			mocks:   []*gock.Mocker{mckrLoginHTMLOk(), mckrLoginPassCgiErr()},
			wantErr: true,
		},
		{
			name:    "fail - login_pass.html",
			mocks:   []*gock.Mocker{mckrLoginHTMLOk(), mckrLoginPassCgiOk(), mckrLoginPassHTMLErr()},
			wantErr: true,
		},
		{
			name:    "fail - login_counter.html?stap=0",
			mocks:   []*gock.Mocker{mckrLoginHTMLOk(), mckrLoginPassCgiOk(), mckrLoginPassHTMLOk(), mckrLoginCounterS0Err()},
			wantErr: true,
		},
		{
			name:    "fail - login_counter.html?stap=1",
			mocks:   []*gock.Mocker{mckrLoginHTMLOk(), mckrLoginPassCgiOk(), mckrLoginPassHTMLOk(), mckrLoginCounterS0Ok(), mckrLoginCounterS1Err()},
			wantErr: true,
		},
		{
			name:    "fail - login_counter.html?stap=2",
			mocks:   []*gock.Mocker{mckrLoginHTMLOk(), mckrLoginPassCgiOk(), mckrLoginPassHTMLOk(), mckrLoginCounterS0Ok(), mckrLoginCounterS1Ok(), mckrLoginCounterS2Err()},
			wantErr: true,
		},
		{
			name:    "fail - login.cgi",
			mocks:   []*gock.Mocker{mckrLoginHTMLOk(), mckrLoginPassCgiOk(), mckrLoginPassHTMLOk(), mckrLoginCounterS0Ok(), mckrLoginCounterS1Ok(), mckrLoginCounterS2Auth1(), mckrLoginCgiErr()},
			wantErr: true,
		},
		{
			name:    "fail - redir error",
			mocks:   []*gock.Mocker{mckrLoginHTMLOk(), mckrLoginPassCgiOk(), mckrLoginPassHTMLOk(), mckrLoginCounterS0Ok(), mckrLoginCounterS1Ok(), mckrLoginCounterS2Auth1(), mckrLoginCgiRedirFail()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(gock.Flush)
			t.Cleanup(gock.CleanUnmatchedRequest)

			for _, mocker := range tt.mocks {
				gock.Register(mocker)
			}
			cp._loggedIn = false
			err := cp.login()
			if (err != nil) != tt.wantErr {
				t.Errorf("unexpected result; wantErr = %v, found %v", tt.wantErr, err)
			}
			for _, mocker := range tt.mocks {
				if mocker.Request().Counter != 0 {
					t.Errorf("missing calls to %v", mocker.Request())
				}
			}
			if gock.HasUnmatchedRequest() {
				for _, request := range gock.GetUnmatchedRequests() {
					t.Errorf("unexpected call to %v", request)
				}
			}
		})
	}
}

func TestCP_logout(t *testing.T) {
	c, _ := NewCyberPower(testhostname, testuser, testpw, false)
	cp := c.(*CP)

	gock.Intercept()
	defer gock.Off()

	tests := []struct {
		name     string
		loggedIn bool
		mocks    []*gock.Mocker
		wantErr  bool
	}{
		{
			name:     "success - logged in",
			loggedIn: true,
			mocks:    []*gock.Mocker{mckrLogoutHTMLOk()},
			wantErr:  false,
		}, {
			name:     "success - logged out",
			loggedIn: false,
			mocks:    []*gock.Mocker{},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(gock.Flush)
			t.Cleanup(gock.CleanUnmatchedRequest)

			for _, mocker := range tt.mocks {
				gock.Register(mocker)
			}

			cp._loggedIn = tt.loggedIn

			err := cp.logout()
			if (err != nil) != tt.wantErr {
				t.Errorf("expected wantErr = %v, found %v", tt.wantErr, err)
			}

			if cp._loggedIn {
				t.Errorf("expected CP._loggedIn == false, found %v", cp._loggedIn)
			}

			for _, mocker := range tt.mocks {
				if mocker.Request().Counter != 0 {
					t.Errorf("missing calls to %v", mocker.Request())
				}
			}

			if gock.HasUnmatchedRequest() {
				for _, request := range gock.GetUnmatchedRequests() {
					t.Errorf("unexpected call to %v", request)
				}
			}
		})
	}
}
