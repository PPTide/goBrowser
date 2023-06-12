package parse

type parser struct {
	tokenizer *tokenizer
}

func newParser(input string) (this *parser) {
	this = &parser{tokenizer: newTokenizer(input)}

	this.tokenizer.setParser(this)
	return
}

func (p parser) run() {
	for true {
		if p.tokenizer.isEOF() {
			break
		}
		p.tokenizer.nextToken()
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
