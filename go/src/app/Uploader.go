package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	url := "http://127.0.0.1:8088/api/upload"
	
	byte, _ := ioutil.ReadFile("1-1.txt")
	req, _ := http.NewRequest("POST", url, bytes.NewReader(byte))
	req.Header.Add("Host", "127.0.0.1:8088")
	req.Header.Add("Content-Length", "18")
	req.Header.Add("Content-Disposition", "attachment; filename=\"1111.txt\"")
	req.Header.Add("X-Content-Range", "bytes 0-17/51")
	req.Header.Add("Session-ID", "lucklrj")
	req.Header.Add("Content-Type", "application/octet-stream")
	
	client := &http.Client{}
	rep, _ := client.Do(req)
	fmt.Println(rep.StatusCode)
	r, _ := ioutil.ReadAll(rep.Body)
	fmt.Println(string(r))
	
	byte, _ = ioutil.ReadFile("1-2.txt")
	req, _ = http.NewRequest("POST", url, bytes.NewReader(byte))
	req.Header.Add("Host", "127.0.0.1:8088")
	req.Header.Add("Content-Length", "33")
	req.Header.Add("Content-Disposition", "attachment; filename=\"1111.txt\"")
	req.Header.Add("X-Content-Range", "bytes 18-50/51")
	req.Header.Add("Session-ID", "lucklrj")
	req.Header.Add("Content-Type", "application/octet-stream")
	
	rep, _ = client.Do(req)
	fmt.Println(rep.StatusCode)
	r, _ = ioutil.ReadAll(rep.Body)
	fmt.Println(string(r))
}
