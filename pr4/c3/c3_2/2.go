package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

func uploadMemo(targetURL string, temp string) error {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	// 写入普通文本字段
	for k, v := range map[string]string{
		"title":  "提权工单",
		"reason": "123456789012",
		"scope":  "audit:read",
	} {
		if err := writer.WriteField(k, v); err != nil {
			return err
		}
	}

	// 创建文件字段（file part）并写入内容
	filePart, err := writer.CreateFormFile("attachment", "notes.txt")
	if err != nil {
		return err
	}

	if _, err := filePart.Write([]byte("请输入文本")); err != nil {
		return err
	}

	// 关键：必须在构造请求前 Close()，否则缺少结尾的 --boundary--

	if err := writer.Close(); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, targetURL, &body)
	if err != nil {
		return err
	}
	// 必须使用 FormDataContentType() 自动带上 boundary
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-Demo-Client", "go-tutorial") // 按需添加认证或自定义请求头
	req.Header.Set("Authorization", "Bearer "+temp)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("状态码: %d\n响应体:\n%s\n", resp.StatusCode, string(respBytes))
	return nil
}

func main() {

	temp := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3OTU4NTg4NzQsImlhdCI6MTc5NTg1MTY3NCwiaXNzIjoiemVtdS1hdXRoIiwianRpIjoiYjg0NjZhYjRlMGIzYTlhNiIsInJvbGUiOiJlbXBsb3llZSIsInN1YiI6Im11eGlhb3hpMDQifQ.RdnOC_EzxuPynCIV_fIm2jV4nllw-soc-aakvdnALy0"

	uploadMemo("http://1.12.252.114:8080/ticket", temp)
	// req, _ := http.NewRequest("GET", "/ticket", nil)

	// req.Header.Set("Authorization", "Bearer "+temp)
	// req.Header.Set("Content-Type", "multipart/form-data")

	// resp, _ := http.DefaultClient.Do(req)

	// resbody, _ := io.ReadAll(resp.Body)
	// fmt.Println(string(resbody))
}
