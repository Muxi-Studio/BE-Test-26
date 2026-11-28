package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	temp := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3OTU4NTg4NzQsImlhdCI6MTc5NTg1MTY3NCwiaXNzIjoiemVtdS1hdXRoIiwianRpIjoiYjg0NjZhYjRlMGIzYTlhNiIsInJvbGUiOiJlbXBsb3llZSIsInN1YiI6Im11eGlhb3hpMDQifQ.RdnOC_EzxuPynCIV_fIm2jV4nllw-soc-aakvdnALy0"
	temp2 := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3OTU4NTM1NDcsImlhdCI6MTc5NTg1MTc0NywiaXNzIjoiemVtdS1hdXRoIiwic2NvcGUiOlsiYXVkaXQ6cmVhZCJdLCJzdWIiOiJtdXhpYW94aTA0IiwidGlja2V0IjoiVC02ZTk2MWUifQ.mm2Qz1yYEaHI7mPMk9Xz7DCAcEnG_nWShI3g7kfm_X0"
	req, _ := http.NewRequest("GET", "http://1.12.252.114:8080/logs", nil)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+temp)
	req.Header.Set("X-Elevation", temp2)

	resp, _ := http.DefaultClient.Do(req)

	resbody, _ := io.ReadAll(resp.Body)
	fmt.Println(string(resbody))
}
