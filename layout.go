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
	cursorX  float32
	cursorY  float32
	fontSize float32
}

//TODO: turn layout into a struct

func (d *Document) Layout() {
	d.displayList = make([]displayItem, 0)
	display := display{
		cursorX:  20,
		cursorY:  20,
		fontSize: 16,
	}
	/*d.displayList[0] = displayItem{
		text:     d.body,
		font:     fonts[0],
		position: rl.NewVector2(20, 20),
		fontSize: 16,
		color:    rl.Black,
	}*/
	for _, n := range d.document.getChildren() {
		d.recourse(n, &display)
	}

	for _, v := range d.displayList {
		println(v.text)
	}
}

func (d *Document) recourse(treeNode node, display *display) {
	if treeNode.isText() {
		d.displayText(treeNode, display)
	} else {
		//TODO: open and close tag
		if treeNode.getTag() == "h1" { //FIXME: only test
			d.flush(display)
			display.fontSize += 10
		}
		if treeNode.getTag() == "script" || treeNode.getTag() == "head" {
			return
		}
		for _, c := range treeNode.getChildren() {
			d.recourse(c, display)
		}
		if treeNode.getTag() == "h1" || treeNode.getTag() == "br" || treeNode.getTag() == "p" { //FIXME: only test
			d.flush(display)
		}
		if treeNode.getTag() == "h1" {
			display.fontSize -= 10
		}
	}
}

func (d *Document) displayText(n node, display *display) {
	for _, w := range strings.Split(n.getText(), " ") {
		w = strings.TrimSpace(w)
		if len(w) == 0 {
			continue
		}
		wSize := rl.MeasureTextEx(fonts[0], w, display.fontSize, 0)

		if display.cursorX+wSize.X > float32(rl.GetScreenWidth()) {
			d.flush(display)
		}

		d.displayList = append(d.displayList, displayItem{ //TODO: add displayItems to line and to displayList in flush
			text:     w,
			font:     fonts[0],
			position: rl.NewVector2(display.cursorX, display.cursorY),
			fontSize: display.fontSize,
			color:    rl.Black,
		})

		display.cursorX += wSize.X + rl.MeasureTextEx(fonts[0], " ", display.fontSize, 0).X
	}
}

func (d *Document) flush(display *display) {
	display.cursorX = 20
	display.cursorY += rl.MeasureTextEx(fonts[0], " ", display.fontSize, 0).Y
	//display.cursorY += float32(fonts[0].BaseSize)
}
