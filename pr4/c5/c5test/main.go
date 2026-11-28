package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	temp := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3OTU4NTg4NzQsImlhdCI6MTc5NTg1MTY3NCwiaXNzIjoiemVtdS1hdXRoIiwianRpIjoiYjg0NjZhYjRlMGIzYTlhNiIsInJvbGUiOiJlbXBsb3llZSIsInN1YiI6Im11eGlhb3hpMDQifQ.RdnOC_EzxuPynCIV_fIm2jV4nllw-soc-aakvdnALy0"
	req, _ := http.NewRequest("GET", "http://1.12.252.114:8080/xiaohei", nil)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+temp)

	resp, _ := http.DefaultClient.Do(req)

	resbody, _ := io.ReadAll(resp.Body)
	fmt.Println(string(resbody))
}
