package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

/*前面的忘记截图了😭
 *不过也没写多少其实
 */

type Resquest struct {
	Username string `json:"username`
	Password string `json:"password`
}

const (
	jwt = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3OTU4NTk0MzQsImlhdCI6MTc5NTg1MjIzNCwiaXNzIjoiemVtdS1hdXRoIiwianRpIjoiNmM2MDQ5NzQyYTMwMzY0MCIsInJvbGUiOiJlbXBsb3llZSIsInN1YiI6Im11eGlhb3hpMTMifQ.G1pQbBne8GLl7O3Dtw4vxb_MZ5Dp8SFcXoEtV1jYygU"
)

func main() {
	/*client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest(
		"GET",
		"http://1.12.252.114:8080/apikey",
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("Authorization", jwt)
	req.Header.Set("Content-Type", "application/json")

	resp, _ := client.Do(req)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	fmt.Println(body)*/

	//什么叫我的key被大肥鱼拿走了😡

	/*client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest(
		"GET",
		"http://1.12.252.114:8080/logs",
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, _ := client.Do(req)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	fmt.Println(body)*/

	uploadMemo("http://1.12.252.114:8080/ticket")

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

func uploadMemo(targetURL string) error {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	filepath := "/aaa"

	// 写入普通文本字段
	for k, v := range map[string]string{
		"title":      "aaaaaaaa",
		"reason":     "bbbbbbbbbbbbbbbbbb",
		"scope":      "audit:read",
		"attachment": filepath,
	} {
		if err := writer.WriteField(k, v); err != nil {
			return err
		}
	}

	// 创建文件字段（file part）并写入内容
	filePart, err := writer.CreateFormFile("file", "notes.txt")
	if err != nil {
		return err
	}
	if _, err := filePart.Write([]byte("我是muxiaoxi13")); err != nil {
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
	req.Header.Set("Authorization", jwt) // 按需添加认证或自定义请求头
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
