package cyberpower

import (
	"bytes"
	"reflect"
	"testing"
	"text/template"

	"golang.org/x/net/html"
)

func TestENV_update(t *testing.T) {
	updateTemplate, err := template.New("EnvStatus").
		Parse(htmlEnvStatusBodyTemplate)
	if err != nil {
		t.Errorf("unable to parse template: %v", err)
		return
	}

	tests := []struct {
		name      string
		updateErr bool
		parseErr  bool
		want      *ENV
	}{
		{
			name: "success",
			want: &ENV{
				Name:     "TestSensor",
				Location: "TestLocation",
				TempF:    77.0,
				Humidity: 45,
			},
			updateErr: false,
			parseErr:  false,
		},
		{
			name:      "fail_get",
			want:      &ENV{},
			updateErr: true,
			parseErr:  false,
		},
		{
			name: "fail_parse",
			want: &ENV{
				Name:     "TestSensor",
				Location: "TestLocation",
				TempF:    77.0,
				Humidity: 45,
			},
			updateErr: false,
			parseErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var parsedHTML bytes.Buffer
			var obj any
			if tt.parseErr {
				obj = struct {
					Name     string
					Location string
					TempF    string
					Humidity string
				}{"TestSensor", "TestLocation", "notanumber", "notanumber"}
			} else {
				obj = tt.want
			}
			err := updateTemplate.Execute(&parsedHTML, obj)
			if err != nil {
				t.Errorf("unable to execute tempalte: %v", err)
				return
			}
			parsed, err := html.Parse(&parsedHTML)
			if err != nil {
				t.Errorf("unable to parse html: %v", err)
				return
			}
			c := &mockCP{updateData: parsed, updateErr: tt.updateErr}
			e := &ENV{parent: c}
			c.env = e

			if err := e.update(); (err != nil) != tt.updateErr {
				t.Errorf("ENV.update() error = %v, wantErr %v", err, tt.updateErr)
			}
			if !c.calledOnce() {
				t.Errorf("expeced CyberPower methods to be called once, found %d", c.calls())
			}
			tt.want.parent = c
			if reflect.DeepEqual(tt.want, e) == tt.parseErr {
				t.Errorf("expected %v, got %v", tt.want, e)
			}
		})
	}
}

func TestENV_getParent(t *testing.T) {
	c := &mockCP{}

	tests := []struct {
		name string
		e    ENV
		want CyberPower
	}{
		{"success", ENV{parent: c}, c},
		{"no_parent", ENV{}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.getParent(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ENV.getParent() = %v, want %v", got, tt.want)
			}
		})
	}
}
