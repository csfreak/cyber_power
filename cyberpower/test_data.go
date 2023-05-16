package cyberpower

import (
	"fmt"
	"io"
	"net/url"
)

var (
	testhostname  = "testhost"
	testuser      = "testuser"
	testpw        = "testpw"
	testhostpath  = "http://" + testhostname
	testpath      = "/test/path"
	testloginform = url.Values{
		"action":   []string{"LOGIN"},
		"username": []string{testuser},
		"password": []string{testpw},
	}
	htmlDefaultBase = `
	<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
	<html x	mlns="http://www.w3.org/1999/xhtml">
	<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
	<link rel="shortcut icon" href="icon/icon.ico" />
	<link href="css/rmc.css" rel="stylesheet" type="text/css" />
	<title> TEST</title>
	</head>
	<body>
	%s
	</body></html>
	`
	htmlUpsStatusBodyTemplate = `
	<span class="caption">Input</span><br/>
	<div class="gap">
		<span class="lb statusLb">Status</span><span class="txt">{{.Input.Status}}</span><br/>
		<span class="hide"><span class="firstItem">Phase1</span><span class="anotherItem">Phase2</span><span class="anotherItem">Phase3</span></br></span>
		<span class="lb statusLb">Voltage</span><span class="firstData">{{.Input.Voltage}} V</span><span class="hide"><span class="anotherData"> V</span><span class="anotherData"> V</span></span><br/>
		<span class=""><span class="lb statusLb">Frequency</span><span class="firstData">{{.Input.Frequency}} Hz</span><span class="hide"><span class="anotherData"> Hz</span><span class="anotherData"> Hz</span></span></br></span>
		<span class="hide"><span class="lb statusLb">Current</span><span class="firstData"> A</span><span class="anotherData"> A</span><span class="anotherData"> A</span></br></span>
		<span class="hide"><span class="lb statusLb">Power Factor</span><span class="firstData"></span><span class="anotherData"></span><span class="anotherData"></span></br></span>
	</div>
	<span class="caption hide">Bypass</span>
	<div class="gap hide">
		<span class="lb statusLb">Status</span><span class="txt"></span><br>
		<span class="firstItem">Phase1</span><span class="anotherItem">Phase2</span><span class="anotherItem">Phase3</span><br>
		<span class="hide"><span class="lb statusLb">Voltage</span><span class="firstData"> V</span><span class="anotherData"> V</span><span class="anotherData"> V</span></br></span>
		<span class="hide"><span class="lb statusLb">Current</span><span class="firstData"> A</span><span class="anotherData"> A</span><span class="anotherData"> A</span></br></span>
		<span class="hide"><span class="lb statusLb">Frequency</span><span class="firstData"> Hz</span><span class="anotherData"> Hz</span><span class="anotherData"> Hz</span></br></span>
		<span class="hide"><span class="lb statusLb">Power Factor</span><span class="firstData"></span><span class="anotherData"></span><span class="anotherData"></span></br></span>
	</div>
	<span class="caption">Output</span><br/>
	<div class="gap">
		<span class="lb statusLb">Status</span><span class="txt">{{.Output.Status}}</span><br/>
		<span class="hide"><span class="firstItem">Phase1</span><span class="anotherItem">Phase2</span><span class="anotherItem">Phase3</span></br></span>
		<span class="lb statusLb">Voltage</span><span class="firstData">{{.Output.Voltage}} V</span><span class="hide"><span class="anotherData"> V</span><span class="anotherData"> V</span></span><br/>
		<span class=""><span class="lb statusLb ">Frequency</span><span class="firstData ">{{.Output.Frequency}} Hz</span><span class="hide"><span class="anotherData"> Hz</span><span class="anotherData"> Hz</span></span></br></span>
		<span class="lb statusLb">Load</span><span class="firstData">{{.Output.LoadPercent}} % ({{.Output.LoadWatts}} Watts)</span><span class="hide"><span class="anotherData"></span><span class="anotherData"></span></span><br/>
		<span class=""><span class="lb statusLb">Current</span><span class="firstData">{{.Output.Current}} A</span><span class="hide"><span class="anotherData"> A</span><span class="anotherData"> A</span></span></br></span>
		<span class="hide"><span class="lb statusLb">Power Factor</span><span class="firstData"></span><span class="anotherData"></span><span class="anotherData"></span></br></span>
		<span class="hide"><span class="lb statusLb">Active Power</span><span class="firstData"> kW</span><span class="anotherData"> kW</span><span class="anotherData"> kW</span></br></span>
		<span class="hide"><span class="lb statusLb">Apparent Power</span><span class="firstData"> kVA</span><span class="anotherData"> kVA</span><span class="anotherData"> kVA</span></br></span>
		<span class="hide"><span class="lb statusLb">Reactive Power</span><span class="firstData"> kVAr</span><span class="anotherData"> kVAr</span><span class="anotherData"> kVAr</span></br></span>
		<span class="lb statusLb hide">CL</span><span class="hide"><span class="firstData">None</span></br></span>
		<span class="lb statusLb ">NCL </span><span class=""><span class="firstData">On</span></span><br />
		<span class="lb statusLb hide">NCL 2</span><span class="hide"><span class="firstData"></span></span>
		<span class="hide">
			<span class="lb statusLb">Energy</span><span class="firstData txt">0.0 kWh</span><span class="txt2">  ( from 05/14/2021	00:04:26	)</span><br/>
			<form name="Form1" action="status.html" method="get">
				<span class="lb statusLb ">&nbsp;</span><span class="firstData"><input class="" style="font-weight:bold" type="submit" name="SumRST" value="Reset" />&nbsp;</span></br>
			</form>
		</span>
	</div>
	<span class="caption">Battery</span><br/>
	<div class="gap">
		<span class="lb statusLb">Status</span><span class="txt">{{.Battery.Status}}</span><br/>
		<span class="hide"><span class="lb statusLb">Charge Mode</span><span class="txt"></span><br/></span>
		<span class="hide"><span class="lb statusLb">Charge State</span><span class="txt"></span><br/></span>
		<span class="lb statusLb">Remaining Capacity</span><span class="txt">{{.Battery.RemainingCapacity}} %</span><br/>
		<span class="lb statusLb">Remaining Runtime</span><span class="txt">{{.Battery.RemainingRuntime | secToMin}}min. </span><br/>
		<span class="hide"><span class="lb statusLb hide">Voltage</span><span class="firstData hide">0 V</span><span class="hide"> V</span></br></span>
		<span class="hide"><span class="lb statusLb">Current</span><span class="firstData"> A</span><span class="anotherData"> A</span></br></span>
		<span class="hide"><span class="lb statusLb">Temperature</span><span class="firstData">&deg;C</span></br></span>
	</div>
	<span class="caption">System</span><br/>
	<div class="gap">
		<span class="lb statusLb">Status</span><span class="txt">{{.Status}}</span><br/>
		<span class="lb statusLb ">Temperature</span><span class="txt ">{{.TempC}}&deg;C{{.TempF}}&deg;F &nbsp;</span><br />
		<span class="hide"><span class="lb statusLb">Maintenance Breaker</span><span class="txt"></span></br></span>
	</div>
	`
	htmlEnvStatusBodyTemplate = `
	<span class="caption">Information</span><br />
	<div class="gap">
	<span class="lb env_statusLb">Name</span><span class="txt2">{{.Name}}</span><br />
	<span class="lb env_statusLb">Location</span><span class="txt2">{{.Location}}</span><br />
	</div>
	<span class="caption">Temperature</span><br />
	<div class="gap">
	<span class="lb env_statusLb">Current Value</span><span class="txt3">{{.TempF}}</span><span class="txt3">  &deg;F </span><br />
	<span class="lb env_statusLb">Maximum</span><span class="txt2">136.4</span><span class="txt2"> &deg;F       ( at 02/28/2023	18:05:40	)</span><br />
	<span class="lb env_statusLb">Minimum</span><span class="txt2">63.1</span><span class="txt2"> &deg;F       ( at 08/09/2022	13:48:15	)</span><br />	
	<span class=""><input class="env_statusBtn" type="submit" name="TemRST" value="Reset" />&nbsp;</span>
	</div>
	<span class="caption">Humidity</span><br />
	<div class="gap">
	<span class="lb env_statusLb">Current Value</span><span class="txt3">{{.Humidity}} %RH</span><br />
	<span class="lb env_statusLb">Maximum</span><span class="txt2">93 %RH</span><span class="txt2">      ( at 08/09/2022	13:49:40	)</span><br />
	<span class="lb env_statusLb">Minimum</span><span class="txt2">4 %RH</span><span class="txt2">      ( at 02/18/2023	17:11:50	)</span><br />	
	<span class=""><input class="env_statusBtn" type="submit" name="HumRST" value="Reset" />&nbsp;</span>
	</div>
	<span class="caption">Contact</span><br />
	<div class="gap">
	<span class="lb env_statusLb">Contact#1</span><span class="txt">Normal</span><br />
	<span class="lb env_statusLb">Contact#2</span><span class="txt">Normal</span><br />
	<span class="lb env_statusLb">Contact#3</span><span class="txt">Normal</span><br />
	<span class="lb env_statusLb">Contact#4</span><span class="txt">Normal</span><br />
	</div>
	`
)

type Body struct {
	body       string
	base       string
	_processed bool
	_raw       []byte
}

func (b *Body) process() {
	if b.base == "" {
		b.base = htmlDefaultBase
	}

	h := fmt.Sprintf(b.base, b.body)
	b._raw = []byte(h)
	b._processed = true
}

func (b *Body) Reset() {
	b._processed = false
	b._raw = make([]byte, 0)

	b.process()
}

func (b *Body) Read(p []byte) (int, error) {
	if !b._processed {
		b.process()
	}

	rlen := len(b._raw)
	plen := len(p)

	if plen > rlen {
		for i := 0; i < rlen; i++ {
			p[i] = b._raw[i]
		}

		b._raw = make([]byte, 0)

		return rlen, io.EOF
	}

	for i := 0; i < plen; i++ {
		p[i] = b._raw[i]
	}

	b._raw = b._raw[plen:]

	return plen, nil
}
