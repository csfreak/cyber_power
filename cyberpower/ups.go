package cyberpower

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

var (
	upsPath          = "/status_update.html"
	runtimeRegex     = regexp.MustCompile(`^([0-9]+)min.`)
	temperatureRegex = regexp.MustCompile(`^([0-9]+)°C([0-9]+)°F`)
)

func (u *UPS) update() error {
	root, err := u.parent.get(upsPath)
	if err != nil {
		return fmt.Errorf("unable to update UPS on %s; %w", u.parent.getHost(), err)
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
				processUpsGroup(currentGroup, labelGroup, u)
			}
		}

		currentGroup = currentGroup.NextSibling
	}

	return nil
}

func (u UPS) getParent() CyberPower {
	return u.parent
}

func processUpsGroup(group *html.Node, labelGroup *html.Node, u *UPS) {
	currentItem := group.FirstChild

	var labelItem *html.Node

	for {
		if currentItem == nil {
			break
		}

		switch {
		case len(currentItem.Attr) == 0:
		case currentItem.Attr[0].Key == parseAttrKey && strings.Trim(currentItem.Attr[0].Val, " ") == "hide":
		case currentItem.Attr[0].Key == parseAttrKey && strings.Trim(currentItem.Attr[0].Val, " ") == "lb statusLb":
			labelItem = currentItem
		case currentItem.Attr[0].Key == parseAttrKey && strings.Trim(currentItem.Attr[0].Val, " ") == "txt":
			if !(labelItem == nil) {
				switch labelItem.FirstChild.Data {
				case "Status":
					switch labelGroup.FirstChild.Data {
					case parseInputKey:
						u.Input.Status = currentItem.FirstChild.Data
					case parseOutputKey:
						u.Output.Status = currentItem.FirstChild.Data
					case parseBatteryKey:
						u.Battery.Status = currentItem.FirstChild.Data
					case "System":
						u.Status = currentItem.FirstChild.Data
					}
				case "Remaining Capacity":
					cs := currentItem.FirstChild.Data
					cs = strings.Split(cs, " ")[0]

					rc, err := strconv.Atoi(cs)
					if err != nil {
						log.Printf("Unable to parse Remaining Capacity for %s", labelGroup.FirstChild.Data)
						break
					}

					if labelGroup.FirstChild.Data == parseBatteryKey {
						u.Battery.RemainingCapacity = rc
					}
				case "Remaining Runtime":
					rs := runtimeRegex.FindStringSubmatch(currentItem.FirstChild.Data)
					if len(rs) < expectedRemainingRuntimeRegexMatch {
						log.Printf("Unable to parse Remaining Runtime for %s", labelGroup.FirstChild.Data)
						break
					}

					rr, err := strconv.Atoi(rs[1])
					if err != nil {
						log.Printf("Unable to parse Remaining Runtime for %s", labelGroup.FirstChild.Data)
						break
					}

					if labelGroup.FirstChild.Data == parseBatteryKey {
						rr *= secInMin
						u.Battery.RemainingRuntime = rr
					}
				case "Temperature":
					ts := temperatureRegex.FindStringSubmatch(currentItem.FirstChild.Data)
					if len(ts) != expectedTemperatureRegexMatch {
						log.Printf("Unable to parse Tempurature for %s", labelGroup.FirstChild.Data)
						break
					}

					tc, err := strconv.Atoi(ts[1])
					if err != nil {
						log.Printf("Unable to parse Tempurature for %s", labelGroup.FirstChild.Data)
						break
					}

					tf, err := strconv.Atoi(ts[2])
					if err != nil {
						log.Printf("Unable to parse Tempurature for %s", labelGroup.FirstChild.Data)
						break
					}

					if labelGroup.FirstChild.Data == "System" {
						u.TempC = tc
						u.TempF = tf
					}
				}
			}
		case currentItem.Attr[0].Key == parseAttrKey && strings.Trim(currentItem.Attr[0].Val, " ") == "firstData":
			if !(labelItem == nil) {
				switch labelItem.FirstChild.Data {
				case "Voltage":
					vs := currentItem.FirstChild.Data
					vs = strings.Split(vs, " ")[0]

					v, err := strconv.ParseFloat(vs, 64)
					if err != nil {
						log.Printf("Unable to parse Voltage for %s", labelGroup.FirstChild.Data)
						break
					}

					switch labelGroup.FirstChild.Data {
					case parseInputKey:
						u.Input.Voltage = v
					case parseOutputKey:
						u.Output.Voltage = v
					}
				case "Frequency":
					fs := currentItem.FirstChild.Data
					fs = strings.Split(fs, " ")[0]

					f, err := strconv.ParseFloat(fs, 64)
					if err != nil {
						log.Printf("Unable to parse Frequency for %s", labelGroup.FirstChild.Data)
						break
					}

					switch labelGroup.FirstChild.Data {
					case parseInputKey:
						u.Input.Frequency = f
					case parseOutputKey:
						u.Output.Frequency = f
					}
				case "Current":
					cs := currentItem.FirstChild.Data
					cs = strings.Split(cs, " ")[0]

					c, err := strconv.ParseFloat(cs, 64)
					if err != nil {
						log.Printf("Unable to parse Current for %s", labelGroup.FirstChild.Data)
						break
					}

					if labelGroup.FirstChild.Data == parseOutputKey {
						u.Output.Current = c
					}
				case "Load":
					ls := currentItem.FirstChild.Data
					lsplit := strings.Split(ls, " ")

					lp, err := strconv.Atoi(lsplit[0])
					if err != nil {
						log.Printf("Unable to parse Load Percent for %s", labelGroup.FirstChild.Data)
						break
					}

					lw, err := strconv.Atoi(strings.Trim(lsplit[2], "()"))
					if err != nil {
						log.Printf("Unable to parse Load Watts for %s", labelGroup.FirstChild.Data)
						break
					}

					if labelGroup.FirstChild.Data == parseOutputKey {
						u.Output.LoadPercent = lp
						u.Output.LoadWatts = lw
					}
				}
			}
		case currentItem.Attr[0].Key == parseAttrKey && strings.Trim(currentItem.Attr[0].Val, " ") == "":
			processUpsGroup(currentItem, labelGroup, u)
		}

		currentItem = currentItem.NextSibling

		continue
	}
}
