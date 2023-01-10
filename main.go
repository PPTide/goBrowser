package main

import (
	"encoding/json"
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

	entitiesTmp map[string]map[string]interface{}
	entities    map[string]entity
)

type entity struct {
	characters string
	codepoints []float64
}

func main() { //TODO: maybe switch to sdl2
	entitiesJSON, err := os.ReadFile("data/entities.json")
	checkErr(err)
	err = json.Unmarshal(entitiesJSON, &entitiesTmp)
	checkErr(err)
	entities = make(map[string]entity, 0)
	//Transform map[string]map[string]interface{} into map[string]entity
	for key1, entityTmp := range entitiesTmp {
		e := entity{}
		for key, value := range entityTmp {
			if key == "characters" {
				e.characters = value.(string)
				continue
			}
			if key == "codepoints" {
				codepointsTmp := value.([]interface{})
				e.codepoints = make([]float64, 0)
				for _, codepoint := range codepointsTmp {
					e.codepoints = append(e.codepoints, codepoint.(float64))
				}
			}
		}
		entities[key1] = e
	}

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
