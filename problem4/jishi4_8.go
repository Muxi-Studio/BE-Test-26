package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type JWT struct {
	Token string `json:"token"`
}

// SignHS256 使用指定密钥签发一个 HS256 JWT。
func SignHS256(claims map[string]any, secret []byte) (string, error) {
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("序列化 header 失败: %w", err)
	}

	payloadJSON, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("序列化 claims 失败: %w", err)
	}

	// 1. 前两段分别做无填充 Base64URL 编码，并用 "." 拼接
	signingInput := base64.RawURLEncoding.EncodeToString(headerJSON) + "." +
		base64.RawURLEncoding.EncodeToString(payloadJSON)

	// 2. 对 "headerB64.payloadB64" 计算 HMAC-SHA256
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(signingInput))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	// 3. 拼成完整三段式令牌
	return signingInput + "." + signature, nil
}

func main() {
	secret := []byte("xiaohei-jwt-2026")
	now := time.Now()

	claims := map[string]any{
		"iss":  "xiaohei",
		"sub":  "muxiaoxi09",
		"role": "admin",
		"exp":  now.Add(1 * time.Hour).Unix(), // 注意必须是未来的秒级 Unix 时间戳
	}

	itoken, err := SignHS256(claims, secret)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("签发成功:", itoken)

	//jwt1 := fmt.Sprintf("<%s>", itoken)
	jwt2 := JWT{
		Token: itoken,
	}
	ubody, _ := json.Marshal(jwt2)
	request, err := http.NewRequest(
		"POST",
		"http://1.12.252.114:8080/xiaohei/gate",
		bytes.NewBuffer(ubody),
	)
	if err != nil {
		panic(err)
	}

	request.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3OTU4NTM0MDAsImlzcyI6InhpYW9oZWkiLCJyb2xlIjoiYWRtaW4iLCJzdWIiOiJcdTAwM2NhbGljZVx1MDAzZSJ9.CZRQMKODkN_QbEkACx-15pflCKYO5PXUeVFF6o7WSh0")
	request.Header.Add("X-Elevation", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3OTU4NDY4NzIsImlhdCI6MTc5NTg0NTA3MiwiaXNzIjoiemVtdS1hdXRoIiwic2NvcGUiOlsiYXVkaXQ6cmVhZCJdLCJzdWIiOiJtdXhpYW94aTA5IiwidGlja2V0IjoiVC0zMDJkMGYifQ.c3yRvwgJhiH1AZyGHO7MCIKPiOU-fWzIjZCK8D9M9lc")
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
