package main

import (
	"log"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	width  = 600
	height = 450

	hStep = 13
	vStep = 18

	scrollSpeed = 5
)

var (
	fonts          = make([]rl.Font, 1)
	scroll float32 = 0
)

func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(width, height, "Browser")

	args := os.Args[1:]
	pageUrl := ""
	if len(args) > 0 {
		pageUrl = args[0]
	} else {
		pageUrl = "https://en.wikipedia.org/wiki/ELISA"
	}

	fonts[0] = rl.LoadFont("fonts/Arial.ttf")

	d := CreateDocument(pageUrl)
	d.parseHTML()
	d.document = newDocumentLayout(d.nodes)
	d.document.layout()
	d.displayList = make([]displayItem, 0)
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
