package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {

	user := User{
		Username: "muxiaoxi09",
		Password: "123456",
	}

	ubody, _ := json.Marshal(user)

	request, err := http.NewRequest(
		"POST",
		"http://1.12.252.114:8080/login",
		bytes.NewBuffer(ubody),
	)
	if err != nil {
		panic(err)
	}

	//request.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3OTU4NDk3OTAsImlhdCI6MTc5NTg0MjU5MCwiaXNzIjoiemVtdS1hdXRoIiwianRpIjoiNTYwNjZkNGQ4ZjAxYWNhMSIsInJvbGUiOiJlbXBsb3llZSIsInN1YiI6Im11eGlhb3hpMDMifQ.bm-jXroazatWzEYtPzM_2CR5bhGIwNLgy6tjeE6Ie5s")
	request.Header.Add("Content-Type", " application/json")

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
