package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func request(url string) (headers []string, body string, er error) {

	split := strings.SplitN(url, "://", 2)
	schema := split[0]

	splits := strings.SplitN(split[1], "/", 2)
	host := splits[0]

	if strings.Contains(host, ":") {
		er = fmt.Errorf("custom port not supported")
		return
	}

	if schema == "http" {
		con, err := net.Dial("tcp", host+":80")
		return http(url, con, err)
	} else if schema == "https" {
		con, err := tls.Dial("tcp", host+":443", nil)
		return http(url, con, err)
	} else if schema == "file" {
		//get the content of the file
		content, err := os.ReadFile(split[1])
		checkErr(err)
		return nil, string(content), nil
	} else {
		er = fmt.Errorf("schema \"" + schema + "\" not implemented") //TODO: add support for other schemas (view-source:, data:)
		return
	}
}

func http(url string, con net.Conn, err error) (headers []string, body string, er error) {
	checkErr(err)

	split := strings.SplitN(url, "://", 2)

	split = strings.SplitN(split[1], "/", 2)
	host := split[0]
	path := "/" + split[1]
	_ = path

	defer con.Close()

	msg := "GET " + path + " HTTP/1.0\r\n" +
		"Host: " + host + "\r\n\r\n"

	_, err = con.Write([]byte(msg))
	checkErr(err)

	reply, err := io.ReadAll(con)
	checkErr(err)

	splitReply := strings.Split(string(reply), "\r\n")
	//----------Split Header-Body------------
	headers = make([]string, 0)
	body = ""
	inHeaders := true
	for _, x := range splitReply {
		if x == "" {
			inHeaders = false
			continue
		}
		if inHeaders {
			headers = append(headers, x)
			continue
		}
		body = x
	}

	return
}
