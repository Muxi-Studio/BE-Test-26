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
	"os"
	"time"
)

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
	signingInput := base64.RawURLEncoding.RawURLEncoding(headerJSON) + "." +
		base64.RawURLEncoding.RawURLEncoding(payloadJSON)

	// 2. 对 "headerB64.payloadB64" 计算 HMAC-SHA256
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(signingInput))
	signature := base64.RawURLEncoding.RawURLEncoding(mac.Sum(nil))

	// 3. 拼成完整三段式令牌
	return signingInput + "." + signature, nil
}

func main() {

	var secret = []byte(os.Getenv("JWT_SECRET"))
	now := time.Now()

	claims := map[string]any{
		"iss":  "xiaohei",
		"sub":  "muxiaoxi04",
		"role": "admin",
		"exp":  now.Add(1 * time.Hour).Unix(), // 注意必须是未来的秒级 Unix 时间戳
	}

	token, err := SignHS256(claims, secret)
	if err != nil {
		log.Fatal(err)
	}

	var requestBody = map[string]string{
		"token": token,
	}

	index, _ := json.Marshal(requestBody)
	temp := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3OTU4NTg4NzQsImlhdCI6MTc5NTg1MTY3NCwiaXNzIjoiemVtdS1hdXRoIiwianRpIjoiYjg0NjZhYjRlMGIzYTlhNiIsInJvbGUiOiJlbXBsb3llZSIsInN1YiI6Im11eGlhb3hpMDQifQ.RdnOC_EzxuPynCIV_fIm2jV4nllw-soc-aakvdnALy0"
	req, _ := http.NewRequest("POST", "http://1.12.252.114:8080/xiaohei/gate", bytes.NewReader(index))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+temp)

	resp, _ := http.DefaultClient.Do(req)

	resbody, _ := io.ReadAll(resp.Body)
	fmt.Println(string(resbody))

}
