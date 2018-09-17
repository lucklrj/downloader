package main

import (
	"bytes"
	"fmt"
	"github.com/ddliu/go-httpclient"
	"io/ioutil"
	
	//"io/ioutil"
	"os"
	"strconv"
)

func main() {
	
	url := "http://127.0.0.1:8088/api/upload"
	host := "127.0.0.1:8088"
	filePath := "/Users/mxj/Desktop/1.txt"
	sessionId := "lucklrj"
	domain := "lrj-test"
	
	uploadPerSize := uint64(3) //1M
	
	fileInfo, _ := os.Stat(filePath)
	var fileSize uint64 = uint64(fileInfo.Size())
	fileContent, _ := ioutil.ReadFile(filePath)
	
	if uploadPerSize > fileSize {
		uploadPerSize = fileSize
	}
	
	var runNum uint64 = (fileSize-1)/uploadPerSize + 1
	fmt.Println("上传次数：", runNum)
	fmt.Println("文件大小：", fileSize)
	
	var runIndex uint64 = 1
	
	var start uint64 = 0
	var end uint64 = uploadPerSize - 1
	
	var lengthInheader string = "0"
	var fileSizeInheader string = strconv.Itoa(int(fileSize))
	var length uint64 = 0
	headers := make(map[string]string)
	
	for runIndex <= runNum {
		
		fmt.Println("第", runIndex, "次上传")
		if runIndex < runNum {
			length = uploadPerSize
		} else {
			length = fileSize - (runIndex-1)*uploadPerSize
		}
		
		startInheader := strconv.Itoa(int(start))
		endInheader := strconv.Itoa(int(end))
		
		if runIndex == 1 {
			headers["Host"] = host
			headers["Content-Length"] = lengthInheader
			headers["X-Content-Range"] = "bytes " + startInheader + "-" + endInheader + "/" + fileSizeInheader
			headers["Content-Disposition"] = "attachment; filename=\"1111.txt\""
			headers["Session-ID"] = sessionId
			headers["Content-Type"] = "application/octet-stream"
			headers["Domain"] = domain
		} else {
			headers["Content-Length"] = lengthInheader
			headers["X-Content-Range"] = " bytes " + startInheader + "-" + endInheader + "/" + fileSizeInheader
		}
		
		rep, _ := httpclient.Do("POST", url, headers, bytes.NewReader(fileContent[start:end+1]))
		r, _ := ioutil.ReadAll(rep.Body)
		if runIndex == runNum {
			fmt.Println(string(r))
		}
		
		start = end + 1;
		if runIndex == runNum-1 {
			end = fileSize - 1
		} else {
			end = start + length - 1
		}
		runIndex++
	}
}
