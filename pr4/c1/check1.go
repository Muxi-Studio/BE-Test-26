package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var requestBody = map[string]string{
	"username": "muxiaoxi04",
	"password": "123456",
}

func main() {
	index, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "http://1.12.252.114:8080/login", bytes.NewReader(index))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, _ := http.DefaultClient.Do(req)

	resbody, _ := io.ReadAll(resp.Body)
	fmt.Println(string(resbody))

}
