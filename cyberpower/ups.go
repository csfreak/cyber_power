package cyberpower

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type UPS struct {
	parent  *CyberPower
	input   INPUT_POWER
	output  OUTPUT_POWER
	battery BATTERY_POWER
	temp_c  int
	temp_f  int
	status  string
}

type INPUT_POWER struct {
	status    string
	voltage   float64
	frequency float64
}

type OUTPUT_POWER struct {
	status       string
	voltage      float64
	frequency    float64
	current      float64
	load_watts   int
	load_percent int
}

type BATTERY_POWER struct {
	status             string
	remaining_capacity int
	remaining_runtime  int
}

var ups_path = "/status_update.html"
var runtime_regex = regexp.MustCompile(`^[0-5]+`)
var temperature_regex = regexp.MustCompile(`^([0-9]+)°C([0-9]+)°F`)

func (u *UPS) update() {
	root, err := u.parent.get(ups_path)
	if err != nil {
		log.Printf("Unable to update UPS on %s", u.parent.hostpath)
	}

	body := root.FirstChild.LastChild

	curr_group := body.FirstChild
	var label_group *html.Node
	for {
		if curr_group == nil {
			break
		}
		switch curr_group.Data {
		case "span":
			if curr_group.Attr[0].Key == "class" && curr_group.Attr[0].Val == "caption" {
				label_group = curr_group
			}
		case "div":
			if curr_group.Attr[0].Key == "class" && curr_group.Attr[0].Val == "gap" {
				process_ups_group(curr_group, label_group, u)

			}
		}

		curr_group = curr_group.NextSibling
	}

}

func process_ups_group(group *html.Node, label_group *html.Node, u *UPS) {
	curr_item := group.FirstChild
	var label_item *html.Node
	for {
		if curr_item == nil {
			break
		}
		if len(curr_item.Attr) == 0 {
			curr_item = curr_item.NextSibling
			continue
		} else if curr_item.Attr[0].Key == "class" && strings.Trim(curr_item.Attr[0].Val, " ") == "hide" {
			curr_item = curr_item.NextSibling
			continue
		} else if curr_item.Attr[0].Key == "class" && strings.Trim(curr_item.Attr[0].Val, " ") == "lb statusLb" {
			label_item = curr_item
		} else if curr_item.Attr[0].Key == "class" && strings.Trim(curr_item.Attr[0].Val, " ") == "txt" {
			if !(label_item == nil) {
				switch label_item.FirstChild.Data {
				case "Status":
					switch label_group.FirstChild.Data {
					case "Input":
						u.input.status = curr_item.FirstChild.Data
					case "Output":
						u.output.status = curr_item.FirstChild.Data
					case "Battery":
						u.battery.status = curr_item.FirstChild.Data
					case "System":
						u.status = curr_item.FirstChild.Data
					}
				case "Remaining Capacity":
					cs := curr_item.FirstChild.Data
					cs = strings.Split(cs, " ")[0]
					rc, err := strconv.Atoi(cs)
					if err != nil {
						log.Printf("Unable to parse Remaining Capacity for %s", label_group.FirstChild.Data)
						break
					}
					switch label_group.FirstChild.Data {
					case "Battery":
						u.battery.remaining_capacity = rc
					}
				case "Remaining Runtime":
					rs := runtime_regex.FindString(curr_item.FirstChild.Data)
					rr, err := strconv.Atoi(rs)
					rr = rr * 60
					if err != nil {
						log.Printf("Unable to parse Remaining Runtime for %s", label_group.FirstChild.Data)
						break
					}
					switch label_group.FirstChild.Data {
					case "Battery":
						u.battery.remaining_runtime = rr
					}
				case "Temperature":
					ts := temperature_regex.FindStringSubmatch(curr_item.FirstChild.Data)
					if len(ts) != 3 {
						log.Printf("Unable to parse Tempurature for %s", label_group.FirstChild.Data)
						break
					}
					tc, err := strconv.Atoi(ts[1])
					if err != nil {
						log.Printf("Unable to parse Tempurature for %s", label_group.FirstChild.Data)
						break
					}
					tf, err := strconv.Atoi(ts[2])

					if err != nil {
						log.Printf("Unable to parse Tempurature for %s", label_group.FirstChild.Data)
						break
					}
					switch label_group.FirstChild.Data {
					case "System":
						u.temp_c = tc
						u.temp_f = tf
					}

				}
			}
		} else if curr_item.Attr[0].Key == "class" && strings.Trim(curr_item.Attr[0].Val, " ") == "firstData" {
			if !(label_item == nil) {
				switch label_item.FirstChild.Data {
				case "Voltage":
					vs := curr_item.FirstChild.Data
					vs = strings.Split(vs, " ")[0]
					v, err := strconv.ParseFloat(vs, 64)
					if err != nil {
						log.Printf("Unable to parse Voltage for %s", label_group.FirstChild.Data)
						break
					}
					switch label_group.FirstChild.Data {
					case "Input":
						u.input.voltage = v
					case "Output":
						u.output.voltage = v
					}
				case "Frequency":
					fs := curr_item.FirstChild.Data
					fs = strings.Split(fs, " ")[0]
					f, err := strconv.ParseFloat(fs, 64)
					if err != nil {
						log.Printf("Unable to parse Frequency for %s", label_group.FirstChild.Data)
						break
					}
					switch label_group.FirstChild.Data {
					case "Input":
						u.input.frequency = f
					case "Output":
						u.output.frequency = f
					}
				case "Current":
					cs := curr_item.FirstChild.Data
					cs = strings.Split(cs, " ")[0]
					c, err := strconv.ParseFloat(cs, 64)
					if err != nil {
						log.Printf("Unable to parse Current for %s", label_group.FirstChild.Data)
						break
					}
					switch label_group.FirstChild.Data {
					case "Output":
						u.output.current = c
					}
				case "Load":
					ls := curr_item.FirstChild.Data
					lsplit := strings.Split(ls, " ")
					lp, err := strconv.Atoi(lsplit[0])
					if err != nil {
						log.Printf("Unable to parse Load Percent for %s", label_group.FirstChild.Data)
						break
					}
					lw, err := strconv.Atoi(strings.Trim(lsplit[2], "()"))
					if err != nil {
						log.Printf("Unable to parse Load Watts for %s", label_group.FirstChild.Data)
						break
					}
					switch label_group.FirstChild.Data {
					case "Output":
						u.output.load_percent = lp
						u.output.load_watts = lw
					}
				}
			}
		} else if curr_item.Attr[0].Key == "class" && strings.Trim(curr_item.Attr[0].Val, " ") == "" {
			process_ups_group(curr_item, label_group, u)
		}
		curr_item = curr_item.NextSibling
		continue
	}
}
