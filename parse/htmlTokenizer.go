package parse

import "strings"

const EOF = -1

type tokenizer struct {
	state        HTMLState
	input        *strings.Reader
	parser       *parser
	returnState  HTMLState
	currentToken HTMLToken
	tokens       []HTMLToken
	buffer       string
}
type HTMLToken interface {
	asTag() HTMLTokenTag
}
type HTMLTokenTag interface {
	TagName() *string
}
type HTMLTokenCharacter struct {
	codePoint int
}

func (t HTMLTokenCharacter) asTag() HTMLTokenTag {
	//TODO implement me
	panic("implement me")
}

type HTMLTokenEOF struct {
}

func (t HTMLTokenEOF) asTag() HTMLTokenTag {
	//TODO implement me
	panic("implement me")
}

type HTMLTokenStartTag struct {
	tagName string
}

func (t HTMLTokenStartTag) TagName() *string {
	//TODO implement me
	panic("implement me")
}

func (t HTMLTokenStartTag) asTag() HTMLTokenTag {
	//TODO implement me
	panic("implement me")
}

type HTMLTokenEndTag struct {
	tagName string
}

func (t HTMLTokenEndTag) TagName() *string {
	//TODO implement me
	panic("implement me")
}

func (t HTMLTokenEndTag) asTag() HTMLTokenTag {
	//TODO implement me
	panic("implement me")
}

type HTMLTokenComment struct {
	data string
}
type HTMLTokenDoctype struct {
	name string
}
type HTMLTokenParseError struct {
}

func newTokenizer(input string) *tokenizer {
	return &tokenizer{
		input: strings.NewReader(input),
	}
}

func (t *tokenizer) setParser(parser *parser) {
	t.parser = parser
}

type HTMLState int

const (
	HTMLState_Data HTMLState = iota
	HTMLState_RCDATA
	HTMLState_RAWTEXT
	HTMLState_ScriptData
	HTMLState_PLAINTEXT
	HTMLState_TagOpen
	HTMLState_EndTagOpen
	HTMLState_TagName
	HTMLState_RCDATALessThanSign
	HTMLState_RCDATAEndTagOpen
	HTMLState_RCDATAEndTagName
	HTMLState_RAWTEXTLessThanSign
	HTMLState_RAWTEXTEndTagOpen
	HTMLState_RAWTEXTEndTagName
	HTMLState_ScriptDataLessThanSign
	HTMLState_ScriptDataEndTagOpen
	HTMLState_ScriptDataEndTagName
	HTMLState_ScriptDataEscapeStart
	HTMLState_ScriptDataEscapeStartDash
	HTMLState_ScriptDataEscaped
	HTMLState_ScriptDataEscapedDash
	HTMLState_ScriptDataEscapedDashDash
	HTMLState_ScriptDataEscapedLessThanSign
	HTMLState_ScriptDataEscapedEndTagOpen
	HTMLState_ScriptDataEscapedEndTagName
	HTMLState_ScriptDataDoubleEscapeStart
	HTMLState_ScriptDataDoubleEscaped
	HTMLState_ScriptDataDoubleEscapedDash
	HTMLState_ScriptDataDoubleEscapedDashDash
	HTMLState_ScriptDataDoubleEscapedLessThanSign
	HTMLState_ScriptDataDoubleEscapeEnd
	HTMLState_BeforeAttributeName
	HTMLState_AttributeName
	HTMLState_AfterAttributeName
	HTMLState_BeforeAttributeValue
	HTMLState_AttributeValueDoubleQuoted
	HTMLState_AttributeValueSingleQuoted
	HTMLState_AttributeValueUnquoted
	HTMLState_AfterAttributeValueQuoted
	HTMLState_SelfClosingStartTag
	HTMLState_BogusComment
	HTMLState_MarkupDeclarationOpen
	HTMLState_CommentStart
	HTMLState_CommentStartDash
	HTMLState_Comment
	HTMLState_CommentLessThanSign
	HTMLState_CommentLessThanSignBang
	HTMLState_CommentLessThanSignBangDash
	HTMLState_CommentLessThanSignBangDashDash
	HTMLState_CommentEndDash
	HTMLState_CommentEnd
	HTMLState_CommentEndBang
	HTMLState_DOCTYPE
	HTMLState_BeforeDOCTYPEName
	HTMLState_DOCTYPEName
	HTMLState_AfterDOCTYPEName
	HTMLState_AfterDOCTYPEPublicKeyword
	HTMLState_BeforeDOCTYPEPublicIdentifier
	HTMLState_DOCTYPEPublicIdentifierDoubleQuoted
	HTMLState_DOCTYPEPublicIdentifierSingleQuoted
	HTMLState_AfterDOCTYPEPublicIdentifier
	HTMLState_BetweenDOCTYPEPublicAndSystemIdentifiers
	HTMLState_AfterDOCTYPESystemKeyword
	HTMLState_BeforeDOCTYPESystemIdentifier
	HTMLState_DOCTYPESystemIdentifierDoubleQuoted
	HTMLState_DOCTYPESystemIdentifierSingleQuoted
	HTMLState_AfterDOCTYPESystemIdentifier
	HTMLState_BogusDOCTYPE
	HTMLState_CDATASection
	HTMLState_CDATASectionBracket
	HTMLState_CDATASectionEnd
	HTMLState_CharacterReference
	HTMLState_NamedCharacterReference
	HTMLState_AmbiguousAmpersand
	HTMLState_NumericCharacterReference
	HTMLState_HexadecimalCharacterReferenceStart
	HTMLState_DecimalCharacterReferenceStart
	HTMLState_HexadecimalCharacterReference
	HTMLState_DecimalCharacterReference
	HTMLState_NumericCharacterReferenceEnd
)

