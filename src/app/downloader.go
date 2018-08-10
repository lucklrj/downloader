package main

import (
	"github.com/ddliu/go-httpclient"
	"strconv"
	"fmt"
	"os"
	"net/http"
	"strings"
	"math"
	"flag"
)

type httpData struct {
	Index    int
	Body     string
	Response *http.Response
	Err      error
}

var (
	downFile = flag.String("file", "", "要下载的文件地址")
	target   = flag.String("target", "", "存放位置")
	thread   = flag.String("thread", "10", "线程数")
)

func init() {
	flag.Parse()
	if *downFile == "" {
		fmt.Println("缺少文件下载地址")
		os.Exit(0)
	}
	if *target == "" {
		*target = "./"
	}
}
func main() {
	
	firstDataSize := 1024
	threadNum, _ := strconv.Atoi(*thread)
	url := *downFile
	fileName := url[strings.LastIndex(url, "/")+1:]
	
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
	singleRequestSize := math.Floor(float64(remainTotal) / float64(threadNum))
	
	c := make([]chan httpData, threadNum)
	
	startPoint := firstDataSize
	endPoint := startPoint + int(singleRequestSize)
	
	for i := 0; i < threadNum; i++ {
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
	for i := 0; i < threadNum; i++ {
		fileContent = fileContent + allData[i]
	}
	
	f, err := os.OpenFile(*target+fileName, os.O_RDWR|os.O_CREATE, 0755)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		f.WriteString(fileContent)
	}
	
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
