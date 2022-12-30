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
	line     []displayItem
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
	d.flush(&display)

	/*for _, v := range d.displayList {
		println(v.text)
	}*/
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
		if treeNode.getTag() == "script" || treeNode.getTag() == "style" || treeNode.getTag() == "head" {
			return
		}
		if treeNode.getTag() == "big" {
			display.fontSize += 4
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
		if treeNode.getTag() == "big" {
			display.fontSize -= 4
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

		display.line = append(display.line, displayItem{
			text:     w,
			font:     fonts[0],
			position: rl.NewVector2(display.cursorX, display.cursorY),
			fontSize: display.fontSize,
			color:    rl.Black,
		})

		display.cursorX += wSize.X + rl.MeasureTextEx(fonts[0], " ", display.fontSize, 0).X
	}
}

func (d *Document) flush(display *display) { //idk why this works... it's magic don't question it
	var biggestItem displayItem
	for _, item := range display.line { // go through every item to find the biggest font size
		if item.fontSize > biggestItem.fontSize { //FIXME: I should really be searching for the biggest height.
			biggestItem = item
		}
	}
	biggestHeight := rl.MeasureTextEx(biggestItem.font, " ", biggestItem.fontSize, 0).Y //calculate the height
	for _, item := range display.line {
		//move every item down by the difference in size between the biggest one in the line and the item and add it to displayList
		itemHeight := rl.MeasureTextEx(item.font, " ", item.fontSize, 0).Y
		item.position.Y += biggestHeight - itemHeight
		d.displayList = append(d.displayList, item)
	}
	display.line = []displayItem{}
	display.cursorX = 20
	display.cursorY += rl.MeasureTextEx(fonts[0], " ", biggestItem.fontSize, 0).Y //move the cursor down for the next line
	//display.cursorY += float32(fonts[0].BaseSize)
}
