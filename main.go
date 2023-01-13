package main

import (
	"goBrowser/data"
	"log"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//TODO: add support for these thingies: "&lt", "&gt"

const (
	width  = 600
	height = 450

	hStep = 13
	vStep = 18

	scrollSpeed = 5
)

type fontBook struct {
	fonts map[int32]rl.Font
	font  string
}

func (fb *fontBook) getSize(size int32) rl.Font {
	font, ok := fb.fonts[size]
	if ok {
		return font
	}
	font = rl.LoadFontEx(fb.font, size, nil)
	fb.fonts[size] = font
	return font
}

func newFontBook(font string) *fontBook {
	return &fontBook{
		font:  font,
		fonts: map[int32]rl.Font{},
	}
}

var (
	fonts          = make([]*fontBook, 1)
	scroll float32 = 0

	entities map[string]data.Entity
)

func main() { //TODO: maybe switch to sdl2
	entities = data.GetEntities()

	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(width, height, "Browser")

	args := os.Args[1:]
	pageUrl := ""
	if len(args) > 0 {
		pageUrl = args[0]
	} else {
		pageUrl = "https://browser.engineering/styles.html"
	}

	fonts[0] = newFontBook("fonts/Arial.ttf")

	d := CreateDocument(pageUrl)
	d.parseHTML()
	style(d.nodes)
	d.document = newDocumentLayout(d.nodes)
	d.document.layout()
	d.displayList = make([]drawItem, 0)
	d.document.paint(&d.displayList)
	for !rl.WindowShouldClose() {
		d.Draw()
	}
}

func checkErr(err error) {

	if err != nil {

		log.Fatal(err)
	}
}
