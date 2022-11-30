package main

/*
import (
	"io"
	"log"
	"strings"
	"unicode"
)

type url struct {
	scheme   string
	path     string
	password string
	username string
}

const (
	schemeStartState                   = iota
	schemeState                        = iota
	noSchemeState                      = iota
	specialRelativeOrAuthorityState    = iota
	fileState                          = iota
	specialAuthoritySlashesState       = iota
	pathOrAuthorityState               = iota
	relativeState                      = iota
	relativeSlashState                 = iota
	opaquePathState                    = iota
	authorityState                     = iota
	pathState                          = iota
	specialAuthorityIgnoreSlashesState = iota
	hostState                          = iota
)

var (
	specialScheme = map[string]int{
		"ftp":   21,
		"file":  -1,
		"http":  80,
		"https": 443,
		"ws":    80,
		"wss":   443,
	}
)

// The basic URL parser takes a string input, with an optional null or base URL base (default null),
// an optional encoding encoding (default UTF-8), an optional URL url, and an optional state
// override state override, and then runs these steps:
func parseURL(i string) url {
	//Spec: https://url.spec.whatwg.org/#concept-basic-url-parser
	//	1. If url is not given:
	//		1. Set url to a new URL.
	url := url{}
	//		2. If input contains any leading or trailing C0 control or space, validation error.
	if i[0] == ' ' || i[len(i)-1] == ' ' { // won't add C0 for the moment
		log.Println("Validation error: Input contains leading or trailing space. ")
		//TODO: 3.Remove any leading and trailing C0 control or space from input
	}
	//	2. If input contains any ASCII tab or newline, validation error.
	if strings.ContainsAny(i, "\t\n") {
		log.Println("Validation error: Input contains ASCII tab or newline. ")
		//	3. Remove all ASCII tab or newline from input
		i = strings.ReplaceAll(i, "\t", "")
		i = strings.ReplaceAll(i, "\n", "")
	}
	//	4. Let state be state override if given, or scheme start state otherwise.
	state := schemeStartState
	//	5. Set encoding to the result of getting an output encoding from encoding.
	//skipped
	//	6. Let buffer be the empty string.
	buffer := make([]rune, 0)
	//	7. Let atSignSeen, insideBrackets, and passwordTokenSeen be false.
	atSignSeen, insideBrackets, passwordTokenSeen := false, false, false
	//	8. Let pointer be a pointer for input.
	reader := strings.NewReader(i)
	//	9. Keep running the following state machine by switching on state. If after a run pointer
	//		points to the EOF code point, go to the next step. Otherwise, increase pointer by 1 and
	//		continue with the state machine.
	for reader.Len() > 0 {
		c, _, err := reader.ReadRune()
		checkErr(err)
		switch state {
		case schemeStartState:
			//	1. If c is an ASCII alpha, append c, lowercased, to buffer, and set state to scheme state.
			if unicode.IsLetter(c) {
				buffer = append(buffer, unicode.ToLower(c))
				state = schemeState
				break
			}
			//	2. Otherwise, if state override is not given, set state to no scheme state and decrease pointer by 1.
			state = noSchemeState
			reader.UnreadRune()
			break
			//	3. Otherwise, validation error, return failure.
			log.Panic("Validation error: Invalid character in schema")

		case schemeState:
			//	1. If c is an ASCII alphanumeric, U+002B (+), U+002D (-), or U+002E (.), append c, lowercased, to buffer.
			if unicode.IsLetter(c) || c == '+' || c == '-' || c == '.' {
				buffer = append(buffer, unicode.ToLower(c))
				state = schemeState
				break
			}
			//	2. Otherwise, if c is U+003A (:), then:
			if c == ':' {
				//	1. If state override is given, then:
				//not implementing override
				//	2. Set url’s scheme to buffer.
				url.scheme = string(buffer)
				//	3. If state override is given, then:
				//not implementing override
				//	4. Set buffer to the empty string.
				buffer = make([]rune, 0)
				//	5. If url’s scheme is "file", then:
				if url.scheme == "file" {
					//	1. If remaining does not start with "//", validation error.
					if GetRemaining(reader)[:1] != "//" {
						log.Println("Validation error: nkgebksjgbjh") //TODO:
					}
					//	2. Set state to file state.
					state = fileState
				} else
				//	6. Otherwise, if url is special, base is non-null, and base’s scheme is url’s scheme:
				//TODO: we don't support non-null base url's
				//	7. Otherwise, if url is special, set state to special authority slashes state.
				if _, ok := specialScheme[url.scheme]; ok {
					state = specialAuthoritySlashesState

				} else
				//	8. Otherwise, if remaining starts with an U+002F (/), set state to path or authority state and increase pointer by 1.
				if GetRemaining(reader)[0] == '/' {
					state = pathOrAuthorityState
					reader.ReadRune()
				} else
				//	9. Otherwise, set url’s path to the empty string and set state to opaque path state.
				{
					url.path = ""
					state = opaquePathState
				}
			}
			//	3. Otherwise, if state override is not given, set buffer to the empty string, state to no scheme
			//		state, and start over (from the first code point in input).
			buffer = make([]rune, 0)
			state = noSchemeState
			reader.Reset(i)
			break //TODO: implement state override
			//	4. Otherwise, validation error, return failure
			log.Panic("If you see this look if your pc is on fire cause i will never have send you here... there just was a break the line before this so if ur coming here there is something seriusly wrong")

		case noSchemeState:
			//	1. If base is null, or base has an opaque path and c is not U+0023 (#), validation error, return failure.
			log.Panic("Validation error: Base not implemented")
			//	2. Otherwise, if base has an opaque path and c is U+0023 (#), set url’s scheme to base’s scheme, url’s path to base’s path, url’s query to base’s query, url’s fragment to the empty string, and set state to fragment state.
			//	3. Otherwise, if base’s scheme is not "file", set state to relative state and decrease pointer by 1.
			//	4. Otherwise, set state to file state and decrease pointer by 1.

		case specialRelativeOrAuthorityState:
			//TODO: implement stuff -> requires implementation of base
			log.Panic("Your pc is still on fire you should really call the fire deparment xD")

		case pathOrAuthorityState:
			//	1. If c is U+002F (/), then set state to authority state.
			if c == '/' {
				state = authorityState
				break
			}
			// 2. Otherwise, set state to path state, and decrease pointer by 1.
			state = pathState
			reader.UnreadRune()
			break

		case relativeState:
			//TODO: implement -> base
			log.Panic("It's to late now... you can't survive the fire anymore... RIP in peace")

		case relativeSlashState:
			log.Panic("How are you still alive? And how is you computer still running? It shoud have long burned to ashes.")

		case specialAuthoritySlashesState:
			//	1. If c is U+002F (/) and remaining starts with U+002F (/), then set state to special authority ignore slashes state and increase pointer by 1.
			if c == '/' && GetRemaining(reader)[0] == '/' {
				state = specialAuthorityIgnoreSlashesState
				reader.ReadRune()
				break
			}
			//	2. Otherwise, validation error, set state to special authority ignore slashes state and decrease pointer by 1.
			log.Println("Validation error: TODO")
			state = specialAuthorityIgnoreSlashesState
			reader.UnreadRune()
			break

		case specialAuthorityIgnoreSlashesState:
			//	1. If c is neither U+002F (/) nor U+005C (\), then set state to authority state and decrease pointer by 1.
			if c != '/' && c != '\\' {
				state = authorityState
				reader.UnreadRune()
				break
			}
			//	2. Otherwise, validation error.
			log.Println("Validation rrror: TODO")
			break

		case authorityState:
			//	1. If c is U+0040 (@), then:
			if c == '@' {
				//	1. Validation error.
				log.Println("Validation error: TODO")
				//	2. If atSignSeen is true, then prepend "%40" to buffer.
				if atSignSeen {
					buffer = append([]rune("%40"), buffer...)
				}
				//	3. Set atSignSeen to true.
				atSignSeen = true
				//	4. For each codePoint in buffer:
				for _, codePoint := range buffer {
					//	1. If codePoint is U+003A (:) and passwordTokenSeen is false, then set passwordTokenSeen to true and continue.
					if codePoint == ':' && !passwordTokenSeen {
						passwordTokenSeen = true
						continue
					}
					//	2. Let encodedCodePoints be the result of running UTF-8 percent-encode codePoint using the userinfo percent-encode set.
					encodedCodePoints := string(codePoint) //TODO: add support for encoding
					//	3. If passwordTokenSeen is true, then append encodedCodePoints to url’s password.
					if passwordTokenSeen {
						url.password += encodedCodePoints
						continue
					}
					//	4. Otherwise, append encodedCodePoints to url’s username.
					url.username += encodedCodePoints
				}
				//	5. Set buffer to the empty string.
				buffer = make([]rune, 0)
				break
			}
			//	2. Otherwise, if one of the following is true:
			//		* c is the EOF code point, U+002F (/), U+003F (?), or U+0023 (#)
			//		* url is special and c is U+005C (\)
			if _, ok := specialScheme[url.scheme]; (c == '/' || c == '?' || c == '#') || (ok && c == '\\') {
				//	1. If atSignSeen is true and buffer is the empty string, validation error, return failure.
				if atSignSeen && len(buffer) == 0 {
					log.Panic("Validation error: TODO")
				}
				//	2. Decrease pointer by the number of code points in buffer plus one, set buffer to the empty string, and set state to host state.
				for i := 0; i < len(buffer)+1; i++ {
					reader.UnreadRune()
				}
				buffer = make([]rune, 0)
				state = hostState
				break
			}
			//	3. Otherwise, append c to buffer.
			buffer = append(buffer, c)
		}
	}

	return url
}

func GetRemaining(reader *strings.Reader) string {
	remaining, err := io.ReadAll(reader)
	checkErr(err)
	reader.Reset(string(remaining))
	return string(remaining)
}
*/
