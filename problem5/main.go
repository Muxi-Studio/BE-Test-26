package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

// 用户账号。
type Credentials struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
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

// sessions 保存“登录令牌 -> 用户名”。
var sessions = map[string]string{}

// users 保存“用户名 -> 用户账号”。
var credentials = map[string]Credentials{}

var users = map[int]User{}

var games = map[int]GameRecord{}
var gsessions = map[string]int{}

// nextID 为下一条新用户分配编号。
var nextID = 1
var nxergame = 1

// mu 用来保护上面的 map。
var mu sync.RWMutex

func main() {
	mux := http.NewServeMux()

	// 注册和登录接口。
	mux.HandleFunc("POST /register", register)
	mux.HandleFunc("POST /login", login)

	// 用户的增删改查接口。请忽略所有的events，将其视为users
	mux.HandleFunc("POST /users", createuser)
	mux.HandleFunc("GET /users/{id}", geteuser)
	mux.HandleFunc("PUT /users/{id}", updateuser)
	mux.HandleFunc("DELETE /users/{id}", deleteuser)

	mux.HandleFunc("POST /games/", creategame)
	fmt.Println("服务器运行在 http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}

// register 接收 JSON 用户信息并注册账号。
func register(w http.ResponseWriter, r *http.Request) {
	var user Credentials
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "JSON 格式错误", http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Password == "" {
		http.Error(w, "用户名和密码不能为空", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if _, exists := credentials[user.Username]; exists {
		http.Error(w, "用户名已存在", http.StatusConflict)
		return
	}

	credentials[user.Username] = user
	fmt.Fprintln(w, "注册成功")
}

// login 验证账号密码，成功后设置 Cookie。
func login(w http.ResponseWriter, r *http.Request) {
	var input Credentials
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "JSON 格式错误", http.StatusBadRequest)
		return
	}

	mu.RLock()
	user, exists := credentials[input.Username]
	mu.RUnlock()

	if !exists || user.Password != input.Password {
		http.Error(w, "用户名或密码错误", http.StatusUnauthorized)
		return
	}

	// 生成一个随机 token，作为这次登录的凭证。
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		http.Error(w, "生成登录状态失败", http.StatusInternalServerError)
		return
	}
	token := hex.EncodeToString(randomBytes)

	// 服务器保存 token 对应的用户名。
	mu.Lock()
	sessions[token] = input.Username
	mu.Unlock()

	// 把 token 放入 Cookie 发给客户端。
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	})

	fmt.Fprintln(w, "登录成功")
}

// currentUser 根据 Cookie 查找当前登录用户。
func currentUser(r *http.Request) string {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return ""
	}

	mu.RLock()
	defer mu.RUnlock()
	return sessions[cookie.Value]
}

// createuser 创建用户，所属人取自登录用户。
func createuser(w http.ResponseWriter, r *http.Request) {
	username := currentUser(r)
	if username == "" {
		http.Error(w, "请先登录", http.StatusUnauthorized)
		return
	}

	var event User
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "JSON 格式错误", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// ID 和归属人由服务器填写，避免客户端伪造。
	event.ID = nextID
	event.Username = username

	users[nextID] = event
	nextID++

	json.NewEncoder(w).Encode(event)
}

// geteuser 根据 ID 查看待办，不限制是不是自己的。
func geteuser(w http.ResponseWriter, r *http.Request) {
	if currentUser(r) == "" {
		http.Error(w, "请先登录", http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID 格式错误", http.StatusBadRequest)
		return
	}

	mu.RLock()
	defer mu.RUnlock()

	event, exists := users[id]
	if !exists {
		http.Error(w, "用户不存在", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(event)
}

// updateuser 修改待办描述和期限，只允许归属人操作。
func updateuser(w http.ResponseWriter, r *http.Request) {
	username := currentUser(r)
	if username == "" {
		http.Error(w, "请先登录", http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID 格式错误", http.StatusBadRequest)
		return
	}

	var input User
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "JSON 格式错误", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	event, exists := users[id]
	if !exists {
		http.Error(w, "待办不存在", http.StatusNotFound)
		return
	}
	if event.Username != username {
		http.Error(w, "不能修改别人的待办", http.StatusForbidden)
		return
	}

	event.Bio = input.Bio
	users[id] = event

	json.NewEncoder(w).Encode(event)
}

// deleteuser 删除待办，只允许归属人操作。
func deleteuser(w http.ResponseWriter, r *http.Request) {
	username := currentUser(r)
	if username == "" {
		http.Error(w, "请先登录", http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID 格式错误", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	event, exists := users[id]
	if !exists {
		http.Error(w, "待办不存在", http.StatusNotFound)
		return
	}
	if event.Username != username {
		http.Error(w, "不能删除别人的用户", http.StatusForbidden)
		return
	}

	delete(users, id)
	fmt.Fprintln(w, "删除成功")
}

//创建对局没时间写了
/*func creategame(w http.ResponseWriter, r *http.Request) {
	var input GameRecord
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "JSON 格式错误", http.StatusBadRequest)
		return
	}

	expected := os.Getenv("GAME_SERVICE_KEY")    // 获取本地密钥
	actual := r.Header.Get("X-Game-Service-Key") // 获取请求携带的密钥
	if expected != actual {
		http.Error(w, "密钥错误", http.StatusBadRequest)
	}
	// 服务器保存 token 对应的用户名。
	mu.Lock()
	gsessions[expected] = input.GameID
	mu.Unlock()

	// 把 token 放入 Cookie 发给客户端。
	users[input.WhiteID].Wins

	fmt.Fprintln(w, "对局创立成功")
}
*/
