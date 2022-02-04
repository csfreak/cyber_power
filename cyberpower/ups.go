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
	Input   INPUT_POWER
	Output  OUTPUT_POWER
	Battery BATTERY_POWER
	Temp_c  int
	Temp_f  int
	Status  string
}

type INPUT_POWER struct {
	Status    string
	Voltage   float64
	Frequency float64
}

type OUTPUT_POWER struct {
	Status      string
	Voltage     float64
	Frequency   float64
	Current     float64
	LoadWatts   int
	LoadPercent int
}

type BATTERY_POWER struct {
	Status            string
	RemainingCapacity int
	RemainingRuntime  int
}

var ups_path = "/status_update.html"
var runtime_regex = regexp.MustCompile(`^([0-9]+)min.`)
var temperature_regex = regexp.MustCompile(`^([0-9]+)°C([0-9]+)°F`)

func (u *UPS) update() {
	root, err := u.parent.get(ups_path)
	if err == nil {

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
	} else {
		log.Printf("Unable to update UPS on %s", u.parent.hostpath)
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
						u.Input.Status = curr_item.FirstChild.Data
					case "Output":
						u.Output.Status = curr_item.FirstChild.Data
					case "Battery":
						u.Battery.Status = curr_item.FirstChild.Data
					case "System":
						u.Status = curr_item.FirstChild.Data
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
						u.Battery.RemainingCapacity = rc
					}
				case "Remaining Runtime":
					rs := runtime_regex.FindStringSubmatch(curr_item.FirstChild.Data)
					if len(rs) < 2 {
						log.Printf("Unable to parse Remaining Runtime for %s", label_group.FirstChild.Data)
						break
					}
					rr, err := strconv.Atoi(rs[1])
					rr = rr * 60
					if err != nil {
						log.Printf("Unable to parse Remaining Runtime for %s", label_group.FirstChild.Data)
						break
					}
					switch label_group.FirstChild.Data {
					case "Battery":
						u.Battery.RemainingRuntime = rr
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
						u.Temp_c = tc
						u.Temp_f = tf
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
						u.Input.Voltage = v
					case "Output":
						u.Output.Voltage = v
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
						u.Input.Frequency = f
					case "Output":
						u.Output.Frequency = f
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
						u.Output.Current = c
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
						u.Output.LoadPercent = lp
						u.Output.LoadWatts = lw
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
