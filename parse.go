package main

import (
	"log"
	"strings"
)

type node interface {
	getParent() node
	getChildren() []node
	addChild(node)
	getText() string
}

type attribute struct{}

type text struct {
	text   string
	parent node
	//children []node
}

type element struct {
	tag        string
	attributes attribute
	parent     node
	children   []node
}

func (e *element) addChild(n node) {
	e.children = append(e.children, n)
}

func (e *text) addChild(n node) {
	//e.children = append(e.children, n)
	log.Panicf("idk why a text node should have children... if your reading this you probably know that...") //TODO: change to return err
}

func (e *text) getText() string {
	return e.text
}

func (e *element) getText() string {
	log.Panicf("idk why a text node should have children... if your reading this you probably know that...") //TODO: change to return err
	return "err"
}

func (t *text) getParent() node {
	return t.parent
}

func (e *element) getParent() node {
	return e.parent
}

func (t *text) getChildren() []node {
	log.Panicf("idk why a text node should have children... if your reading this you probably know that...") //TODO: change to return err
	//return t.children
	return make([]node, 0)
}

func (e *element) getChildren() []node {
	return e.children
}

func (d *document) parseHTML() {
	//TODO: implement Documentation: https://html.spec.whatwg.org/multipage/parsing.html#tokenization
	t := ""
	inTag := false
	r := strings.NewReader(d.body)
	d.document = &element{}
	for r.Len() > 0 {
		c, _, err := r.ReadRune()
		checkErr(err)
		switch c {
		case '<':
			inTag = true
			if len(strings.TrimSpace(t)) > 0 {
				//d.test += text //TODO: add real support for text nodes
				d.document.addChild(&text{
					text: t,
				})
			}
			t = ""
		case '>':
			inTag = false
			//TODO: add support for tag nodes
			t = ""
		default:
			t += string(c)
		}
	}
	if !inTag && len(t) > 0 {
		d.test += t
	}
}
