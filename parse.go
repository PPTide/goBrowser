package main

import (
	"fmt"
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

var RAWTEXT_TYPES = [...]string{
	"script",
	"style",
}

type node interface {
	getParent() node
	getChildren() []node
	addChild(node)
	getText() string
	getTag() string
	isText() bool
	printTree(int)
}

type text struct {
	text   string
	parent node
	//children []node
}

type element struct {
	tag        string
	attributes map[string]string
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

func (t *text) getTag() string {
	log.Panicf("idk why a text node should have children... if your reading this you probably know that...") //TODO: change to return err
	return "err"
}

func (e *element) addChild(n node) {
	e.children = append(e.children, n)
}

func (t *text) addChild(n node) {
	//e.children = append(e.children, n)
	log.Panicf("idk why a text node should have children... if your reading this you probably know that...") //TODO: change to return err
}

func (t *text) getText() string {
	return t.text
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

func (d *Document) parseHTML() {
	//TODO: implement Documentation: https://html.spec.whatwg.org/multipage/parsing.html#tokenization
	currentText := ""
	inTag := false
	inRawtext := false
	rawtextType := ""
	inComment := false
	r := strings.NewReader(d.body)
	d.document = &element{}
	for r.Len() > 0 {
		c, _, err := r.ReadRune()
		checkErr(err)
		if inComment && c == '-' {
			r1, _, err := r.ReadRune()
			checkErr(err)
			r2, _, err := r.ReadRune()
			checkErr(err)
			if r1 == '-' && r2 == '>' {
				inComment = false
				currentText = ""
			} else {
				currentText += string(c)
				_, err = r.Seek(-2, 1)
				checkErr(err)
			}
			continue
		}
		if inComment {
			currentText += string(c)
			continue
		}
		switch c {
		case '<':
			inTag = true
			d.addText(currentText)
			r1, _, err := r.ReadRune()
			checkErr(err)
			r2, _, err := r.ReadRune()
			checkErr(err)
			r3, _, err := r.ReadRune()
			checkErr(err)
			if r1 == '!' && r2 == '-' && r3 == '-' {
				inComment = true
			} else {
				_, err = r.Seek(-3, 1)
				checkErr(err)
			}
			currentText = ""
		case '>':
			inTag = false
			if inRawtext && currentText != "/"+rawtextType {
				currentText = "<" + currentText + ">"
				continue
			}
			d.addTag(currentText)
			/* TODO: if added tag has special use here switch to the appropriate state:
			 *  https://html.spec.whatwg.org/multipage/parsing.html#parsing-html-fragments:rawtext-state
			 */
			if arrayContains(RAWTEXT_TYPES[:], strings.Split(currentText, " ")[0]) {
				inRawtext = true
				rawtextType = strings.Split(currentText, " ")[0]
			}
			if inRawtext && currentText == "/"+rawtextType {
				inRawtext = false
			}
			currentText = ""
		default:
			currentText += string(c)
		}
	}
	if !inTag && len(currentText) > 0 {
		d.addText(currentText)
	}
	d.finishParsing()
}

func (d *Document) finishParsing() {
	if len(d.unfinished) == 0 {
		d.addTag("html")
	}
	for len(d.unfinished) > 1 {
		n := d.unfinished[len(d.unfinished)-1]
		d.unfinished = d.unfinished[:len(d.unfinished)-1]
		parent := d.unfinished[len(d.unfinished)-1]
		parent.addChild(n)
	}
	d.document.addChild(d.unfinished[0])
	d.unfinished[0].printTree(0)
}

func (d *Document) addText(t string) {
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

func (d *Document) implicitTags(tag string) {
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

func getAttributes(attribString string) (attribs map[string]string) {
	if len(attribString) == 0 {
		return nil
	}
	attribs = map[string]string{}
	attribPairs := strings.Split(attribString, " ")
	for i := 0; i < len(attribPairs); i++ {
		attribPair := attribPairs[i]
		if !strings.Contains(attribPair, "=") {
			attribs[attribPair] = ""
			continue
		}
		split := strings.SplitN(attribPair, "=", 2)
		value := split[1]
		if len(value) > 2 && (value[0] == '"' || value[0] == '\'') {
			for value[len(value)-1] != value[0] {
				i++
				value += " " + attribPairs[i]
			}
			value = value[1 : len(value)-1]
		}
		attribs[split[0]] = value
	}
	return
}

func (d *Document) addTag(tag string) {
	if tag[0] == '!' { //TODO: This is temporary (it will stay here for ever :-) )
		return
	}
	splitTag := strings.SplitN(tag, " ", 2)
	tagName := splitTag[0]
	if len(splitTag) < 2 {
		splitTag = append(splitTag, "")
	}
	attributes := getAttributes(splitTag[1])
	d.implicitTags(tagName)
	if len(tagName) == 0 {
		panic("ahhhhhh")
	}
	if tagName[0] == '/' { //closing Tag
		if len(d.unfinished) == 1 {
			return
		}
		n := d.unfinished[len(d.unfinished)-1]
		d.unfinished = d.unfinished[:len(d.unfinished)-1]
		parent := d.unfinished[len(d.unfinished)-1]
		parent.addChild(n)
		return
	}

	if tag[len(tag)-1] == '/' || arrayContains(SELF_CLOSING_TAGS[:], tagName) { //self closing tag
		parent := d.unfinished[len(d.unfinished)-1]
		n := &element{
			tag:        tagName,
			parent:     parent,
			attributes: attributes,
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
		tag:        tagName,
		parent:     parent,
		attributes: attributes,
	}
	d.unfinished = append(d.unfinished, n)
}

func (e *element) printTree(indent int) {
	fmt.Printf("%s<%s %s>\n", strings.Repeat(" ", indent), e.tag, e.attributes)
	for _, child := range e.children {
		child.printTree(indent + 2)
	}
}

func (t *text) printTree(indent int) {
	println(strings.Repeat(" ", indent) + strings.ReplaceAll(t.text, "\n", "\\n"))
}
