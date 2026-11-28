// 账号：muxiaoxi01 密码：123456
// http://1.12.252.114:8080/start
package main

import (
	// 	"fmt"
	 	"net/http"
	     //"time"
	 	"io"
	"bytes"
	"fmt"
		//"encoding/json"
	"mime/multipart"
	"log"
)

  type User struct{
  	Username  string `json:"username"`
      Password  string `json:"password"`
	  }
// func main() {
// 	client := &http.Client{
//     Timeout : 30 * time.Second,
// 	}
// 	user := User{
// 		Username: "muxiaoxi01",
// 		Password: "123456",
// 	}
// 	user1,err := json.Marshal(user)
// 	rep,err := http.NewRequest(
// 		http.MethodPost,
// 		"http://1.12.252.114:8080/login",
// 		bytes.NewReader(user1),
// 	)
//     if err != nil {
// 		fmt.Println("请求体创建失败")
// 		return
// 	}
// 	rep.Header.Set("Content-Type","application/json")

// 	resp,err := client.Do(rep)
// 	if err != nil {
// 		fmt.Println("请求发送失败")
// 		return
// 	}
// 	body,err := io.ReadAll(resp.Body)
//     if err != nil {
// 		fmt.Println("读取失败")
// 		return
// 	}
// 	fmt.Println(string(body))
// 	fmt.Println(resp.StatusCode)
// 	fmt.Println(resp.Header)

//  }
//通行证：Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3OTU4MzgzMTgsImlhdCI6MTc5NTgzMTExOCwiaXNzIjoiemVtdS1hdXRoIiwianRpIjoiMWM5YzRmOTAwZTk1OWQwYiIsInJvbGUiOiJlbXBsb3llZSIsInN1YiI6Im11eGlhb3hpMDEifQ.G0y4I5RjQboVXlcpbhF9HvMd1RvD0gnf_4mM1Qiu_KE

// func main() {
// 	client := &http.Client{
//     Timeout : 30 * time.Second,
// 	}
// 	user := User{
// 		Username: "muxiaoxi01",
// 		Password: "123456",
// 	}
// 	user1,err := json.Marshal(user)
// 	rep,err := http.NewRequest(
// 		http.MethodGet,
// 		"http://1.12.252.114:8080/apikey",
// 		bytes.NewReader(user1),
// 	)
//     if err != nil {
// 		fmt.Println("请求体创建失败")
// 		return
// 	}
// 	rep.Header.Set("Content-Type","application/json")
// 	rep.Header.Set("Authorization","Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3OTU4NTM2ODIsImlhdCI6MTc5NTg0NjQ4MiwiaXNzIjoiemVtdS1hdXRoIiwianRpIjoiMDk2OWFiNTM3MGQ0YmFkMiIsInJvbGUiOiJlbXBsb3llZSIsInN1YiI6Im11eGlhb3hpMDEifQ.LktdfrsQjnNu1nn5vMlWTAcPCbY94-BL_MhuP2zd_G8")

// 	resp,err := client.Do(rep)
// 	if err != nil {
// 		fmt.Println("请求发送失败")
// 		return
// 	}
// 	body,err := io.ReadAll(resp.Body)
//     if err != nil {
// 		fmt.Println("读取失败")
// 		return
// 	}
// 	fmt.Println(string(body))
// 	fmt.Println(resp.StatusCode)
// 	fmt.Println(resp.Header)
// }
//需查看审计日志
//X-Request-Id:[req-69ccc895]]
// func main() {
// 	client := &http.Client{
//     Timeout : 30 * time.Second,
// 	}
// 	user := User{
// 		Username: "muxiaoxi01",
// 		Password: "123456",
// 	}
// 	user1,err := json.Marshal(user)
// 	rep,err := http.NewRequest(
// 		http.MethodPost,
// 		"http://1.12.252.114:8080/ticket",
// 		bytes.NewReader(user1),
// 	)
//     if err != nil {
// 		fmt.Println("请求体创建失败")
// 		return
// 	}
// 	rep.Header.Set("Content-Type","application/json")
// 	rep.Header.Set("Authorization","Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3OTU4NTM2ODIsImlhdCI6MTc5NTg0NjQ4MiwiaXNzIjoiemVtdS1hdXRoIiwianRpIjoiMDk2OWFiNTM3MGQ0YmFkMiIsInJvbGUiOiJlbXBsb3llZSIsInN1YiI6Im11eGlhb3hpMDEifQ.LktdfrsQjnNu1nn5vMlWTAcPCbY94-BL_MhuP2zd_G8")
// 	rep.Header.Set("Content-Type","multipart/form-data")
// 	resp,err := client.Do(rep)
// 	if err != nil {
// 		fmt.Println("请求发送失败")
// 		return
// 	}
// 	body,err := io.ReadAll(resp.Body)
//     if err != nil {
// 		fmt.Println("读取失败")
// 		return
// 	}
// 	fmt.Println(string(body))
// 	fmt.Println(resp.StatusCode)
// 	fmt.Println(resp.Header)

//  }
func uploadMemo(targetURL string) error {
	var body bytes.Buffer
        writer := multipart.NewWriter(&body)
		for k, v := range map[string]string{
                "title": "alice",
                "reason":"dsjskksdsjkwj",
				"scope" : "audit:read",
        } {
                if err := writer.WriteField(k, v); err != nil {
                        return err
                }
        }
		filePart, err := writer.CreateFormFile("file", "notes.txt")
        if err != nil {
                return err
        }
        if _, err := filePart.Write([]byte("hello multipart world")); err != nil {
                return err
        }
        if err := writer.Close(); err != nil {
                return err
        }
        req, err := http.NewRequest(http.MethodPost, targetURL, &body)
		
        // 必须使用 FormDataContentType() 自动带上 boundary
        req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization","Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3OTU4NTM2ODIsImlhdCI6MTc5NTg0NjQ4MiwiaXNzIjoiemVtdS1hdXRoIiwianRpIjoiMDk2OWFiNTM3MGQ0YmFkMiIsInJvbGUiOiJlbXBsb3llZSIsInN1YiI6Im11eGlhb3hpMDEifQ.LktdfrsQjnNu1nn5vMlWTAcPCbY94-BL_MhuP2zd_G8")
        req.Header.Set("Content-Disposition","form-data;name="file"; filename="notes.txt"")
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
        if err := uploadMemo("http://1.12.252.114:8080/ticket"); err != nil {
                log.Fatal(err)
        }
}
		
