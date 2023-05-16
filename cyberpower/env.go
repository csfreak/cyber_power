package cyberpower

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

var envPath = "/env_status_update.html"

func (e *ENV) update() error {
	root, err := e.parent.get(envPath)
	if err != nil {
		return fmt.Errorf("unable to update ENV on %s; %w", e.parent.getHost(), err)
	}

	body := root.FirstChild.LastChild
	currentGroup := body.FirstChild

	var labelGroup *html.Node

	for {
		if currentGroup == nil {
			break
		}

		switch currentGroup.Data {
		case "span":
			if currentGroup.Attr[0].Key == parseAttrKey && currentGroup.Attr[0].Val == "caption" {
				labelGroup = currentGroup
			}
		case "div":
			if currentGroup.Attr[0].Key == parseAttrKey && currentGroup.Attr[0].Val == "gap" {
				processEnvGroup(currentGroup, labelGroup, e)
			}
		}

		currentGroup = currentGroup.NextSibling
	}

	return nil
}

func (e ENV) getParent() CyberPower {
	return e.parent
}

func processEnvGroup(group *html.Node, labelGroup *html.Node, e *ENV) {
	currentItem := group.FirstChild

	var labelItem *html.Node

ItemIter:
	for {
		switch {
		case currentItem == nil:
			break ItemIter
		case len(currentItem.Attr) == 0:
		case currentItem.Attr[0].Key == parseAttrKey && strings.Trim(currentItem.Attr[0].Val, " ") == "hide":
		case currentItem.Attr[0].Key == parseAttrKey && strings.Trim(currentItem.Attr[0].Val, " ") == "lb env_statusLb":
			labelItem = currentItem
		case currentItem.Attr[0].Key == parseAttrKey && strings.Trim(currentItem.Attr[0].Val, " ") == "txt2":
			if !(labelItem == nil) {
				switch labelItem.FirstChild.Data {
				case "Name":
					e.Name = currentItem.FirstChild.Data
				case "Location":
					e.Location = currentItem.FirstChild.Data
				}
			}
		case currentItem.Attr[0].Key == parseAttrKey && strings.Trim(currentItem.Attr[0].Val, " ") == "txt3":
			if !(labelItem == nil) {
				if labelItem.FirstChild.Data == "Current Value" {
					switch labelGroup.FirstChild.Data {
					case "Temperature":
						ts := currentItem.FirstChild.Data
						ts = strings.Trim(ts, " ")
						t, err := strconv.ParseFloat(ts, 64)
						if err != nil {
							if !(currentItem.PrevSibling.Attr[0].Key == parseAttrKey && strings.Trim(currentItem.PrevSibling.Attr[0].Val, " ") == "txt3") {
								log.Printf("Unable to parse Temp Value for %s", labelGroup.FirstChild.Data)
							}
							break
						}
						e.TempF = t
					case "Humidity":
						hs := currentItem.FirstChild.Data
						hs = strings.Split(hs, " ")[0]
						f, err := strconv.Atoi(hs)
						if err != nil {
							log.Printf("Unable to parse Humidity Value for %s", labelGroup.FirstChild.Data)
							break
						}
						e.Humidity = f
					}
				}
			}
		}
		currentItem = currentItem.NextSibling
	}
}
