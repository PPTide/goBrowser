package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"strings"
)

var blockElements = []string{"html", "body", "article", "section", "nav", "aside",
	"h1", "h2", "h3", "h4", "h5", "h6", "hgroup", "header",
	"footer", "address", "p", "hr", "pre", "blockquote",
	"ol", "ul", "menu", "li", "dl", "dt", "dd", "figure",
	"figcaption", "main", "div", "table", "form", "fieldset",
	"legend", "details", "summary"}

type displayItem struct {
	text     string
	font     rl.Font
	position rl.Vector2
	fontSize float32
	color    rl.Color
}

type drawItem interface {
	Type() string
	Text() drawText
	Rect() drawRect
	Execute()
}

type drawRect struct {
	rect  rl.Rectangle
	color rl.Color
}

func (dRect drawRect) Execute() {
	dRect.rect.Y += scroll
	rl.DrawRectangleRec(dRect.rect, dRect.color)
}

func (dRect drawRect) Type() string {
	return "rect"
}

func (dRect drawRect) Text() drawText {
	panic("tried to get drawText from drawRect")
}

func (dRect drawRect) Rect() drawRect {
	return dRect
}

type drawText struct {
	text     string
	font     rl.Font
	position rl.Vector2
	fontSize float32
	color    rl.Color
}

func (dText drawText) Execute() {
	dText.position.Y = dText.position.Y + scroll
	rl.DrawTextEx(dText.font, dText.text, dText.position, dText.fontSize, 0, dText.color)
}

func (dText drawText) Type() string {
	return "text"
}

func (dText drawText) Text() drawText {
	return dText
}

func (dText drawText) Rect() drawRect {
	panic("tried to get drawRect from drawText")
}

type layout interface {
	layout()
	paint(*[]drawItem)

	Width() float32
	Height() float32
	X() float32
	Y() float32
}

type blockLayout struct {
	node     node
	parent   layout
	previous layout
	children []layout

	width  float32
	height float32
	x      float32
	y      float32
}

func (l *blockLayout) paint(displayList *[]drawItem) {
	for _, child := range l.children {
		child.paint(displayList)
	}
}

func (l *blockLayout) X() float32 {
	return l.x
}

func (l *blockLayout) Y() float32 {
	return l.y
}

func (l *blockLayout) Width() float32 {
	return l.width
}

func (l *blockLayout) Height() float32 {
	return l.height
}

func newBlockLayout(node node, parent layout, previous layout) *blockLayout {
	return &blockLayout{node: node, parent: parent, previous: previous}
}

func (l *blockLayout) layout() {
	if !l.node.isText() && l.node.getTag() == "head" {
		return
	}
	var previous layout = nil
	for _, child := range l.node.getChildren() {
		var next layout
		if layoutMode(child) == "inline" {
			next = newInlineLayout(child, l, previous)
		} else {
			next = newBlockLayout(child, l, previous)
		}
		l.children = append(l.children, next)
		previous = next
	}
	l.width = l.parent.Width()
	l.x = l.parent.X()
	if l.previous != nil {
		l.y = l.previous.Y() + l.previous.Height()
	} else {
		l.y = l.parent.Y()
	}
	for _, child := range l.children {
		child.layout()
	}
	l.height = 0
	for _, child := range l.children {
		l.height += child.Height()
	}
}

type inlineLayout struct {
	node     node
	parent   layout
	previous layout
	children []layout

	width  float32
	height float32
	x      float32
	y      float32

	displayList []displayItem
	cursorX     float32
	cursorY     float32
	fontSize    float32
	line        []displayItem
	//TODO: add support for `style` and `weight`
}

func (l *inlineLayout) paint(drawList *[]drawItem) {
	var xOffset float32 = 0
	//----------------------Start-Handling-Tags--------------------FIXME: temporary solution
	if !l.node.isText() && l.node.getTag() == "head" {
		return
	}
	if !l.node.isText() && l.node.getTag() == "pre" {
		*drawList = append(*drawList, drawRect{
			rect: rl.Rectangle{
				X:      l.x,
				Y:      l.y,
				Width:  l.width,
				Height: l.height,
			},
			color: rl.Gray,
		})
	}
	if !l.node.isText() && (l.node.getTag() == "nav" && l.node.Attributes()["class"] == "links") {
		*drawList = append(*drawList, drawRect{
			rect: rl.Rectangle{
				X:      l.x,
				Y:      l.y,
				Width:  l.width,
				Height: l.height,
			},
			color: rl.LightGray,
		})
	}
	if !l.node.isText() && l.node.getTag() == "li" {
		*drawList = append(*drawList, drawRect{
			rect: rl.Rectangle{
				X:      l.x + 2,
				Y:      l.y + 5, //FIXME: This only works for fontsize 16
				Width:  4,
				Height: 4,
			},
			color: rl.Black,
		})
		xOffset += 8
	}
	//----------------------End-Handling-Tags--------------------
	for _, item := range l.displayList {
		item.position.X += xOffset
		*drawList = append(*drawList, drawText{
			text:     item.text,
			font:     item.font,
			position: item.position,
			fontSize: item.fontSize,
			color:    item.color,
		})
	}
}

