package main

import rl "github.com/gen2brain/raylib-go/raylib"

func (l *inlineLayout) openTag(tag string) {
	if tag == "h1" { //FIXME: only test
		l.flush(false)
		l.fontSize += 10
	}
	if tag == "title" {
		l.line = []displayItem{}
	}
	if tag == "big" {
		l.fontSize += 4
	}
}

func (l *inlineLayout) closeTag(tag string) {
	if tag == "h1" || tag == "br" || tag == "p" { //FIXME: only test
		l.flush(false)
	}
	if tag == "h1" {
		l.fontSize -= 10
	}
	if tag == "title" {
		text := ""
		for _, item := range l.line {
			text += item.text + " "
		}
		rl.SetWindowTitle(text)
	}
	if tag == "big" {
		l.fontSize -= 4
	}
}
