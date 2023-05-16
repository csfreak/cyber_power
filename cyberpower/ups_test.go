package cyberpower

import (
	"bytes"
	"reflect"
	"testing"
	"text/template"

	"golang.org/x/net/html"
)

func TestUPS_update(t *testing.T) {
	updateTemplate, err := template.New("UpsStatus").
		Funcs(map[string]any{"secToMin": func(s any) int { return s.(int) / 60 }}).
		Parse(htmlUpsStatusBodyTemplate)
	if err != nil {
		t.Errorf("unable to parse template: %v", err)
		return
	}

	tests := []struct {
		name      string
		updateErr bool
		want      *UPS
	}{
		{
			name: "success",
			want: &UPS{
				Input: InputPower{
					Status:    "Normal",
					Voltage:   120.0,
					Frequency: 60.0,
				},
				Output: OutputPower{
					Status:      "Normal",
					Voltage:     120.0,
					Frequency:   60.0,
					Current:     6.0,
					LoadWatts:   720,
					LoadPercent: 70,
				},
				Battery: BatteryPower{
					Status:            "Normal",
					RemainingCapacity: 100,
					RemainingRuntime:  1080,
				},
				TempC:  25.0,
				TempF:  77.0,
				Status: "Normal",
			},
			updateErr: false,
		},
		{
			name:      "fail_get",
			want:      &UPS{},
			updateErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var parsedHTML bytes.Buffer
			err := updateTemplate.Execute(&parsedHTML, tt.want)
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
			u := &UPS{parent: c}
			c.ups = u

			if err := u.update(); (err != nil) != tt.updateErr {
				t.Errorf("UPS.update() error = %v, wantErr %v", err, tt.updateErr)
			}
			if !c.calledOnce() {
				t.Errorf("expeced CyberPower methods to be called once, found %d", c.calls())
			}
			tt.want.parent = c
			if !reflect.DeepEqual(tt.want, u) {
				t.Errorf("expected %v, got %v", tt.want, u)
			}
		})
	}
}

func TestUPS_getParent(t *testing.T) {
	c := &mockCP{}

	tests := []struct {
		name string
		u    UPS
		want CyberPower
	}{
		{"success", UPS{parent: c}, c},
		{"no_parent", UPS{}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.getParent(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UPS.getParent() = %v, want %v", got, tt.want)
			}
		})
	}
}
