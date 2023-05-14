package cyberpower

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"

	_ "github.com/golang/mock/mockgen/model"
	"github.com/h2non/gock"
	"golang.org/x/net/html"
)

var (
	testhostname  string     = "testhost"
	testuser      string     = "testuser"
	testpw        string     = "testpw"
	testhostpath  string     = "http://" + testhostname
	testpath      string     = "/test/path"
	testloginform url.Values = url.Values{
		"action":   []string{"LOGIN"},
		"username": []string{testuser},
		"password": []string{testpw},
	}
)

func TestNewCyberPower(t *testing.T) {
	type args struct {
		host     string
		username string
		password string
		validate bool
	}
	default_args := args{testhostname, testuser, testpw, false}

	tests := []struct {
		name  string
		args  args
		check func(got *CP, err error)
	}{
		{
			"hostpath set",
			default_args,
			func(got *CP, err error) {
				if got.hostpath != testhostpath {
					t.Errorf("expected hostpath %s found %s", testhostpath, got.hostpath)
				}
			},
		},
		{
			"loginform set",
			default_args,
			func(got *CP, err error) {
				if !reflect.DeepEqual(got.loginForm, testloginform) {
					t.Errorf("expected loginform %v found %v", testloginform, got.loginForm)
				}
			},
		},
		{
			"env set",
			default_args,
			func(got *CP, err error) {
				if got.env.getParent() != got {
					t.Errorf("expected env.parent to be set, found %v", got.env)
				}
			},
		},
		{
			"ups set",
			default_args,
			func(got *CP, err error) {
				if got.ups.getParent() != got {
					t.Errorf("expected env.parent to be set, found %v", got.ups)
				}
			},
		},
		{
			"logged_in false",
			default_args,
			func(got *CP, err error) {
				if got._logged_in {
					t.Errorf("expected logged_in to be false, found %v", got._logged_in)
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

func TestCP_logged_in(t *testing.T) {
	type fields struct {
		_logged_in bool
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
				_logged_in: tt.fields._logged_in,
			}
			if got := c.logged_in(); got != tt.want {
				t.Errorf("CP.logged_in() = %v, want %v", got, tt.want)
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
		env   CyberPowerModule
		check func(got CyberPowerModule, ok bool)
	}{
		{"env type", testenv, func(got CyberPowerModule, ok bool) {
			if !ok {
				t.Errorf("failed get env")
			}
			if got != testenv {
				t.Errorf("expected %v, got %v", testenv, got)
			}
		}},
		{"ups type", &UPS{}, func(got CyberPowerModule, ok bool) {
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
		ups   CyberPowerModule
		check func(got CyberPowerModule, ok bool)
	}{
		{"ups type", testups, func(got CyberPowerModule, ok bool) {
			if !ok {
				t.Errorf("failed get ups")
			}
			if got != testups {
				t.Errorf("expected %v, got %v", testups, got)
			}
		}},
		{"env type", &ENV{}, func(got CyberPowerModule, ok bool) {
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
		ups *mock_cpmodule
		env *mock_cpmodule
	}
	tests := []struct {
		name  string
		f     fields
		check func(fields, error)
	}{
		{
			"update success",
			fields{
				ups: &mock_cpmodule{update_error: nil},
				env: &mock_cpmodule{update_error: nil}},
			func(f fields, err error) {
				if err != nil {
					t.Errorf("expected no error, found %v", err)
				}
				if !f.ups.called_once() {
					t.Errorf("expected ups update called once, found %d", f.ups.calls())
				}
				if !f.env.called_once() {
					t.Errorf("expected env update called once, found %d", f.env.calls())
				}
			},
		},
		{
			"update ups failed",
			fields{
				ups: &mock_cpmodule{update_error: fmt.Errorf("ups update failed")},
				env: &mock_cpmodule{update_error: nil}},
			func(f fields, err error) {
				if err == nil {
					t.Error("expected error, none found")
				}
				if !f.ups.called_once() {
					t.Errorf("expected ups update called once, found %d", f.ups.calls())
				}
			},
		},
		{
			"update env failed",
			fields{
				ups: &mock_cpmodule{update_error: nil},
				env: &mock_cpmodule{update_error: fmt.Errorf("env update failed")}},
			func(f fields, err error) {
				if err == nil {
					t.Error("expected error, none found")
				}
				if !f.env.called_once() {
					t.Errorf("expected env update called once, found %d", f.env.calls())
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CP{
				ups:        tt.f.ups,
				env:        tt.f.env,
				_logged_in: false,
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
		name      string
		path      string
		logged_in bool
		body      *Body
		resp_code int
		resp_err  error
		wantErr   bool
	}{
		{
			"get empty body",
			testpath,
			true,
			&Body{body: ""},
			200,
			nil,
			false,
		},
		{
			"get full body",
			testpath,
			true,
			&Body{body: htmlUpsStatusBody},
			200,
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want, parseErr := html.Parse(tt.body)
			tt.body.Reset()
			req := gock.New(testhostpath).Get(tt.path)
			if tt.resp_err != nil {
				req.ReplyError(tt.resp_err)
			} else {
				req.Reply(tt.resp_code).Body(tt.body)
			}
			cp._logged_in = tt.logged_in
			got, err := c.get(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("got error %v, expected %v", err, tt.wantErr)
			}
			if err != tt.resp_err {
				t.Errorf("expected response error %v, got %v", tt.resp_err, err)
			}
			if err != parseErr {
				t.Errorf("expected parse errer %v, got %v", parseErr, err)
			}
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
		})
	}
}
