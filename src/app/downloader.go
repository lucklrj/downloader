package main

import (
	"github.com/ddliu/go-httpclient"
	"strconv"
	"fmt"
	"os"
	"net/http"
	"strings"
	"math"
)

type httpData struct {
	Index    int
	Body     string
	Response *http.Response
	Err      error
}

func main() {
	firstDataSize := 1024
	progress := 10
	url := "http://pic-bucket.nosdn.127.net/photo/0003/2018-08-10/DOR5JCRU00AJ0003NOS.jpg"
	fileContent := ""
	fileSize := 0
	data := down(url, 0, firstDataSize, 0)
	if data.Err != nil {
		die(data.Err.Error())
	}
	
	if _, ok := data.Response.Header["Content-Range"]; ok {
		contentRange := strings.Split(data.Response.Header["Content-Range"][0], "/")
		fileSize, _ = strconv.Atoi(contentRange[1])
	} else {
		die("无法获取文件大小")
	}
	
	fileContent = data.Body
	
	//分包抓取
	remainTotal := fileSize - firstDataSize - 1
	singleRequestSize := math.Floor(float64(remainTotal) / float64(progress))
	
	c := make([]chan httpData, progress)
	
	startPoint := firstDataSize
	endPoint := startPoint + int(singleRequestSize)
	
	for i := 0; i < progress; i++ {
		c[i] = make(chan httpData)
		go func(link chan httpData, url string, startPoint int, endPoint int, index int) {
			line := down(url, startPoint, endPoint, index)
			link <- line
		}(c[i], url, startPoint, endPoint, i+1)
		
		startPoint = endPoint + 1
		endPoint = startPoint + int(singleRequestSize)
		if endPoint > fileSize {
			endPoint = fileSize
		}
	}
	
	//获取内容
	allData := make(map[int]string)
	singleData := httpData{}
	for _, ch := range c {
		singleData = <-ch
		allData[singleData.Index] = singleData.Body
	}
	
	//拼接内容
	for i := 0; i < progress; i++ {
		fileContent = fileContent + allData[i]
	}
	fmt.Println(fileContent)
	
}

func down(url string, startPoint int, endPoint int, index int) httpData {
	res, _ := httpclient.
		Begin().
		WithHeader("Range", "bytes="+strconv.Itoa(startPoint)+"-"+strconv.Itoa(endPoint)).
		Get(url)
	
	bodyBytes, err := res.ReadAll()
	return httpData{Index: index, Body: string(bodyBytes), Response: res.Response, Err: err}
	
}

func die(msg string) {
	fmt.Println(msg)
	os.Exit(0)
}
