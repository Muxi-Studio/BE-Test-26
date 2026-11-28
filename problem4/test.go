package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func uploadMemo(targetURL string) error {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3OTU4NTMwNDcsImlhdCI6MTc5NTg0NTg0NywiaXNzIjoiemVtdS1hdXRoIiwianRpIjoiZDk3ZmM4NmYxYzRmZjFiMyIsInJvbGUiOiJlbXBsb3llZSIsInN1YiI6Im11eGlhb3hpMDgifQ.wqsYB-Bx471vjVqTPl9puW5ACRJpFYC5Js8UA9pjdcM"
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	// 写入普通文本字段
	for k, v := range map[string]string{
		"title":  "alice",
		"reason": "2026年Q1会议纪要备份",
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
	if _, err := filePart.Write([]byte("hello multipart world")); err != nil {
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
	req.Header.Add("Authorization", "Bearer "+token) // 按需添加认证或自定义请求头
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
func addFileFromDisk(writer *multipart.Writer, fieldName, filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	part, err := writer.CreateFormFile(fieldName, filepath.Base(filePath))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, f)
	return err
}
func main() {
	if err := uploadMemo("http://1.12.252.114:8080/ticket"); err != nil {
		log.Fatal(err)
	}
}
