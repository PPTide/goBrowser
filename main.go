package main

import (
	"log"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	width  = 600
	height = 450

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
		pageUrl = "https://pptie.de/en/"
	}

	fonts[0] = rl.LoadFont("fonts/Arial.ttf")

	d := CreateDocument(pageUrl)
	d.parseHTML()
	d.Layout()
	for !rl.WindowShouldClose() {
		d.Draw()
	}
}

func checkErr(err error) {

	if err != nil {

		log.Fatal(err)
	}
}
