package parse

import (
	"strings"
	"testing"
)

func TestTokenizationDom_CharacterReferenceState(t *testing.T) {
	td := tokenizationDom{
		tokens: make([]token, 0),
		reader: strings.NewReader("lt;"),
		tmp:    "",
	}
	td.characterReferenceState()
	if !(td.tokens[0].(characterToken).char == '<') {
		t.FailNow()
	}
}
