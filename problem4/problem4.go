package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	//JWT := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3OTU4NTM0NjgsImlzcyI6InhpYW9oZWkiLCJyb2xlIjoiYWRtaW4iLCJzdWIiOiJtdXhpYW94aTA4In0.wJ7sRzqghVPcXAwp8Ke4DybOUUTeG6onUS-Lutmh4hk"
	token1 := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3OTU4NTA5ODcsImlhdCI6MTc5NTg0OTE4NywiaXNzIjoiemVtdS1hdXRoIiwic2NvcGUiOlsiYXVkaXQ6cmVhZCJdLCJzdWIiOiJtdXhpYW94aTA4IiwidGlja2V0IjoiVC1mODM2MDMifQ.P9HVXMXzNdQwyBHVRGzKhhxn33voAtvyXeDGLwju5W8"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3OTU4NTMwNDcsImlhdCI6MTc5NTg0NTg0NywiaXNzIjoiemVtdS1hdXRoIiwianRpIjoiZDk3ZmM4NmYxYzRmZjFiMyIsInJvbGUiOiJlbXBsb3llZSIsInN1YiI6Im11eGlhb3hpMDgifQ.wqsYB-Bx471vjVqTPl9puW5ACRJpFYC5Js8UA9pjdcM"
	url := "http://1.12.252.114:8080/xiaohei/gate"
	method := "POST"

	payload := strings.NewReader(`{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3OTU4NTYwNDgsImlzcyI6InhpYW9oZWkiLCJyb2xlIjoiYWRtaW4iLCJzdWIiOiJtdXhpYW94aTA4In0.94zf_L0jIEYrKYIkQZ1ET6dyKGMmAOpWUHhdDOquVc4"}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Elevation", token1)
	req.Header.Add("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

/*
心服口服。拿走你的 apikey 吧。
apikey：sk-zemu-e73263efefd71c3fdec8837d4573b702
账号：muxiaoxi08
*/
