package main

import (
	"testing"
)

func TestRequestExample(t *testing.T) {
	url := "http://example.org/"
	_, body, err := request(url)
	checkErr(err)

	correctBody := `<!doctype html>
<html>
<head>
    <title>Example Domain</title>

    <meta charset="utf-8" />
    <meta http-equiv="Content-type" content="text/html; charset=utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <style type="text/css">
    body {
        background-color: #f0f0f2;
        margin: 0;
        padding: 0;
        font-family: -apple-system, system-ui, BlinkMacSystemFont, "Segoe UI", "Open Sans", "Helvetica Neue", Helvetica, Arial, sans-serif;
        
    }
    div {
        width: 600px;
        margin: 5em auto;
        padding: 2em;
        background-color: #fdfdff;
        border-radius: 0.5em;
        box-shadow: 2px 3px 7px 2px rgba(0,0,0,0.02);
    }
    a:link, a:visited {
        color: #38488f;
        text-decoration: none;
    }
    @media (max-width: 700px) {
        div {
            margin: 0 auto;
            width: auto;
        }
    }
    </style>    
</head>

<body>
<div>
    <h1>Example Domain</h1>
    <p>This domain is for use in illustrative examples in documents. You may use this
    domain in literature without prior coordination or asking for permission.</p>
    <p><a href="https://www.iana.org/domains/example">More information...</a></p>
</div>
</body>
</html>
`

	if body != correctBody {
		t.Fatal("FAIL: \"http://example.org/\"")
		//t.Fatalf(`Request("http://example.org/") = _, %q, \n --------------- \n want %q`, body, correctBody)
	}
}

func TestRequestHTTPS(t *testing.T) {
	url := "https://example.org/"
	_, body, err := request(url)
	checkErr(err)

	correctBody := `<!doctype html>
<html>
<head>
    <title>Example Domain</title>

    <meta charset="utf-8" />
    <meta http-equiv="Content-type" content="text/html; charset=utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <style type="text/css">
    body {
        background-color: #f0f0f2;
        margin: 0;
        padding: 0;
        font-family: -apple-system, system-ui, BlinkMacSystemFont, "Segoe UI", "Open Sans", "Helvetica Neue", Helvetica, Arial, sans-serif;
        
    }
    div {
        width: 600px;
        margin: 5em auto;
        padding: 2em;
        background-color: #fdfdff;
        border-radius: 0.5em;
        box-shadow: 2px 3px 7px 2px rgba(0,0,0,0.02);
    }
    a:link, a:visited {
        color: #38488f;
        text-decoration: none;
    }
    @media (max-width: 700px) {
        div {
            margin: 0 auto;
            width: auto;
        }
    }
    </style>    
</head>

<body>
<div>
    <h1>Example Domain</h1>
    <p>This domain is for use in illustrative examples in documents. You may use this
    domain in literature without prior coordination or asking for permission.</p>
    <p><a href="https://www.iana.org/domains/example">More information...</a></p>
</div>
</body>
</html>
`

	if body != correctBody {
		t.Fatal("FAIL: \"https://example.org\"")
	}
}

func TestRequestPptiede(t *testing.T) {
	url := "https://pptie.de/test.html"
	_, body, err := request(url)
	checkErr(err)

	correctBody := "test\n"

	if body != correctBody {
		t.Fatal("FAIL: \"https://pptie.de/test.html\" with body: \n" + body)
	}
}

func TestRequestPortErr(t *testing.T) {
	url := "http://pptie.de:3001/"
	_, _, err := request(url)
	if err.Error() != "custom port not supported" {
		t.Fatalf("custom ports should not be supported")
	}
}

func TestRequestSchemeErr(t *testing.T) {
	url := "gbrk://pptie.de/"
	_, _, err := request(url)

	if err.Error() != "schema \"gbrk\" not implemented" {
		t.Fatalf("Doesn't throw schema unimplemented error")
	}
}