type HTMLParseError int

const (
	HTMLParseError_AbruptClosingOfEmptyComment HTMLParseError = iota
	HTMLParseError_AbruptDoctypePublicIdentifier
	HTMLParseError_AbruptDoctypeSystemIdentifier
	HTMLParseError_AbsenceOfDigitsInNumericCharacterReference
	HTMLParseError_CdataInHtmlContent
	HTMLParseError_CharacterReferenceOutsideUnicodeRange
	HTMLParseError_ControlCharacterReference
	HTMLParseError_EndTagWithAttributes
	HTMLParseError_duplicateAttribute
	HTMLParseError_endTagWithTrailingSolidus
	HTMLParseError_EofBeforeTagName
	HTMLParseError_EofInCdata
	HTMLParseError_EofInComment
	HTMLParseError_EofInDoctype
	HTMLParseError_EofInScriptHtmlCommentLikeText
	HTMLParseError_EofInTag
	HTMLParseError_IncorrectlyClosedComment
	HTMLParseError_IncorrectlyOpenedComment
	HTMLParseError_InvalidCharacterSequenceAfterDoctypeName
	HTMLParseError_InvalidFirstCharacterOfTagName
	HTMLParseError_MissingAttributeValue
	HTMLParseError_MissingDoctypeName
	HTMLParseError_MissingDoctypePublicIdentifier
	HTMLParseError_MissingDoctypeSystemIdentifier
	HTMLParseError_MissingEndTagName
	HTMLParseError_MissingQuoteBeforeDoctypePublicIdentifier
	HTMLParseError_MissingQuoteBeforeDoctypeSystemIdentifier
	HTMLParseError_MissingSemicolonAfterCharacterReference
	HTMLParseError_MissingWhitespaceAfterDoctypePublicKeyword
	HTMLParseError_MissingWhitespaceAfterDoctypeSystemKeyword
	HTMLParseError_MissingWhitespaceBeforeDoctypeName
	HTMLParseError_MissingWhitespaceBetweenAttributes
	HTMLParseError_MissingWhitespaceBetweenDoctypePublicAndSystemIdentifiers
	HTMLParseError_NestedComment
	HTMLParseError_NoncharacterCharacterReference
	HTMLParseError_NoncharacterInInputStream
	HTMLParseError_NonVoidHtmlElementStartTagWithTrailingSolidus
	HTMLParseError_NullCharacterReference
	HTMLParseError_SurrogateCharacterReference
	HTMLParseError_SurrogateInInputStream
	HTMLParseError_UnexpectedCharacterAfterDoctypeSystemIdentifier
	HTMLParseError_UnexpectedCharacterInAttributeName
	HTMLParseError_UnexpectedCharacterInUnquotedAttributeValue
	HTMLParseError_UnexpectedEqualsSignBeforeAttributeName
	HTMLParseError_UnexpectedNullCharacter
	HTMLParseError_UnexpectedQuestionMarkInsteadOfTagName
	HTMLParseError_UnexpectedSolidusInTag
	HTMLParseError_UnknownNamedCharacterReference
)

