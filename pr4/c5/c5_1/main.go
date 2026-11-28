package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	temp := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3OTU4NDkyMTMsImlhdCI6MTc5NTg0MjAxMywiaXNzIjoiemVtdS1hdXRoIiwianRpIjoiZDdkNGVjOTg3NzE4NmRlZSIsInJvbGUiOiJlbXBsb3llZSIsInN1YiI6Im11eGlhb3hpMDQifQ.VZFRM5t-IXqKySMxQXZC2Y7iBvxY1mwB3odKUMekE3U"
	req, _ := http.NewRequest("GET", "http://1.12.252.114:8080/xiaohei", nil)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+temp)

	resp, _ := http.DefaultClient.Do(req)

	resbody, _ := io.ReadAll(resp.Body)
	fmt.Println(string(resbody))
}
