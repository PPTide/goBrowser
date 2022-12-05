package main

import (
	"log"
	"strings"
)

var SELF_CLOSING_TAGS = [...]string{
	"area",
	"base",
	"br",
	"col",
	"embed",
	"hr",
	"img",
	"input",
	"link",
	"meta",
	"param",
	"source",
	"track",
	"wbr",
}

type node interface {
	getParent() node
	getChildren() []node
	addChild(node)
	getText() string
	getTag() string
	isText() bool
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

func (e *element) isText() bool {
	return false
}

func (t *text) isText() bool {
	return true
}

func (e *element) getTag() string {
	return e.tag
}

func (e *text) getTag() string {
	log.Panicf("idk why a text node should have children... if your reading this you probably know that...") //TODO: change to return err
	return "err"
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
	log.Panicf("idk why an element node should have text... if your reading this you probably know that...") //TODO: change to return err
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
	currentText := ""
	inTag := false
	r := strings.NewReader(d.body)
	d.document = &element{}
	for r.Len() > 0 {
		c, _, err := r.ReadRune()
		checkErr(err)
		switch c {
		case '<':
			inTag = true
			d.addText(currentText)
			currentText = ""
		case '>':
			inTag = false
			d.addTag(currentText)
			currentText = ""
		default:
			currentText += string(c)
		}
	}
	if !inTag && len(currentText) > 0 {
		d.addText(currentText)
	}
	if len(d.unfinished) != 1 {
		log.Panic("Not one unfinished node")
	}
	d.document.addChild(d.unfinished[0])
}

func (d *document) addText(t string) {
	if len(strings.TrimSpace(t)) == 0 {
		return
	}

	d.implicitTags("")

	parent := d.unfinished[len(d.unfinished)-1]
	parent.addChild(&text{
		text:   t,
		parent: parent,
	})
}

func (d *document) implicitTags(tag string) {
	for {
		openTags := make([]string, 0)
		for _, v := range d.unfinished {
			openTags = append(openTags, v.getTag())
		}
		if len(openTags) == 0 && tag != "html" {
			d.addTag("html")
		} else {
			break
		}
		//TODO: add other cases
	}
}

func arrayContains(array []string, contains string) bool {
	for _, v := range array {
		if v == contains {
			return true
		}
	}
	return false
}

func (d *document) addTag(tag string) {
	if tag[0] == '!' { //TODO: This is temporary (it will stay here for ever :-) )
		return
	}
	//TODO: split off attributes for real
	tag = strings.Split(tag, " ")[0]
	d.implicitTags(tag)
	//TODO: add support for ending tags and self closing tags
	if tag[0] == '/' { //closing Tag
		if len(d.unfinished) == 1 {
			return
		}
		n := d.unfinished[len(d.unfinished)-1]
		d.unfinished = d.unfinished[:len(d.unfinished)-1]
		parent := d.unfinished[len(d.unfinished)-1]
		parent.addChild(n)
		return
	}

	if tag[len(tag)-1] == '/' || arrayContains(SELF_CLOSING_TAGS[:], tag) { //self closing tag
		parent := d.unfinished[len(d.unfinished)-1]
		n := &element{
			tag:    tag,
			parent: parent, //TODO: add attribute
		}
		parent.addChild(n)
		return
	}

	//starter tag
	parent := d.document
	if len(d.unfinished) >= 1 {
		parent = d.unfinished[len(d.unfinished)-1]
	}
	n := &element{
		tag:    tag,
		parent: parent, //TODO: add attribute
	}
	d.unfinished = append(d.unfinished, n)
}
