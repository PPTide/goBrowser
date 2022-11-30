package main

import (
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type displayItem struct {
	text     string
	font     rl.Font
	position rl.Vector2
	fontSize float32
	color    rl.Color
}

type display struct {
	cursor_x float32
}

func (d *document) Layout() {
	d.displayList = make([]displayItem, 0)
	display := display{
		cursor_x: 20,
	}
	/*d.displayList[0] = displayItem{
		text:     d.body,
		font:     fonts[0],
		position: rl.NewVector2(20, 20),
		fontSize: 16,
		color:    rl.Black,
	}*/
	for _, n := range d.document.getChildren() {
		if true { //TODO: check if text
			d.addText(n, &display)
		}
	}

	for _, v := range d.displayList {
		println(v.text)
	}
}

func (d *document) addText(n node, display *display) {
	for _, w := range strings.Split(n.getText(), " ") {
		w = strings.TrimSpace(w)
		if len(w) == 0 {
			continue
		}
		wSize := rl.MeasureTextEx(fonts[0], w, 16, 0)

		d.displayList = append(d.displayList, displayItem{
			text:     w,
			font:     fonts[0],
			position: rl.NewVector2(display.cursor_x, 20),
			fontSize: 16,
			color:    rl.Black,
		})

		display.cursor_x += wSize.X + rl.MeasureTextEx(fonts[0], " ", 16, 0).X
	}
}
