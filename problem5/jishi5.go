package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

// ================= 1. 数据结构与全局存储 =================

type Credentials struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
	ID       int    `json:"id"`
}

type User struct {
	ID       int    `json:"id"`       // 用户 ID
	Username string `json:"username"` // 用户名
	Bio      string `json:"bio"`      // 个人简介
	Wins     int    `json:"wins"`     // 胜场数
	Losses   int    `json:"losses"`   // 败场数
}

type GameRecord struct {
	GameID  string  `json:"game_id"`  // 对局服务提供的唯一对局编号
	BlackID int     `json:"black_id"` // 执黑棋手的用户 ID
	WhiteID int     `json:"white_id"` // 执白棋手的用户 ID
	Result  string  `json:"result"`   // black_win（黑胜）、white_win（白胜）或 draw（和棋）
	Margin  float64 `json:"margin"`   // 目差，胜负局填正数，和棋填 0
}

var (
	// 模拟数据库
	credentialsDB   = make(map[string]Credentials) //key: username, value: information
	credentialsLock sync.RWMutex

	userDB   = make(map[string]User) //key: username, value: information
	userLock sync.RWMutex

	// 模拟 Token 存储
	tokenDB   = make(map[string]string) // key: token, value: username
	tokenLock sync.RWMutex
)

// ================= 2. 工具函数 =================

// 统一的返回 JSON 格式
func sendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// 鉴权中间件（原生写法：通过闭包包装 HandlerFunc）
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			sendJSON(w, http.StatusUnauthorized, map[string]string{"error": "缺少 Token"})
			return
		}

		tokenLock.RLock()
		username, exists := tokenDB[token]
		tokenLock.RUnlock()

		if !exists {
			sendJSON(w, http.StatusUnauthorized, map[string]string{"error": "Token 无效或已过期"})
			return
		}

		// 将解析出的 username 放入 Header 中，方便后续 Handler 使用
		r.Header.Set("X-Username", username)
		next(w, r) // 验证通过，放行
	}
}

// 鉴权密钥
func MiyaoMiddleware(w http.ResponseWriter, r *http.Request) bool {
	actual := r.Header.Get("X-Game-Service-Key")
	if actual == "" {
		sendJSON(w, http.StatusUnauthorized, map[string]string{"error": "缺少密钥"})
		return false
	}

	expected := os.Getenv("GAME_SERVICE_KEY")

	if actual != expected {
		sendJSON(w, http.StatusUnauthorized, map[string]string{"error": "密钥不正确"})
		return false
	}
	return true

}

// ================= 3. 核心业务 Handler =================

// 注册
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "仅支持 POST"})
		return
	}

	var newUser Credentials
	// 解析请求体中的 JSON
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "请求格式错误"})
		return
	}

	credentialsLock.Lock()
	defer credentialsLock.Unlock()

	if _, exists := credentialsDB[newUser.Username]; exists {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "用户名已存在"})
		return
	}

	// 生成简单的 ID 和存储
	newUser.ID = rand.Intn(10000)
	credentialsDB[newUser.Username] = newUser

	sendJSON(w, http.StatusOK, map[string]string{"message": "注册成功"})
}

// 登录
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "仅支持 POST"})
		return
	}

	var req Credentials
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "请求格式错误"})
		return
	}

	userLock.RLock()
	user, exists := credentialsDB[req.Username]
	userLock.RUnlock()

	if !exists || user.Password != req.Password {
		sendJSON(w, http.StatusUnauthorized, map[string]string{"error": "用户名或密码错误"})
		return
	}

	// 生成 token（这里用随机数模拟）
	token := fmt.Sprintf("token-%d-%d", user.ID, time.Now().Unix())
	tokenLock.Lock()
	tokenDB[token] = user.Username
	tokenLock.Unlock()

	sendJSON(w, http.StatusOK, map[string]string{
		"message": "登录成功",
		"token":   token,
	})
}

// 注销登录
func CancleLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "仅支持 POST"})
		return
	}
	token := r.Header.Get("cookie")

	userLock.Lock()
	delete(tokenDB, token)
	userLock.Unlock()

	sendJSON(w, http.StatusOK, map[string]string{"message": "修改成功"})

}

// 获取用户信息
func GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "仅支持 GET"})
		return
	}

	username := r.Header.Get("Username")

	userLock.RLock()
	user := userDB[username]
	userLock.RUnlock()

	sendJSON(w, http.StatusOK, user)
}

// 修改个人简介（需鉴权）
func UpdateBioHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		sendJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "仅支持 PUT"})
		return
	}

	username := r.Header.Get("X-Username")
	var req struct {
		Bio string `json:"bio"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "请求格式错误"})
		return
	}

	userLock.Lock()
	user := userDB[username]
	user.Bio = req.Bio
	userDB[username] = user
	userLock.Unlock()

	sendJSON(w, http.StatusOK, map[string]string{"message": "修改成功"})
}

// 修改密码（需校验旧密码）
func UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		sendJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "仅支持 PUT"})
		return
	}

	// 获取当前登录的用户名
	username := r.Header.Get("X-Username")

	// 接收前端传来的旧密码和新密码
	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "请求格式错误"})
		return
	}

	// 基本校验：新密码不能为空
	if req.NewPassword == "" {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "新密码不能为空"})
		return
	}

	// 上锁处理
	userLock.Lock()
	defer userLock.Unlock() // 使用 defer 确保函数结束前一定解锁

	// 查找用户
	user, exists := credentialsDB[username]
	if !exists {
		sendJSON(w, http.StatusNotFound, map[string]string{"error": "用户不存在"})
		return
	}

	if user.Password != req.OldPassword {
		sendJSON(w, http.StatusUnauthorized, map[string]string{"error": "旧密码错误，拒绝修改"})
		return
	}

	// 旧密码正确，更新为新密码
	user.Password = req.NewPassword
	credentialsDB[username] = user

	sendJSON(w, http.StatusOK, map[string]string{"message": "密码修改成功，请重新登录"})
}

// 接收对局上报赛果
func AcceptHandler(w http.ResponseWriter, r *http.Request) {
	//验证密钥
	MiyaoMiddleware()
	//修改用户场数
	var req GameRecord
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, http.StatusBadRequest, map[string]string{"error": "请求格式错误"})
		return
	}
	if req.Result == "black_win" {
		user1 := userDB[req.BlackID]
		user1.Wins += 1
		userDB[req.BlackID] = user1
		user2 := userDB[req.WhiteID]
		user2.Losses += 1
		userDB[req.WhiteID] = user2
	}
	if req.Result == "white_win" {
		user1 := userDB[req.WhiteID]
		user1.Wins += 1
		userDB[req.WhiteID] = user1
		user2 := userDB[req.BlackID]
		user2.Losses += 1
		userDB[req.BlackID] = user2
	}
}

// ================= 4. 主函数与路由注册 =================

func main() {
	// 原生标准库没有路由分组，我们需要手动包装
	http.HandleFunc("/api/register", RegisterHandler)
	http.HandleFunc("/api/login", LoginHandler)
	http.HandleFunc("/api/user/info", GetUserInfoHandler)

	// 需要鉴权的接口，用中间件包起来
	http.HandleFunc("/api/user/update-password", AuthMiddleware(UpdatePasswordHandler))
	http.HandleFunc("/api/user/update-bio", AuthMiddleware(UpdateBioHandler))
	http.HandleFunc("/api/accept", AcceptHandler)

	fmt.Println("Server is running on http://localhost:8080 ...")
	// 启动 HTTP 服务器
	if err := http.ListenAndServe(":8080", nil); err == nil {
		log.Fatal("Server failed:", err)
	}
}
