package cyberpower

import (
	"log"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type ENV struct {
	parent   *CyberPower
	Name     string
	Location string
	Temp_f   float64
	Humidity int
}

var env_path = "/env_status_update.html"

func (e *ENV) update() {
	root, err := e.parent.get(env_path)
	if err != nil {
		log.Printf("Unable to update ENV on %s", e.parent.hostpath)
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
				process_env_group(curr_group, label_group, e)

			}
		}

		curr_group = curr_group.NextSibling
	}

}

func process_env_group(group *html.Node, label_group *html.Node, e *ENV) {
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
		} else if curr_item.Attr[0].Key == "class" && strings.Trim(curr_item.Attr[0].Val, " ") == "lb env_statusLb" {
			label_item = curr_item
		} else if curr_item.Attr[0].Key == "class" && strings.Trim(curr_item.Attr[0].Val, " ") == "txt2" {
			if !(label_item == nil) {
				switch label_item.FirstChild.Data {
				case "Name":
					e.Name = curr_item.FirstChild.Data
				case "Location":
					e.Location = curr_item.FirstChild.Data

				}
			}
		} else if curr_item.Attr[0].Key == "class" && strings.Trim(curr_item.Attr[0].Val, " ") == "txt3" {
			if !(label_item == nil) {
				switch label_item.FirstChild.Data {
				case "Current Value":
					switch label_group.FirstChild.Data {
					case "Temperature":
						ts := curr_item.FirstChild.Data
						ts = strings.Trim(ts, " ")
						t, err := strconv.ParseFloat(ts, 64)
						if err != nil {
							if !(curr_item.PrevSibling.Attr[0].Key == "class" && strings.Trim(curr_item.PrevSibling.Attr[0].Val, " ") == "txt3") {
								log.Printf("Unable to parse Current Value for %s", label_group.FirstChild.Data)
							}
							break
						}
						e.Temp_f = t
					case "Humidity":
						hs := curr_item.FirstChild.Data
						hs = strings.Split(hs, " ")[0]
						f, err := strconv.Atoi(hs)
						if err != nil {
							log.Printf("Unable to parse Current Value for %s", label_group.FirstChild.Data)
							break
						}
						e.Humidity = f
					}

				}
			}
		}
		curr_item = curr_item.NextSibling
		continue
	}
}