func (e HTMLParseError) String() string {
	//FIXME: fuck this is ugly
	switch e {
	case HTMLParseError_AbruptClosingOfEmptyComment:
		return "abrupt-closing-of-empty-comment"
	case HTMLParseError_AbruptDoctypePublicIdentifier:
		return "abrupt-doctype-public-identifier"
	case HTMLParseError_AbruptDoctypeSystemIdentifier:
		return "abrupt-doctype-system-identifier"
	case HTMLParseError_AbsenceOfDigitsInNumericCharacterReference:
		return "absence-of-digits-in-numeric-character-reference"
	case HTMLParseError_CdataInHtmlContent:
		return "cdata-in-html-content"
	case HTMLParseError_CharacterReferenceOutsideUnicodeRange:
		return "character-reference-outside-unicode-range"
	case HTMLParseError_ControlCharacterReference:
		return "control-character-reference"
	case HTMLParseError_EndTagWithAttributes:
		return "end-tag-with-attributes"
	case HTMLParseError_duplicateAttribute:
		return "duplicate-attribute"
	case HTMLParseError_endTagWithTrailingSolidus:
		return "end-tag-with-trailing-solidus"
	case HTMLParseError_EofBeforeTagName:
		return "eof-before-tag-name"
	case HTMLParseError_EofInCdata:
		return "eof-in-cdata"
	case HTMLParseError_EofInComment:
		return "eof-in-comment"
	case HTMLParseError_EofInDoctype:
		return "eof-in-doctype"
	case HTMLParseError_EofInScriptHtmlCommentLikeText:
		return "eof-in-script-html-comment-like-text"
	case HTMLParseError_EofInTag:
		return "eof-in-tag"
	case HTMLParseError_IncorrectlyClosedComment:
		return "incorrectly-closed-comment"
	case HTMLParseError_IncorrectlyOpenedComment:
		return "incorrectly-opened-comment"
	case HTMLParseError_InvalidCharacterSequenceAfterDoctypeName:
		return "invalid-character-sequence-after-doctype-name"
	case HTMLParseError_InvalidFirstCharacterOfTagName:
		return "invalid-first-character-of-tag-name"
	case HTMLParseError_MissingAttributeValue:
		return "missing-attribute-value"
	case HTMLParseError_MissingDoctypeName:
		return "missing-doctype-name"
	case HTMLParseError_MissingDoctypePublicIdentifier:
		return "missing-doctype-public-identifier"
	case HTMLParseError_MissingDoctypeSystemIdentifier:
		return "missing-doctype-system-identifier"
	case HTMLParseError_MissingEndTagName:
		return "missing-end-tag-name"
	case HTMLParseError_MissingQuoteBeforeDoctypePublicIdentifier:
		return "missing-quote-before-doctype-public-identifier"
	case HTMLParseError_MissingQuoteBeforeDoctypeSystemIdentifier:
		return "missing-quote-before-doctype-system-identifier"
	case HTMLParseError_MissingSemicolonAfterCharacterReference:
		return "missing-semicolon-after-character-reference"
	case HTMLParseError_MissingWhitespaceAfterDoctypePublicKeyword:
		return "missing-whitespace-after-doctype-public-keyword"
	case HTMLParseError_MissingWhitespaceAfterDoctypeSystemKeyword:
		return "missing-whitespace-after-doctype-system-keyword"
	case HTMLParseError_MissingWhitespaceBeforeDoctypeName:
		return "missing-whitespace-before-doctype-name"
	case HTMLParseError_MissingWhitespaceBetweenAttributes:
		return "missing-whitespace-between-attributes"
	case HTMLParseError_MissingWhitespaceBetweenDoctypePublicAndSystemIdentifiers:
		return "missing-whitespace-between-doctype-public-and-system-identifiers"
	case HTMLParseError_NestedComment:
		return "nested-comment"
	case HTMLParseError_NoncharacterCharacterReference:
		return "noncharacter-character-reference"
	case HTMLParseError_NoncharacterInInputStream:
		return "noncharacter-in-input-stream"
	case HTMLParseError_NonVoidHtmlElementStartTagWithTrailingSolidus:
		return "non-void-html-element-start-tag-with-trailing-solidus"
	case HTMLParseError_NullCharacterReference:
		return "null-character-reference"
	case HTMLParseError_SurrogateCharacterReference:
		return "surrogate-character-reference"
	case HTMLParseError_SurrogateInInputStream:
		return "surrogate-in-input-stream"
	case HTMLParseError_UnexpectedCharacterAfterDoctypeSystemIdentifier:
		return "unexpected-character-after-doctype-system-identifier"
	case HTMLParseError_UnexpectedCharacterInAttributeName:
		return "unexpected-character-in-attribute-name"
	case HTMLParseError_UnexpectedCharacterInUnquotedAttributeValue:
		return "unexpected-character-in-unquoted-attribute-value"
	case HTMLParseError_UnexpectedEqualsSignBeforeAttributeName:
		return "unexpected-equals-sign-before-attribute-name"
	case HTMLParseError_UnexpectedNullCharacter:
		return "unexpected-null-character"
	case HTMLParseError_UnexpectedQuestionMarkInsteadOfTagName:
		return "unexpected-question-mark-instead-of-tag-name"
	case HTMLParseError_UnexpectedSolidusInTag:
		return "unexpected-solidus-in-tag"
	case HTMLParseError_UnknownNamedCharacterReference:
		return "unknown-named-character-reference"
	}
	return "unknown"
}

