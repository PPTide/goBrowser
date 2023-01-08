package main

import (
	"errors"
	"fmt"
)

type cssParser struct {
	string string
	i      int
}

func newCssParser(string string) *cssParser {
	return &cssParser{
		string: string,
		i:      0,
	}
}

func (p *cssParser) whitespace() {
	for p.i < len(p.string) && p.string[p.i] == ' ' {
		p.i++
	}
}

func (p *cssParser) word() (string, error) {
	start := p.i
	for p.i < len(p.string) {
		if checkAlphaNum(p.string[p.i]) || arrayContains([]string{"#", "-", ".", "%"}, string(p.string[p.i])) {
			p.i++
		} else {
			break
		}
	}
	if !(p.i > start) {
		//panic("[CSSParser] Error parsing word")
		return "", errors.New("[CSSParser] Error parsing word")
	}
	return p.string[start:p.i], nil
}

func (p *cssParser) literal(literal uint8) error {
	if !(p.i < len(p.string) && p.string[p.i] == literal) { //FIXME: i should use rune here instead of uint8
		//panic("[CSSParser] Error parsing literal")
		return errors.New("[CSSParser] Error parsing literal")
	}
	p.i++
	return nil
}

func (p *cssParser) pair() (string, string, error) {
	var err error
	property, err := p.word()
	if err != nil {
		return "", "", err
	}
	p.whitespace()
	err = p.literal(':')
	if err != nil {
		return "", "", err
	}
	p.whitespace()
	value, err := p.word()
	if err != nil {
		return "", "", err
	}
	return property, value, nil
}

func (p *cssParser) body() map[string]string {
	pairs := map[string]string{}
	for p.i < len(p.string) {
		property, value, err := p.pair()
		if err != nil {
			fmt.Printf("[CSSParser] Ignoring error: %s\n", err)
			why := p.ignoreUntil([]string{";"})
			if why == ';' {
				_ = p.literal(';') //Why tf would there be an error here???
				p.whitespace()
			} else {
				break
			}
		}
		pairs[property] = value
		p.whitespace()
		err = p.literal(';')
		if err != nil {
			//panic("[CSSParser] While parsing body couldn't find ';'")
			println("[CSSParser] While parsing body couldn't find ';'")
		}
		p.whitespace()
	}
	return pairs
}

func (p *cssParser) ignoreUntil(chars []string) uint8 {
	for p.i < len(p.string) {
		if arrayContains(chars, string(p.string[p.i])) {
			return p.string[p.i]
		} else {
			p.i++
		}
	}
	return 0
}

func checkAlphaNum(char uint8) bool {
	if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') && (char < '0' || char > '9') {
		return false
	}
	return true
}

func style(node node) {
	if node.isText() {
		return
	}
	styleMap := node.Style()
	*styleMap = map[string]string{}

	value, ok := node.Attributes()["style"]
	if ok {
		*styleMap = newCssParser(value).body()
	}

	for _, child := range node.getChildren() {
		style(child)
	}
}
