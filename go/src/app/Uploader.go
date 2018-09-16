package main

import (
	"bytes"
	"fmt"
	"github.com/ddliu/go-httpclient"
	"io/ioutil"
)

func main() {
	url := "http://127.0.0.1:8088/api/upload"
	
	byte, _ := ioutil.ReadFile("1.txt")
	headers := make(map[string]string)
	
	headers["Host"] = "127.0.0.1:8088"
	headers["Content-Length"] = "18"
	headers["X-Content-Range"] = "bytes 0-17/51"
	headers["Content-Disposition"] = "attachment; filename=\"1111.txt\""
	headers["Session-ID"] = "lucklrj"
	headers["Content-Type"] = "application/octet-stream"
	
	rep, _ := httpclient.Do("POST", url, headers, bytes.NewReader(byte[0:18]))
	fmt.Println(rep.StatusCode)
	r, _ := ioutil.ReadAll(rep.Body)
	fmt.Println(string(r))
	
	//////////
	headers["Content-Length"] = "31"
	headers["X-Content-Range"] = "bytes 18-50/51"
	
	rep, _ = httpclient.Do("POST", url, headers, bytes.NewReader(byte[18:]))
	fmt.Println(rep.StatusCode)
	r, _ = ioutil.ReadAll(rep.Body)
	fmt.Println(string(r))
	
}
