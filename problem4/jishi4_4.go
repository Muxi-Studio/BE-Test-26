package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

func main() {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	// 写入普通文本字段
	for k, v := range map[string]string{
		"title":  "title",
		"reason": "reasonreason",
		"scope":  "audit:read",
		"alice":  "muxiaoxi09",
	} {
		if err := writer.WriteField(k, v); err != nil {
			fmt.Println(err)
		}
	}
	// 创建文件字段（file part）并写入内容
	filePart, err := writer.CreateFormFile("attachment", "notes.txt")
	if err != nil {
		fmt.Println(err)
	}
	if _, err := filePart.Write([]byte("hello multipart world")); err != nil {
		fmt.Println(err)
	}
	// 关键：必须在构造请求前 Close()，否则缺少结尾的 --boundary--
	if err := writer.Close(); err != nil {
		fmt.Println(err)
	}

	request, err := http.NewRequest(
		"POST",
		"http://1.12.252.114:8080/ticket",
		&body,
	)
	if err != nil {
		panic(err)
	}

	request.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3OTU4NTAxMDIsImlhdCI6MTc5NTg0MjkwMiwiaXNzIjoiemVtdS1hdXRoIiwianRpIjoiMmJlZmUyMDBjZGQ1NmY1NCIsInJvbGUiOiJlbXBsb3llZSIsInN1YiI6Im11eGlhb3hpMDkifQ.eyMBnNDUEd-PU0R2dt18GVeyIFjQDvBFqk869adUSMg")
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("X-Demo-Client", "go-tutorial")

	client := &http.Client{} // 创建一个 Client
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	value := response.Header.Get("Content-Type") // 获取响应头
	ibody, err := io.ReadAll(response.Body)      // 读取响应体
	if err != nil {
		panic(err)

	}

	fmt.Println(string(value))
	fmt.Println(string(ibody))
}