func (l *inlineLayout) X() float32 {
	return l.x
}

func (l *inlineLayout) Y() float32 {
	return l.y
}

func (l *inlineLayout) Width() float32 {
	return l.width
}

func (l *inlineLayout) Height() float32 {
	return l.height
}

func newInlineLayout(node node, parent layout, previous layout) *inlineLayout {
	return &inlineLayout{
		node:     node,
		parent:   parent,
		previous: previous,
	}
}

func (l *inlineLayout) layout() {
	if !l.node.isText() && (l.node.getTag() == "head" || l.node.getTag() == "script" || l.node.getTag() == "style") {
		return
	}
	l.width = l.parent.Width()
	l.x = l.parent.X()
	if l.previous != nil {
		l.y = l.previous.Y() + l.previous.Height()
	} else {
		l.y = l.parent.Y()
	}

	l.displayList = make([]displayItem, 0)

	l.cursorX = l.x
	l.cursorY = l.y
	l.fontSize = 16

	l.line = make([]displayItem, 0)

	l.recourse(l.node)
	l.flush()
	l.height = l.cursorY - l.y
}

func (l *inlineLayout) recourse(n node) {
	if n.isText() {
		l.displayText(n)
	} else {
		l.openTag(n.getTag())
		for _, child := range n.getChildren() {
			l.recourse(child)
		}
		l.closeTag(n.getTag())
	}
}

func (l *inlineLayout) displayText(n node) {
	for _, w := range strings.Split(n.getText(), " ") {
		w = strings.TrimSpace(w)
		if len(w) == 0 {
			continue
		}
		wSize := rl.MeasureTextEx(fonts[0], w, l.fontSize, 0)

		if l.cursorX+wSize.X > l.width {
			l.flush()
		}

		l.line = append(l.line, displayItem{
			text:     w,
			font:     fonts[0],
			position: rl.NewVector2(l.cursorX, l.cursorY),
			fontSize: l.fontSize,
			color:    rl.Black,
		})

		l.cursorX += wSize.X + rl.MeasureTextEx(fonts[0], " ", l.fontSize, 0).X
	}
}

func (l *inlineLayout) flush() {
	var biggestItem displayItem
	for _, item := range l.line { // go through every item to find the biggest font size
		if item.fontSize > biggestItem.fontSize { //FIXME: I should really be searching for the biggest height.
			biggestItem = item
		}
	}
	biggestHeight := rl.MeasureTextEx(biggestItem.font, " ", biggestItem.fontSize, 0).Y //calculate the height
	for _, item := range l.line {
		//move every item down by the difference in size between the biggest one in the line and the item and add it to displayList
		itemHeight := rl.MeasureTextEx(item.font, " ", item.fontSize, 0).Y
		item.position.Y += biggestHeight - itemHeight
		l.displayList = append(l.displayList, item)
	}
	l.line = []displayItem{}
	l.cursorX = l.x
	l.cursorY += rl.MeasureTextEx(fonts[0], " ", biggestItem.fontSize, 0).Y //move the cursor down for the next line
	//display.cursorY += float32(fonts[0].BaseSize)
}

type documentLayout struct {
	node     node
	children []layout

	width  float32
	height float32
	x      float32
	y      float32
}

func (l *documentLayout) paint(displayList *[]drawItem) {
	l.children[0].paint(displayList)
}

func newDocumentLayout(node node) *documentLayout {
	return &documentLayout{node: node}
}

func (l *documentLayout) X() float32 {
	return l.x
}

func (l *documentLayout) Y() float32 {
	return l.y
}

func (l *documentLayout) Width() float32 {
	return l.width
}

func (l *documentLayout) Height() float32 {
	return l.height
}

func (l *documentLayout) layout() {
	child := newBlockLayout(l.node, l, nil)
	l.children = []layout{child}
	l.width = float32(rl.GetRenderWidth()) - 2*hStep
	l.x = hStep
	l.y = vStep
	child.layout()
	l.height = child.height + 2*vStep
}

func layoutMode(n node) string {
	if n.isText() {
		return "inline"
	} else if len(n.getChildren()) > 0 {
		for _, child := range n.getChildren() {
			if child.isText() {
				continue
			}
			if arrayContains(blockElements, child.getTag()) {
				return "block"
			}
		}
		return "inline"
	} else {
		return "block"
	}
}