func (t *tokenizer) nextToken() {
	for true {
	start:
		currentInputCharacter := t.nextCodePoint()
		switch t.state {
		case HTMLState_Data:
			if currentInputCharacter == '&' {
				t.returnState = HTMLState_Data
				t.state = HTMLState_CharacterReference
				goto start
			}
			if currentInputCharacter == '<' {
				t.state = HTMLState_TagOpen
				goto start
			}
			if currentInputCharacter == 0 {
				t.parseError(HTMLParseError_UnexpectedNullCharacter)
				t.emitCharacter(currentInputCharacter)
			}
			if currentInputCharacter == -1 {
				t.emitEOF()
				return
			}
			t.emitCharacter(currentInputCharacter)
		case HTMLState_RCDATA:
			if currentInputCharacter == '&' {
				t.returnState = HTMLState_RCDATA
				t.state = HTMLState_CharacterReference
				goto start
			}
			if currentInputCharacter == '<' {
				t.state = HTMLState_RCDATALessThanSign
				goto start
			}
			if currentInputCharacter == 0 {
				t.parseError(HTMLParseError_UnexpectedNullCharacter)
				t.emitCharacter(0xFFFD)
			}
			if currentInputCharacter == -1 {
				t.emitEOF()
				return
			}
			t.emitCharacter(currentInputCharacter)
		case HTMLState_RAWTEXT:
			if currentInputCharacter == '<' {
				t.state = HTMLState_RAWTEXTLessThanSign
				goto start
			}
			if currentInputCharacter == 0 {
				t.parseError(HTMLParseError_UnexpectedNullCharacter)
				t.emitCharacter(0xFFFD)
			}
			if currentInputCharacter == -1 {
				t.emitEOF()
				return
			}
			t.emitCharacter(currentInputCharacter)
		case HTMLState_ScriptData:
			if currentInputCharacter == '<' {
				t.state = HTMLState_ScriptDataLessThanSign
				goto start
			}
			if currentInputCharacter == 0 {
				t.parseError(HTMLParseError_UnexpectedNullCharacter)
				t.emitCharacter(0xFFFD)
			}
			if currentInputCharacter == -1 {
				t.emitEOF()
				return
			}
			t.emitCharacter(currentInputCharacter)
		case HTMLState_PLAINTEXT:
			if currentInputCharacter == 0 {
				t.parseError(HTMLParseError_UnexpectedNullCharacter)
				t.emitCharacter(0xFFFD)
			}
			if currentInputCharacter == -1 {
				t.emitEOF()
				return
			}
			t.emitCharacter(currentInputCharacter)
		case HTMLState_TagOpen:
			if currentInputCharacter == '!' {
				t.state = HTMLState_MarkupDeclarationOpen
				goto start
			}
			if currentInputCharacter == '/' {
				t.state = HTMLState_EndTagOpen
				goto start
			}
			if isASCIIAlpha(currentInputCharacter) {
				t.emitTokenStartTag()
				t.reconsumeIn(HTMLState_TagName)
				goto start
			}
			if currentInputCharacter == '?' {
				t.parseError(HTMLParseError_UnexpectedQuestionMarkInsteadOfTagName)
				t.state = HTMLState_BogusComment
				goto start
			}
			if currentInputCharacter == -1 {
				t.parseError(HTMLParseError_EofBeforeTagName)
				t.emitCharacter('<')
				t.emitEOF()
				return
			}
			t.parseError(HTMLParseError_InvalidFirstCharacterOfTagName)
			t.emitCharacter('<')
			t.reconsumeIn(HTMLState_Data)
			goto start
		case HTMLState_EndTagOpen:
			if isASCIIAlpha(currentInputCharacter) {
				t.emitTokenEndTag("")
				t.reconsumeIn(HTMLState_TagName)
				goto start
			}
		case HTMLState_TagName:
			if isCharacterTabulation(currentInputCharacter) || isLineFeed(currentInputCharacter) || isFormFeed(currentInputCharacter) || isSpaceCharacter(currentInputCharacter) {
				t.state = HTMLState_BeforeAttributeName
				goto start
			}
			if currentInputCharacter == '/' {
				t.state = HTMLState_SelfClosingStartTag
				goto start
			}
			if currentInputCharacter == '>' {
				t.emitCurrentToken()
				t.state = HTMLState_Data
				goto start
			}
			if isASCIIUpper(currentInputCharacter) {
				*t.currentToken.asTag().TagName() += string(toASCIILower(currentInputCharacter))
				goto start
			}
			if currentInputCharacter == 0 {
				t.parseError(HTMLParseError_UnexpectedNullCharacter)
				*t.currentToken.asTag().TagName() += string(0xFFFD)
				goto start
			}
			if currentInputCharacter == -1 {
				t.parseError(HTMLParseError_EofInTag)
				t.emitEOF()
				return
			}
			*t.currentToken.asTag().TagName() += string(currentInputCharacter)
		case HTMLState_RCDATALessThanSign:
			if currentInputCharacter == '/' {
				t.buffer = ""
				t.state = HTMLState_RCDATAEndTagOpen
				goto start
			}
			t.emitCharacter('<')
			t.reconsumeIn(HTMLState_RCDATA)
		case HTMLState_RCDATAEndTagOpen:
			if isASCIIAlpha(currentInputCharacter) {
				//FIXME: Create a new end tag token, set its tag name to the empty string
				t.state = HTMLState_RCDATAEndTagName
				goto start
			}
			t.emitCharacter('<')
			t.emitCharacter('/')
			t.reconsumeIn(HTMLState_RCDATA)
		case HTMLState_RCDATAEndTagName:
			if isCharacterTabulation(currentInputCharacter) || isLineFeed(currentInputCharacter) || isFormFeed(currentInputCharacter) || isSpaceCharacter(currentInputCharacter) {
				if *t.currentToken.asTag().TagName() == "title" || *t.currentToken.asTag().TagName() == "textarea" {
					t.state = HTMLState_BeforeAttributeName
					goto start
				}
			}
			if currentInputCharacter == '/' {
				if *t.currentToken.asTag().TagName() == "title" || *t.currentToken.asTag().TagName() == "textarea" {
					t.state = HTMLState_SelfClosingStartTag
					goto start
				}
			}
			if currentInputCharacter == '>' {
				if *t.currentToken.asTag().TagName() == "title" || *t.currentToken.asTag().TagName() == "textarea" {
					t.emitCurrentToken()
					t.state = HTMLState_Data
					goto start
				}
			}
			if isASCIIUpper(currentInputCharacter) {
				*t.currentToken.asTag().TagName() += string(toASCIILower(currentInputCharacter))
				t.buffer += string(currentInputCharacter)
				goto start
			}
			if isASCIILower(currentInputCharacter) {
				*t.currentToken.asTag().TagName() += string(currentInputCharacter)
				t.buffer += string(currentInputCharacter)
				goto start
			}
			t.emitCharacter('<')
			t.emitCharacter('/')
			t.emitCharacters(&t.buffer)
			t.reconsumeIn(HTMLState_RCDATA)
		case HTMLState_RAWTEXTLessThanSign:
			if currentInputCharacter == '/' {
				t.buffer = ""
				t.state = HTMLState_RAWTEXTEndTagOpen
				goto start
			}
			t.emitCharacter('<')
			t.reconsumeIn(HTMLState_RAWTEXT)
		case HTMLState_RAWTEXTEndTagOpen:
			if isASCIIAlpha(currentInputCharacter) {
				//FIXME: Create a new end tag token, set its tag name to the empty string
				t.reconsumeIn(HTMLState_RAWTEXTEndTagName)
				goto start
			}
			t.emitCharacter('<')
			t.emitCharacter('/')
			t.reconsumeIn(HTMLState_RAWTEXT)
		case HTMLState_RAWTEXTEndTagName:
			if isCharacterTabulation(currentInputCharacter) || isLineFeed(currentInputCharacter) || isFormFeed(currentInputCharacter) || isSpaceCharacter(currentInputCharacter) {
				if strings.Contains("style xmp iframe noembed noframes", *t.currentToken.asTag().TagName()) { //FIXME: This is a hack
					t.state = HTMLState_BeforeAttributeName
					goto start
				}
			}
			if currentInputCharacter == '/' {
				if strings.Contains("style xmp iframe noembed noframes", *t.currentToken.asTag().TagName()) { //FIXME: This is a hack
					t.state = HTMLState_SelfClosingStartTag
					goto start
				}
			}
			if currentInputCharacter == '>' {
				if strings.Contains("style xmp iframe noembed noframes", *t.currentToken.asTag().TagName()) { //FIXME: This is a hack
					t.emitCurrentToken()
					t.state = HTMLState_Data
					goto start
				}
			}
			if isASCIIUpper(currentInputCharacter) {
				*t.currentToken.asTag().TagName() += string(toASCIILower(currentInputCharacter))
				t.buffer += string(currentInputCharacter)
				goto start
			}
			if isASCIILower(currentInputCharacter) {
				*t.currentToken.asTag().TagName() += string(currentInputCharacter)
				t.buffer += string(currentInputCharacter)
				goto start
			}
			t.emitCharacter('<')
			t.emitCharacter('/')
			t.emitCharacters(&t.buffer)
			t.reconsumeIn(HTMLState_RAWTEXT)
		case HTMLState_ScriptDataLessThanSign:
			if currentInputCharacter == '/' {
				t.buffer = ""
				t.state = HTMLState_ScriptDataEndTagOpen
				goto start
			}
			if currentInputCharacter == '!' {
				t.emitCharacter('<')
				t.emitCharacter('!')
				t.state = HTMLState_ScriptDataEscapeStart
				goto start
			}
			t.emitCharacter('<')
			t.reconsumeIn(HTMLState_ScriptData)
		case HTMLState_ScriptDataEndTagOpen:
			if isASCIIAlpha(currentInputCharacter) {
				//FIXME: Create a new end tag token, set its tag name to the empty string
				t.reconsumeIn(HTMLState_ScriptDataEndTagName)
				goto start
			}
			t.emitCharacter('<')
			t.emitCharacter('/')
			t.reconsumeIn(HTMLState_ScriptData)
		case HTMLState_ScriptDataEndTagName:
			if isCharacterTabulation(currentInputCharacter) || isLineFeed(currentInputCharacter) || isFormFeed(currentInputCharacter) || isSpaceCharacter(currentInputCharacter) {
				if *t.currentToken.asTag().TagName() == "script" {
					t.state = HTMLState_BeforeAttributeName
					goto start
				}
			}
			if currentInputCharacter == '/' {
				if *t.currentToken.asTag().TagName() == "script" {
					t.state = HTMLState_SelfClosingStartTag
					goto start
				}
			}
			if currentInputCharacter == '>' {
				if *t.currentToken.asTag().TagName() == "script" {
					t.emitCurrentToken()
					t.state = HTMLState_Data
					goto start
				}
			}
			if isASCIIUpper(currentInputCharacter) {
				*t.currentToken.asTag().TagName() += string(toASCIILower(currentInputCharacter))
				t.buffer += string(currentInputCharacter)
				goto start
			}
			if isASCIILower(currentInputCharacter) {
				*t.currentToken.asTag().TagName() += string(currentInputCharacter)
				t.buffer += string(currentInputCharacter)
				goto start
			}
			t.emitCharacter('<')
			t.emitCharacter('/')
			t.emitCharacters(&t.buffer)
			t.reconsumeIn(HTMLState_ScriptData)
		case HTMLState_ScriptDataEscapeStart:
			if currentInputCharacter == '-' {
				t.emitCharacter('-')
				t.state = HTMLState_ScriptDataEscapeStartDash
				goto start
			}
			t.reconsumeIn(HTMLState_ScriptData)
		case HTMLState_ScriptDataEscapeStartDash:
			if currentInputCharacter == '-' {
				t.emitCharacter('-')
				t.state = HTMLState_ScriptDataEscapedDashDash
				goto start
			}
			t.reconsumeIn(HTMLState_ScriptData)
		case HTMLState_ScriptDataEscaped:
			if currentInputCharacter == '-' {
				t.emitCharacter('-')
				t.state = HTMLState_ScriptDataEscapedDash
				goto start
			}
			if currentInputCharacter == '<' {
				t.state = HTMLState_ScriptDataEscapedLessThanSign
				goto start
			}
			if currentInputCharacter == 0 {
				t.parseError(HTMLParseError_UnexpectedNullCharacter)
				t.emitCharacter(0xFFFD)
				goto start
			}
			if currentInputCharacter == EOF {
				t.parseError(HTMLParseError_EofInScriptHtmlCommentLikeText)
				t.emitEOF()
				return
			}
			t.emitCharacter(currentInputCharacter)
			goto start
		case HTMLState_ScriptDataEscapedDash:
			if currentInputCharacter == '-' {
				t.emitCharacter('-')
				t.state = HTMLState_ScriptDataEscapedDashDash
				goto start
			}
			if currentInputCharacter == '<' {
				t.state = HTMLState_ScriptDataEscapedLessThanSign
				goto start
			}
			if currentInputCharacter == 0 {
				t.parseError(HTMLParseError_UnexpectedNullCharacter)
				t.emitCharacter(0xFFFD)
				t.state = HTMLState_ScriptDataEscaped
				goto start
			}
			if currentInputCharacter == EOF {
				t.parseError(HTMLParseError_EofInScriptHtmlCommentLikeText)
				t.emitEOF()
				return
			}
			t.emitCharacter(currentInputCharacter)
			t.state = HTMLState_ScriptDataEscaped
			goto start
		case HTMLState_ScriptDataEscapedDashDash:
			if currentInputCharacter == '-' {
				t.emitCharacter('-')
				goto start
			}
			if currentInputCharacter == '<' {
				t.state = HTMLState_ScriptDataEscapedLessThanSign
				goto start
			}
			if currentInputCharacter == '>' {
				t.emitCharacter('>')
				t.state = HTMLState_ScriptData
				goto start
			}
			if currentInputCharacter == 0 {
				t.parseError(HTMLParseError_UnexpectedNullCharacter)
				t.emitCharacter(0xFFFD)
				t.state = HTMLState_ScriptDataEscaped
				goto start
			}
			if currentInputCharacter == EOF {
				t.parseError(HTMLParseError_EofInScriptHtmlCommentLikeText)
				t.emitEOF()
				return
			}
			t.emitCharacter(currentInputCharacter)
			t.state = HTMLState_ScriptDataEscaped
			goto start
		case HTMLState_ScriptDataEscapedLessThanSign:
			if currentInputCharacter == '/' {
				t.buffer = ""
				t.state = HTMLState_ScriptDataEscapedEndTagOpen
				goto start
			}
			if isASCIIAlpha(currentInputCharacter) {
				t.buffer = ""
				t.emitCharacter('<')
				t.state = HTMLState_ScriptDataDoubleEscapeStart
				goto start
			}
			t.emitCharacter('<')
			t.reconsumeIn(HTMLState_ScriptDataEscaped)
		case HTMLState_ScriptDataEscapedEndTagOpen:
			if isASCIIAlpha(currentInputCharacter) {
				//FIXME: create a new end tag token
				t.state = HTMLState_ScriptDataEscapedEndTagName
				goto start
			}
			t.emitCharacter('<')
			t.emitCharacter('/')
			t.reconsumeIn(HTMLState_ScriptDataEscaped)
		case HTMLState_ScriptDataEscapedEndTagName:
			if isSpaceCharacter(currentInputCharacter) || isLineFeed(currentInputCharacter) || isFormFeed(currentInputCharacter) || isSpaceCharacter(currentInputCharacter) {
				if *t.currentToken.asTag().TagName() == "script" {
					t.state = HTMLState_BeforeAttributeName
					goto start
				}
			}
			if currentInputCharacter == '/' {
				if *t.currentToken.asTag().TagName() == "script" {
					t.state = HTMLState_SelfClosingStartTag
					goto start
				}
			}
			if currentInputCharacter == '>' {
				if *t.currentToken.asTag().TagName() == "script" {
					t.emitCurrentToken()
					t.state = HTMLState_Data
					goto start
				}
			}
			if isASCIIUpper(currentInputCharacter) {
				*t.currentToken.asTag().TagName() += string(toASCIILower(currentInputCharacter))
				t.buffer += string(currentInputCharacter)
				goto start
			}
			if isASCIILower(currentInputCharacter) {
				*t.currentToken.asTag().TagName() += string(currentInputCharacter)
				t.buffer += string(currentInputCharacter)
				goto start
			}
			t.emitCharacter('<')
			t.emitCharacter('/')
			t.emitCharacters(&t.buffer)
			t.reconsumeIn(HTMLState_ScriptDataEscaped)
		case HTMLState_ScriptDataDoubleEscapeStart:
			if isCharacterTabulation(currentInputCharacter) || isLineFeed(currentInputCharacter) || isFormFeed(currentInputCharacter) || isSpaceCharacter(currentInputCharacter) || currentInputCharacter == '/' || currentInputCharacter == '>' {
				if t.buffer == "script" {
					t.state = HTMLState_ScriptDataDoubleEscaped
				} else {
					t.state = HTMLState_ScriptDataEscaped
				}
				t.emitCharacter(currentInputCharacter)
				goto start
			}
			if isASCIIUpper(currentInputCharacter) {
				t.buffer += string(toASCIILower(currentInputCharacter))
				t.emitCharacter(currentInputCharacter)
				goto start
			}
			if isASCIILower(currentInputCharacter) {
				t.buffer += string(currentInputCharacter)
				t.emitCharacter(currentInputCharacter)
				goto start
			}
			t.reconsumeIn(HTMLState_ScriptDataEscaped)
		case HTMLState_ScriptDataDoubleEscaped:
			if currentInputCharacter == '-' {
				t.emitCharacter('-')
				t.state = HTMLState_ScriptDataDoubleEscapedDash
				goto start
			}
			if currentInputCharacter == '<' {
				t.emitCharacter('<')
				t.state = HTMLState_ScriptDataDoubleEscapedLessThanSign
				goto start
			}
			if currentInputCharacter == 0 {
				t.parseError(HTMLParseError_UnexpectedNullCharacter)
				t.emitCharacter(0xFFFD)
				goto start
			}
			if currentInputCharacter == EOF {
				t.parseError(HTMLParseError_EofInScriptHtmlCommentLikeText)
				t.emitEOF()
				return
			}
			t.emitCharacter(currentInputCharacter)
		case HTMLState_ScriptDataDoubleEscapedDash:
			if currentInputCharacter == '-' {
				t.emitCharacter('-')
				t.state = HTMLState_ScriptDataDoubleEscapedDashDash
				goto start
			}
			if currentInputCharacter == '<' {
				t.emitCharacter('<')
				t.state = HTMLState_ScriptDataDoubleEscapedLessThanSign
				goto start
			}
			if currentInputCharacter == 0 {
				t.parseError(HTMLParseError_UnexpectedNullCharacter)
				t.state = HTMLState_ScriptDataDoubleEscaped
				t.emitCharacter(0xFFFD)
				goto start
			}
			if currentInputCharacter == EOF {
				t.parseError(HTMLParseError_EofInScriptHtmlCommentLikeText)
				t.emitEOF()
				return
			}
			t.emitCharacter(currentInputCharacter)
			t.state = HTMLState_ScriptDataDoubleEscaped
		case HTMLState_ScriptDataDoubleEscapedDashDash:
		case HTMLState_ScriptDataDoubleEscapedLessThanSign:
		case HTMLState_ScriptDataDoubleEscapeEnd:
		case HTMLState_BeforeAttributeName:
		case HTMLState_AttributeName:
		case HTMLState_AfterAttributeName:
		case HTMLState_BeforeAttributeValue:
		case HTMLState_AttributeValueDoubleQuoted:
		case HTMLState_AttributeValueSingleQuoted:
		case HTMLState_AttributeValueUnquoted:
		case HTMLState_AfterAttributeValueQuoted:
		case HTMLState_SelfClosingStartTag:
		case HTMLState_BogusComment:
		case HTMLState_MarkupDeclarationOpen:
		case HTMLState_CommentStart:
		case HTMLState_CommentStartDash:
		case HTMLState_Comment:
		case HTMLState_CommentLessThanSign:
		case HTMLState_CommentLessThanSignBang:
		case HTMLState_CommentLessThanSignBangDash:
		case HTMLState_CommentLessThanSignBangDashDash:
		case HTMLState_CommentEndDash:
		case HTMLState_CommentEnd:
		case HTMLState_CommentEndBang:
		case HTMLState_DOCTYPE:
		case HTMLState_BeforeDOCTYPEName:
		case HTMLState_DOCTYPEName:
		case HTMLState_AfterDOCTYPEName:
		case HTMLState_AfterDOCTYPEPublicKeyword:
		case HTMLState_BeforeDOCTYPEPublicIdentifier:
		case HTMLState_DOCTYPEPublicIdentifierDoubleQuoted:
		case HTMLState_DOCTYPEPublicIdentifierSingleQuoted:
		case HTMLState_AfterDOCTYPEPublicIdentifier:
		case HTMLState_BetweenDOCTYPEPublicAndSystemIdentifiers:
		case HTMLState_AfterDOCTYPESystemKeyword:
		case HTMLState_BeforeDOCTYPESystemIdentifier:
		case HTMLState_DOCTYPESystemIdentifierDoubleQuoted:
		case HTMLState_DOCTYPESystemIdentifierSingleQuoted:
		case HTMLState_AfterDOCTYPESystemIdentifier:
		case HTMLState_BogusDOCTYPE:
		case HTMLState_CDATASection:
		case HTMLState_CDATASectionBracket:
		case HTMLState_CDATASectionEnd:
		case HTMLState_CharacterReference:
		case HTMLState_NamedCharacterReference:
		case HTMLState_AmbiguousAmpersand:
		case HTMLState_NumericCharacterReference:
		case HTMLState_HexadecimalCharacterReferenceStart:
		case HTMLState_DecimalCharacterReferenceStart:
		case HTMLState_HexadecimalCharacterReference:
		case HTMLState_DecimalCharacterReference:
		case HTMLState_NumericCharacterReferenceEnd:
		}
	}
}

