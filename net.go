package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"strings"
)

func request(url string) (headers []string, body string, er error) {

	split := strings.SplitN(url, "://", 2)
	schema := split[0]

	split = strings.SplitN(split[1], "/", 2)
	host := split[0]
	path := "/" + split[1]
	_ = path

	if strings.Contains(host, ":") {
		er = fmt.Errorf("custom port not supported")
		return
	}

	var con net.Conn
	var err error
	if schema == "http" {
		con, err = net.Dial("tcp", host+":80")
	} else if schema == "https" {
		con, err = tls.Dial("tcp", host+":443", nil)
	} else {
		er = fmt.Errorf("schema \"" + schema + "\" not implemented")
		return
	}
	checkErr(err)

	defer con.Close()

	msg := "GET " + path + " HTTP/1.0\r\n" +
		"Host: " + host + "\r\n\r\n"

	_, err = con.Write([]byte(msg))
	checkErr(err)

	reply, err := ioutil.ReadAll(con)
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
