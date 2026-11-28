package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {

	request, err := http.NewRequest(
		"GET",
		"http://1.12.252.114:8080/logs",
		nil,
	)
	if err != nil {
		panic(err)
	}

	request.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3OTU4NTAxMDIsImlhdCI6MTc5NTg0MjkwMiwiaXNzIjoiemVtdS1hdXRoIiwianRpIjoiMmJlZmUyMDBjZGQ1NmY1NCIsInJvbGUiOiJlbXBsb3llZSIsInN1YiI6Im11eGlhb3hpMDkifQ.eyMBnNDUEd-PU0R2dt18GVeyIFjQDvBFqk869adUSMg")
	request.Header.Add("X-Elevation", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3OTU4NDY4NzIsImlhdCI6MTc5NTg0NTA3MiwiaXNzIjoiemVtdS1hdXRoIiwic2NvcGUiOlsiYXVkaXQ6cmVhZCJdLCJzdWIiOiJtdXhpYW94aTA5IiwidGlja2V0IjoiVC0zMDJkMGYifQ.c3yRvwgJhiH1AZyGHO7MCIKPiOU-fWzIjZCK8D9M9lc")

	client := &http.Client{} // 创建一个 Client
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	value := response.Header.Get("Content-Type") // 获取响应头
	body, err := io.ReadAll(response.Body)       // 读取响应体
	if err != nil {
		panic(err)

	}

	fmt.Println(string(value))
	fmt.Println(string(body))
}
