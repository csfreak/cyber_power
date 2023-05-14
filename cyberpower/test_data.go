package cyberpower

import (
	"fmt"
	"io"
)

var (
	htmlEmptyBody string = `
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
	htmlUpsStatusBody = `
	<table cellspacing="0" cellpadding="0" class="top" style="border-radius: 12px;">
	<tr> <td colspan="3"></td> </tr> <tr>
	<td class="appName" nowrap="nowrap" valign="center" rowspan="2">
	<div style="margin-left:20px"> <table style="width:300px"><tr><td>
	UPS Remote Management </td><td style='border-right:2px solid #1B60AE'>
	</td></tr> </table> </div> </td> <td>
<div style="margin-right:10px">
	<table style="width:500px">
		<tr>
			<td>
				<table>
				<tr><td>
					<table>
						<tr>
							<td>
								<div align="left" class="userInfo">Administrator login from localhost </div>
							</td>
							<td>
								<span id="admin_viewer_pic"></span>
							</td>
							<td>
								<div align="right" style="font-size: 10px;color: #04B;line-height: 18px;">
									[<a href="logout.html" class="logoutState">Logout</a>]&nbsp;&nbsp;&nbsp;
								</div>
							</td>
						</tr>
					</table>
				</td><td>
							<div style="width:20px;">
				<span title="Change Language">
				<table  onclick="ftable_show()" onmouseover="document.getElementById('flag_down').style.display='block';" onmouseout="document.getElementById('flag_down').style.display='none';">
				<tr>
					<td>
						<span ><span id="now_country_pic"></span></span>
					</td>
					<td>
						<span id="flag_down" style="display:none"><span id="down_pic"></span></span>
					</td>
				</tr>
				<tr>
					<td>
						<table id="flag_table" class="flag_table">
						<tr>
							<td>
							<div style="height:4px;"></div>
							<div id="us_flag"  class="flag" onmouseover="mouse_over(this.id);" onmouseout="mouse_out(this.id);" onclick="select_language(this.id);">
								<table>
									<tr>
										<td><span id="us_pic"></span></td>	
										<td><a style="font-size:12px;">&nbsp;&nbsp;English&nbsp;</a></td>
									</tr>	
								</table>	
							</div>
							<div id="spa_flag"  class="flag" onmouseover="mouse_over(this.id);" onmouseout="mouse_out(this.id);" onclick="select_language(this.id);">
								<table>
									<tr>
										<td><span id="sp_pic"></span></td>	
										<td><a style="font-size:12px;">&nbsp;&nbsp;Español&nbsp;</a></td>
									</tr>	
								</table>	
							</div>
							<div id="fr_flag"  class="flag" onmouseover="mouse_over(this.id);" onmouseout="mouse_out(this.id);" onclick="select_language(this.id);">
								<table>
									<tr>
										<td><span id="fr_pic"></span></td>	
										<td><a style="font-size:12px;">&nbsp;&nbsp;Français&nbsp;</a></td>
									</tr>	
								</table>	
							</div>
							<div id="tw_flag"  class="flag" onmouseover="mouse_over(this.id);" onmouseout="mouse_out(this.id);" onclick="select_language(this.id);">
								<table>
									<tr>
										<td><span id="roc_pic"></span></td>	
										<td><a style="font-size:12px;">&nbsp;&nbsp;繁體中文</a></td>
									</tr>	
								</table>	
							</div>
							<div id="de_flag"  class="flag" onmouseover="mouse_over(this.id);" onmouseout="mouse_out(this.id);" onclick="select_language(this.id);">
								<table>
									<tr>
										<td><span id="ge_pic"></span></td>	
										<td><a style="font-size:12px;">&nbsp;&nbsp;Deutsch&nbsp;</a></td>
									</tr>	
								</table>	
							</div>
							<div id="zh_cn_flag"  class="flag" onmouseover="mouse_over(this.id);" onmouseout="mouse_out(this.id);" onclick="select_language(this.id);">
								<table>
									<tr>
										<td><span id="cn_pic"></span></td>	
										<td><a style="font-size:12px;">&nbsp;&nbsp;简体中文</a></td>
									</tr>	
								</table>	
							</div>
							</td>
						</tr>
						</table>
					</td>
				</tr>
				</table>
				</span>
				</div>
				</td></tr>
				</table>
			</td>
		</tr>
	</table>
</div>
		</td>
		<td rowspan="2">
			<div align="right" style="margin-right:5px;">
				<a href=http://www.cyberpowersystems.com/index.html target=_BLANK><span id="logo_pic" title="Cyber Power Systems, Inc."></span></a>
			</div>
		</td>

	</tr>
	<tr>
		<td>
			
			<ul class="nav">
				<li ><div class="leftEdge"></div><a href="summary.html">Summary</a><div class="rightEdge"></div></li>
				<li class="tabHit"><div class="leftEdge"></div><a href="status.html">UPS</a><div class="rightEdge"></div></li>
				<li><div class="leftEdge"></div><a href="env_status.html">Envir</a><div class="rightEdge"></div></li>
				<li ><div class="leftEdge"></div><a href="logs.html">Log</a><div class="rightEdge"></div></li>
				<li><div class="leftEdge"></div><a href="date.html">System</a><div class="rightEdge"></div></li>
				<li id="noSeperator">
					<div></div><a href="/help/status.html" target="w" onclick="var w=open('/help/status.html','w','width=840,height=560,menubar=no,directories=no,resizable=yes,scrollbars=yes'); w.focus();return false;">Help</a>
				</li>
				</br><div style="border:5px solid #FFF"></div>

			</ul>
		</td>
	</tr>
	<tr>
		<td colspan="3"></td>
	</tr>
</table>
<div class="main">
	<table cellspacing="0" cellpadding="0" width="100%">
		<tr>
			<td valign="top">
				<div class="menu">
					<a id="SubMenu3" href="status.html" class="item">Status</a>
					<a id="SubMenu12" href="battstatus.html" class="item" style="display:none">Battery Status</a>
					<a id="SubMenu11" href="module_status.html" class="item" style="display:none">Module Status</a>
					<a id="SubMenu4" href="info.html" class="item">Information</a>
					<a id="SubMenu5" href="config.html" class="item" >Configuration</a>
					<a id="SubMenu6" href="switch.html" class="item" >Master Switch</a>
					<a id="SubMenu7" href="outlet_bank.html" class="item" >Bank Control</a>
					<a id="SubMenu8" href="diagnostics.html" class="item" >Diagnostics</a>
					<a id="SubMenu9" href="schedule.html" class="item" >Schedule</a>
					<div  id="SubName1"></div>
					<div id="SubMenu1" class="SubMenu" style='display:none'>
						<a href="wol.html" class="titem">Features</a>
						<a href="wol_list.html" class="titem">Lists</a>
					</div>
					<div  id="SubName2"></div>
					<div id="SubMenu2" class="SubMenu" style='display:none'>
						<a href="energywise_config.html" class="titem">Configuration</a>
						<a href="energywise.html" class="titem">Node List</a>
					</div>
					<a id="SubMenu10" href="clients.html" class="widthLitem" >PowerPanel<sup>&reg;</sup> List</a>

				</div>
				<script>
					LoadMenu();
					ShowSubMenu(0,3);
				</script>
				<div class="bottomSharp">&nbsp;</div>
			</td>
			<td style="" width="100%" class="workspace" valign="top">
					<div ></div>
<div class="header">Status</div>
<center><span class="lb hide"><font color="red">Communication has not been established.</font></span></center>
<div class="content " id="content">
	<span class="caption">Input</span><br/>
	<div class="gap">
		<span class="lb statusLb">Status</span><span class="txt">Normal</span><br/>
		<span class="hide"><span class="firstItem">Phase1</span><span class="anotherItem">Phase2</span><span class="anotherItem">Phase3</span></br></span>
		<span class="lb statusLb">Voltage</span><span class="firstData">119.0 V</span><span class="hide"><span class="anotherData"> V</span><span class="anotherData"> V</span></span><br/>
		<span class=""><span class="lb statusLb">Frequency</span><span class="firstData">60.0 Hz</span><span class="hide"><span class="anotherData"> Hz</span><span class="anotherData"> Hz</span></span></br></span>
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
		<span class="lb statusLb">Status</span><span class="txt">Normal</span><br/>
		<span class="hide"><span class="firstItem">Phase1</span><span class="anotherItem">Phase2</span><span class="anotherItem">Phase3</span></br></span>
		<span class="lb statusLb">Voltage</span><span class="firstData">119.0 V</span><span class="hide"><span class="anotherData"> V</span><span class="anotherData"> V</span></span><br/>
		<span class=""><span class="lb statusLb ">Frequency</span><span class="firstData ">60.0 Hz</span><span class="hide"><span class="anotherData"> Hz</span><span class="anotherData"> Hz</span></span></br></span>
		<span class="lb statusLb">Load</span><span class="firstData">41 % (615 Watts)</span><span class="hide"><span class="anotherData"></span><span class="anotherData"></span></span><br/>
		<span class=""><span class="lb statusLb">Current</span><span class="firstData">5.0 A</span><span class="hide"><span class="anotherData"> A</span><span class="anotherData"> A</span></span></br></span>
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
		<span class="lb statusLb">Status</span><span class="txt">Fully Charged</span><br/>
		<span class="hide"><span class="lb statusLb">Charge Mode</span><span class="txt"></span><br/></span>
		<span class="hide"><span class="lb statusLb">Charge State</span><span class="txt"></span><br/></span>
		<span class="lb statusLb">Remaining Capacity</span><span class="txt">100 %</span><br/>
		<span class="lb statusLb">Remaining Runtime</span><span class="txt">19min. </span><br/>
		<span class="hide"><span class="lb statusLb hide">Voltage</span><span class="firstData hide">0 V</span><span class="hide"> V</span></br></span>
		<span class="hide"><span class="lb statusLb">Current</span><span class="firstData"> A</span><span class="anotherData"> A</span></br></span>
		<span class="hide"><span class="lb statusLb">Temperature</span><span class="firstData">&deg;C</span></br></span>
	</div>
	<span class="caption">System</span><br/>
	<div class="gap">
		<span class="lb statusLb">Status</span><span class="txt">Normal</span><br/>
		<span class="lb statusLb ">Temperature</span><span class="txt ">25&deg;C77&deg;F &nbsp;</span><br />
		<span class="hide"><span class="lb statusLb">Maintenance Breaker</span><span class="txt"></span></br></span>
	</div>
</div>
</td></tr>
<tr >
	<td></td>
	<td colspan="2" style="background:#FFF;">
	</td>
</tr>
</table></div><div class="footer">&copy; 2010-2018, CyberPower Systems, Inc. All rights reserved.</div>
<div id="ch_pass" style="display:none" class="ch_def_pass">
	<form method="post" action="default_setup.cgi">
		<input id="direct_url" name="direct_url" type="hidden" value="test"></input>
		<span class="txt"><b style="color:blach">You have logged in with the factory default User Name and Password. For enhanced security, please setup a new User Name and Password.</b></span></br>
		<table width="100%">
			<tr>
				<td>
					&nbsp;
				</td>
			</tr>
			<tr>
				<td>
					<table width="95%" align="right">
						<tr style="height:25px;">
							<td>
								<span id="account_string" class="txt" style="color:blach">New User Name</span>
							</td>
							<td>
								<input id="def_account" maxlength="63" name="def_account" type="text" width="15"></input><span class="txt">&nbsp;[1-63 characters]</span>
							</td>
						</tr>
						<tr style="height:25px;">
							<td>
								<span id="def_new_string" class="txt" style="color:blach">New Password</span>
							</td>
							<td>
								<input id="def_new_pass" maxlength="63" name="def_new_pass" type="password" width="15"></input><span class="txt">&nbsp;[1-63 characters]</span>
							</td>
						</tr>
						<tr style="height:25px;">
							<td>
								<span id="def_new_vstring" class="txt" style="color:blach">Confirm Password</span>
							</td>
							<td>
								<input id="def_new_vpass" maxlength="63" name="def_new_vpass" type="password" width="15"></input><span class="txt">&nbsp;[1-63 characters]</span>
							</td>
						</tr>
					</table>
				</td>
			</tr>
			<tr>
				<td align="center">
					</br><input type="button" value="Apply" style="width:100px" onclick="default_setup(this.form)"></br>
					<div id="default_change_msg" class="txt"></div>
				</td>
			</tr>
		</table>
	</form>
</div>
<div id="out_border" style="display:none" class="ch_def_out_border">
</div>
<div id="ch_upgrade" style="display: none; z-index: -1;height:100px;width:500px;" class="ch_def_pass">
	<div id="ups_fwup_status">UPS firmware waiting for updates</div>
</div>
	`
)

type Body struct {
	body       string
	_processed bool
	_raw       []byte
}

func (b *Body) process() {
	h := fmt.Sprintf(htmlEmptyBody, b.body)
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
