package cyberpower

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"

	_ "github.com/golang/mock/mockgen/model"
	"github.com/h2non/gock"
	"golang.org/x/net/html"
)

func TestNewCyberPower(t *testing.T) {
	type args struct {
		host     string
		username string
		password string
		validate bool
	}

	defaultArgs := args{testhostname, testuser, testpw, false}

	tests := []struct {
		name  string
		args  args
		check func(got *CP, err error)
	}{
		{
			"hostpath set",
			defaultArgs,
			func(got *CP, err error) {
				if got.hostpath != testhostpath {
					t.Errorf("expected hostpath %s found %s", testhostpath, got.hostpath)
				}
			},
		},
		{
			"loginform set",
			defaultArgs,
			func(got *CP, err error) {
				if !reflect.DeepEqual(got.loginForm, testloginform) {
					t.Errorf("expected loginform %v found %v", testloginform, got.loginForm)
				}
			},
		},
		{
			"env set",
			defaultArgs,
			func(got *CP, err error) {
				if got.env.getParent() != got {
					t.Errorf("expected env.parent to be set, found %v", got.env)
				}
			},
		},
		{
			"ups set",
			defaultArgs,
			func(got *CP, err error) {
				if got.ups.getParent() != got {
					t.Errorf("expected env.parent to be set, found %v", got.ups)
				}
			},
		},
		{
			"loggedIn false",
			defaultArgs,
			func(got *CP, err error) {
				if got._loggedIn {
					t.Errorf("expected loggedIn to be false, found %v", got._loggedIn)
				}
			},
		},
		{
			"validate failed",
			args{testhostname, testuser, testpw, true},
			func(got *CP, err error) {
				if err == nil {
					t.Error("expected login validation failure, found none")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCyberPower(tt.args.host, tt.args.username, tt.args.password, tt.args.validate)
			tt.check(got.(*CP), err)
		})
	}
}

func TestCP_loggedIn(t *testing.T) {
	type fields struct {
		_loggedIn bool
	}

	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"logged in", fields{true}, true},
		{"not loggedin", fields{false}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CP{
				_loggedIn: tt.fields._loggedIn,
			}
			if got := c.loggedIn(); got != tt.want {
				t.Errorf("CP.loggedIn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCP_getHost(t *testing.T) {
	t.Run(testhostname, func(t *testing.T) {
		c := &CP{
			hostpath: testhostpath,
		}
		if got := c.getHost(); got != testhostname {
			t.Errorf("CP.getHost() = %v, want %v", got, testhostname)
		}
	})
}

func TestCP_getEnv(t *testing.T) {
	testenv := &ENV{}
	tests := []struct {
		name  string
		env   CPModule
		check func(got CPModule, ok bool)
	}{
		{"env type", testenv, func(got CPModule, ok bool) {
			if !ok {
				t.Errorf("failed get env")
			}
			if got != testenv {
				t.Errorf("expected %v, got %v", testenv, got)
			}
		}},
		{"ups type", &UPS{}, func(got CPModule, ok bool) {
			if ok {
				t.Errorf("expected error for type assertion, found none")
			}
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CP{
				env: tt.env,
			}
			tt.check(c.getEnv())
		})
	}
}

func TestCP_getUps(t *testing.T) {
	testups := &UPS{}
	tests := []struct {
		name  string
		ups   CPModule
		check func(got CPModule, ok bool)
	}{
		{"ups type", testups, func(got CPModule, ok bool) {
			if !ok {
				t.Errorf("failed get ups")
			}
			if got != testups {
				t.Errorf("expected %v, got %v", testups, got)
			}
		}},
		{"env type", &ENV{}, func(got CPModule, ok bool) {
			if ok {
				t.Errorf("expected error for type assertion, found none")
			}
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CP{
				ups: tt.ups,
			}
			tt.check(c.getUps())
		})
	}
}

func TestCP_update(t *testing.T) {
	type fields struct {
		ups *mockCpModule
		env *mockCpModule
	}

	tests := []struct {
		name  string
		f     fields
		check func(fields, error)
	}{
		{
			"update success",
			fields{
				ups: &mockCpModule{updateErr: nil},
				env: &mockCpModule{updateErr: nil},
			},
			func(f fields, err error) {
				if err != nil {
					t.Errorf("expected no error, found %v", err)
				}
				if !f.ups.calledOnce() {
					t.Errorf("expected ups update called once, found %d", f.ups.calls())
				}
				if !f.env.calledOnce() {
					t.Errorf("expected env update called once, found %d", f.env.calls())
				}
			},
		},
		{
			"update ups failed",
			fields{
				ups: &mockCpModule{updateErr: fmt.Errorf("ups update failed")},
				env: &mockCpModule{updateErr: nil},
			},
			func(f fields, err error) {
				if err == nil {
					t.Error("expected error, none found")
				}
				if !f.ups.calledOnce() {
					t.Errorf("expected ups update called once, found %d", f.ups.calls())
				}
			},
		},
		{
			"update env failed",
			fields{
				ups: &mockCpModule{updateErr: nil},
				env: &mockCpModule{updateErr: fmt.Errorf("env update failed")},
			},
			func(f fields, err error) {
				if err == nil {
					t.Error("expected error, none found")
				}
				if !f.env.calledOnce() {
					t.Errorf("expected env update called once, found %d", f.env.calls())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CP{
				ups:       tt.f.ups,
				env:       tt.f.env,
				_loggedIn: false,
			}
			tt.check(tt.f, c.update())
		})
	}
}

func TestCP_get(t *testing.T) {
	c, _ := NewCyberPower(testhostname, testuser, testpw, false)
	cp := c.(*CP)

	gock.InterceptClient(&cp.client)
	defer gock.RestoreClient(&cp.client)
	defer gock.Off()

	tests := []struct {
		name     string
		path     string
		loggedIn bool
		body     *Body
		respCode int
		respErr  error
		wantErr  bool
	}{
		{
			name:     "get empty body",
			path:     testpath,
			loggedIn: true,
			body:     &Body{body: ""},
			respCode: http.StatusOK,
			respErr:  nil,
			wantErr:  false,
		},
		{
			name:     "get full body",
			path:     testpath,
			loggedIn: true,
			body:     &Body{body: ""},
			respCode: http.StatusOK,
			respErr:  nil,
			wantErr:  false,
		},
		{
			name:     "get 404",
			path:     testpath,
			loggedIn: true,
			body:     &Body{body: ""},
			respCode: http.StatusNotFound,
			respErr:  nil,
			wantErr:  false,
		},
		{
			name:     "get login failed",
			path:     testpath,
			loggedIn: false,
			body:     nil,
			respCode: http.StatusOK,
			respErr:  nil,
			wantErr:  true,
		},
		{
			name:     "get dial error",
			path:     testpath,
			loggedIn: true,
			body:     nil,
			respCode: 0,
			respErr:  fmt.Errorf("dial error"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer gock.Flush()
			defer gock.CleanUnmatchedRequest()

			var want *html.Node

			if tt.body != nil {
				want, _ = html.Parse(tt.body)
				tt.body.Reset()
			}

			req := gock.New(testhostpath).Get(tt.path)
			if tt.respErr != nil {
				req.ReplyError(tt.respErr)
			} else if tt.body != nil {
				req.Reply(tt.respCode).Body(tt.body)
			}
			cp._loggedIn = tt.loggedIn
			got, err := c.get(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("got error %v, expected %v", err, tt.wantErr)
			}
			if tt.respErr != nil && !errors.Is(err, tt.respErr) {
				t.Errorf("expected response error %v, got %v", tt.respErr, err)
			}

			if want != nil {
				if !reflect.DeepEqual(got, want) {
					t.Errorf("expected CyberPower.get() = %v, got  %v", want, got)
				}
				if req.Counter != 0 {
					t.Errorf("missing calls: %d", req.Counter)
				}
				if gock.HasUnmatchedRequest() {
					for _, r := range gock.GetUnmatchedRequests() {
						t.Errorf("unexpected request: %v", r)
					}
				}
			}
		})
	}
}

func TestFromENV(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		validate bool
		check    func(got CyberPower, err error)
	}{
		{
			"TestFromEnv success",
			map[string]string{
				"CYBERPOWER_HOST":     testhostname,
				"CYBERPOWER_USERNAME": testuser,
				"CYBERPOWER_PASSWORD": testpw,
			},
			false,
			func(got CyberPower, err error) {
				if err != nil {
					t.Errorf("expected no error; found %v", err)
				}
				gotcp, ok := got.(*CP)
				if !ok {
					t.Errorf("expected *CP, got %v", reflect.TypeOf(got))
				}
				if gotcp.hostpath != testhostpath {
					t.Errorf("epected hostpath %s; found %s", testhostpath, gotcp.hostpath)
				}
				if gotcp.loginForm.Get("username") != testuser {
					t.Errorf("epected username %s; found %s", testuser, gotcp.loginForm.Get("username"))
				}
				if gotcp.loginForm.Get("password") != testpw {
					t.Errorf("epected password %s; found %s", testpw, gotcp.loginForm.Get("password"))
				}
			},
		},
		{
			"TestFromEnv no host",
			map[string]string{
				"CYBERPOWER_HOST":     "UNSET",
				"CYBERPOWER_USERNAME": testuser,
				"CYBERPOWER_PASSWORD": testpw,
			},
			false,
			func(got CyberPower, err error) {
				if err == nil || !strings.Contains(err.Error(), "CYBERPOWER_HOST") {
					t.Errorf("expected error %v, found %v", fmt.Errorf("unable to load CYBERPOWER_HOST"), err)
				}
				if got != nil {
					t.Errorf("expected nil, got %v", got)
				}
			},
		},
		{
			"TestFromEnv no user",
			map[string]string{
				"CYBERPOWER_HOST":     testhostname,
				"CYBERPOWER_USERNAME": "UNSET",
				"CYBERPOWER_PASSWORD": testpw,
			},
			false,
			func(got CyberPower, err error) {
				if err == nil || !strings.Contains(err.Error(), "CYBERPOWER_USERNAME") {
					t.Errorf("expected error %v, found %v", fmt.Errorf("unable to load CYBERPOWER_USERNAME"), err)
				}
				if got != nil {
					t.Errorf("expected nil, got %v", got)
				}
			},
		},
		{
			"TestFromEnv no password",
			map[string]string{
				"CYBERPOWER_HOST":     testhostname,
				"CYBERPOWER_USERNAME": testuser,
				"CYBERPOWER_PASSWORD": "UNSET",
			},
			false,
			func(got CyberPower, err error) {
				if err == nil || !strings.Contains(err.Error(), "CYBERPOWER_PASSWORD") {
					t.Errorf("expected error %v, found %v", fmt.Errorf("unable to load CYBERPOWER_PASSWORD"), err)
				}
				if got != nil {
					t.Errorf("expected nil, got %v", got)
				}
			},
		},
		{
			"TestFromEnv failed validation",
			map[string]string{
				"CYBERPOWER_HOST":     testhostname,
				"CYBERPOWER_USERNAME": testuser,
				"CYBERPOWER_PASSWORD": testpw,
			},
			true,
			func(got CyberPower, err error) {
				if err == nil || !strings.Contains(err.Error(), "login") {
					t.Errorf("expected error %v, found %v", fmt.Errorf("unable to login"), err)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				if v != "UNSET" {
					t.Setenv(k, v)
				} else {
					os.Unsetenv(k)
				}
			}
			got, err := FromENV(tt.validate)
			tt.check(got, err)
		})
	}
}
