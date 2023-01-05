package main

import (
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Document struct {
	headers     map[string]string
	body        string
	nodes       node
	document    layout
	displayList []displayItem
	unfinished  []node
}

func CreateDocument(pageUrl string) *Document {
	head, body, err := request(pageUrl)
	checkErr(err)

	var status = make([]string, 2)
	if head == nil {
		status[1] = "200"
	} else {
		status = strings.Split(head[0], " ")
	}

	headers := make(map[string]string)
	if head != nil {
		for _, header := range head[1:] {
			h := strings.SplitN(header, ": ", 2)
			headers[strings.ToLower(h[0])] = h[1]
		}
	}

	switch status[1] {
	case "200":
		break
	case "301":
	case "302":
		locationURL := headers["location"]
		return CreateDocument(locationURL)
	default:
		panic("Got unsupported status code: " + strings.Join(status[1:], " "))
	}

	d := Document{
		headers: headers,
		body:    body,
	}

	return &d
}

func (d *Document) Draw() {
	//TODO: make better
	if rl.IsWindowResized() {
		d.document.layout()
		d.displayList = make([]displayItem, 0)
		d.document.paint(&d.displayList)
	}

	rl.BeginDrawing()

	rl.ClearBackground(rl.White)

	updateScroll()

	for _, x := range d.displayList {
		x.position.Y = x.position.Y + scroll
		rl.DrawTextEx(x.font, x.text, x.position, x.fontSize, 0, x.color)
	}

	rl.EndDrawing()
}

func updateScroll() {
	scroll += rl.GetMouseWheelMove() * scrollSpeed
	if scroll > 0 {
		scroll = 0
	}
}
