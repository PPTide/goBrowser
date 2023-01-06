package main

func (l *inlineLayout) openTag(tag string) {
	if tag == "h1" { //FIXME: only test
		l.flush()
		l.fontSize += 10
	}
	if tag == "script" || tag == "style" || tag == "head" {
		return //FIXME: this doesn't work here xD
	}
	if tag == "big" {
		l.fontSize += 4
	}
}

func (l *inlineLayout) closeTag(tag string) {
	if tag == "h1" || tag == "br" || tag == "p" { //FIXME: only test
		l.flush()
	}
	if tag == "h1" {
		l.fontSize -= 10
	}
	if tag == "big" {
		l.fontSize -= 4
	}
}