func isASCIILower(character int) bool {
	return character >= 'a' && character <= 'z'
}

func toASCIILower(character int) int {
	if !isASCIIAlpha(character) {
		//FIXME: Add error
		return character
	}
	if isASCIIUpper(character) {
		return character + 32
	}
	return character
}

func isASCIIUpper(character int) bool {
	return character >= 'A' && character <= 'Z'
}

func isSpaceCharacter(character int) bool {
	return character == ' '
}

func isFormFeed(character int) bool {
	return character == '\f'
}

func isLineFeed(character int) bool {
	return character == '\n'
}

func isCharacterTabulation(character int) bool {
	return character == '\t'
}

func isASCIIAlpha(character int) bool {
	return (character >= 'a' && character <= 'z') || (character >= 'A' && character <= 'Z')
}

func (t *tokenizer) nextCodePoint() int {
	rune, _, _ := t.input.ReadRune()
	return int(rune)
}

func (t *tokenizer) isEOF() bool {
	return false
}

func (t *tokenizer) emitCharacter(character int) {
	t.currentToken = HTMLTokenCharacter{
		codePoint: character,
	}
}

func (t *tokenizer) emitTokenStartTag() {
	t.currentToken = HTMLTokenStartTag{}
}

func (t *tokenizer) reconsumeIn(name HTMLState) {
	t.state = name
	t.input.UnreadRune()
}

func (t *tokenizer) emitEOF() {
	t.currentToken = HTMLTokenEOF{}
}

func (t *tokenizer) parseError(error HTMLParseError) {
	//TODO: Add option to ignore errors
	println("Parse error: " + error.String())
	//TODO: Add option to panic on errors
	//TODO: Better error handling
}

func (t *tokenizer) emitTokenEndTag(s string) {
	t.currentToken = HTMLTokenEndTag{
		tagName: s,
	}
}

func (t *tokenizer) emitCurrentToken() {
	t.tokens = append(t.tokens, t.currentToken)
}

func (t *tokenizer) emitCharacters(buffer *string) {
	for _, char := range *buffer {
		t.emitCharacter(int(char))
	}
	*buffer = ""
}
