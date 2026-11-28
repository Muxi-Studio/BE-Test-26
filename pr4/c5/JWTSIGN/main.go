package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

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
	secret := []byte("JWT_SECRET")
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
	fmt.Println("签发成功:", token)
}
