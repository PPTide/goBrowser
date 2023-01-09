package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"log"
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
	inlineNodes := make([]node, 0)
	for _, child := range l.node.getChildren() {
		var next layout
		if layoutMode(child) == "inline" {
			inlineNodes = append(inlineNodes, child)
			//next = newInlineLayout(child, l, previous)
		} else {
			if len(inlineNodes) > 0 {
				next = newInlineLayout(inlineNodes, l, previous)
				inlineNodes = make([]node, 0)
				previous = next
				l.children = append(l.children, next)
			}
			next = newBlockLayout(child, l, previous)
			previous = next
			l.children = append(l.children, next)
		}
	}
	if len(inlineNodes) > 0 {
		next := newInlineLayout(inlineNodes, l, previous)
		l.children = append(l.children, next)
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
	nodes    []node
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
	for _, node := range l.nodes {
		xOffset = 0
		if node.isText() {
			continue
		}
		bgcolor := rl.Color{A: 0}
		//----------------------Start-Handling-Tags--------------------FIXME: temporary solution
		if node.getTag() == "head" {
			return
		}
		if node.getTag() == "nav" && node.Attributes()["class"] == "links" {
			bgcolor = rl.LightGray
		}
		if node.getTag() == "li" {
			*drawList = append(*drawList, drawRect{
				rect: rl.Rectangle{
					X:      l.x + 2, //FIXME: This doesn't work with line wraps
					Y:      l.y + 5, //FIXME: This only works for fontsize 16
					Width:  4,
					Height: 4,
				},
				color: rl.Black,
			})
			xOffset += 8
		}
		//----------------------End-Handling-Tags--------------------
		bgcolorTmp, ok := (*node.Style())["background-color"]
		if ok {
			switch bgcolorTmp {
			case "lightblue":
				bgcolor = rl.NewColor(173, 216, 230, 255)
			default:
				log.Panicf("[InlineLayout paint()] background-color is '%s', which is not supported", bgcolorTmp)
			}
		}
		if bgcolor.A != 0 {
			*drawList = append(*drawList, drawRect{
				rect: rl.Rectangle{ //FIXME: This uses the position and size of all nodes, not only the one wanted
					X:      l.x,
					Y:      l.y,
					Width:  l.width,
					Height: l.height,
				},
				color: bgcolor,
			})
		}
	}
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

func newInlineLayout(nodes []node, parent layout, previous layout) *inlineLayout {
	return &inlineLayout{
		nodes:    nodes,
		parent:   parent,
		previous: previous,
	}
}

func (l *inlineLayout) layout() {
	l.displayList = make([]displayItem, 0)
	l.line = make([]displayItem, 0)

	l.width = l.parent.Width()
	l.x = l.parent.X()

	if l.previous != nil {
		l.y = l.previous.Y() + l.previous.Height() //Maybe they don't
	} else {
		l.y = l.parent.Y()
	}

	l.cursorX = l.x
	l.cursorY = l.y
	l.fontSize = 16
	for _, node := range l.nodes {
		if !node.isText() && (node.getTag() == "head" || node.getTag() == "script" || node.getTag() == "style") {
			continue
		}

		l.recourse(node)
	}
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
		wSize := rl.MeasureTextEx(fonts[0].getSize(int32(l.fontSize)), w, l.fontSize, 0)

		if l.cursorX+wSize.X > l.width {
			l.flush()
		}

		l.line = append(l.line, displayItem{
			text:     w,
			font:     fonts[0].getSize(int32(l.fontSize)),
			position: rl.NewVector2(l.cursorX, l.cursorY),
			fontSize: l.fontSize,
			color:    rl.Black,
		})

		l.cursorX += wSize.X + rl.MeasureTextEx(fonts[0].getSize(int32(l.fontSize)), " ", l.fontSize, 0).X
	}
}

func (l *inlineLayout) flush() {
	if len(l.line) < 1 {
		return
	}
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
	l.cursorY += rl.MeasureTextEx(fonts[0].getSize(int32(biggestItem.fontSize)), " ", biggestItem.fontSize, 0).Y //move the cursor down for the next line
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

func (l *documentLayout) paint(drawList *[]drawItem) {
	l.children[0].paint(drawList)
	//------------------Draw-scroll-bar------------------------
	if l.height > float32(rl.GetRenderHeight()) {
		*drawList = append(*drawList, newDrawScroll(*l))
	}
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
	}
	if len(n.getChildren()) > 0 {
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

type drawScroll struct {
	document documentLayout

	width  float32
	height float32
	x      float32
	y      float32
}

func newDrawScroll(document documentLayout) *drawScroll {
	return &drawScroll{document: document}
}

func (dScroll drawScroll) Type() string {
	return "scroll"
}

func (dScroll drawScroll) Text() drawText {
	panic("tried to get drawText from drawScroll")
}

func (dScroll drawScroll) Rect() drawRect {
	panic("tried to get drawRect from drawScroll")
}

func (dScroll drawScroll) Execute() {
	dScroll.width = 8
	dScroll.x = float32(rl.GetRenderWidth()) - dScroll.width

	screenContentRatio := float32(rl.GetRenderHeight()) / dScroll.document.height
	dScroll.height = screenContentRatio * float32(rl.GetRenderHeight())
	dScroll.y = -scroll * screenContentRatio
	rec := rl.NewRectangle(dScroll.x, dScroll.y, dScroll.width, dScroll.height)
	rl.DrawRectangleRec(rec, rl.Blue)
}
