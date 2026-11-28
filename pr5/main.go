package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Credentials struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}

var CreBank = make(map[string]Credentials)

var users = make(map[string]User)

var sessions = make(map[string]string)

var gameRecords = make(map[string]GameRecord)

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

// 注册
func register(w http.ResponseWriter, r *http.Request) {
	var requser Credentials
	json.NewDecoder(r.Body).Decode(&requser)
	CreBank[requser.Username] = requser
	fmt.Fprintln(w, "注册成功")
}

// 登录
func login(w http.ResponseWriter, r *http.Request) {
	var requser Credentials
	json.NewDecoder(r.Body).Decode(&requser)
	u, ok := CreBank[requser.Username]
	if !ok || u.Password != requser.Password {
		fmt.Fprintln(w, "用户名或密码错误")
		return
	}
	token := requser.Username + "123456"
	sessions[token] = requser.Username
	fmt.Fprintln(w, token)
}

// 获取信息
func getInfo(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	username := sessions[token]
	user := users[username]
	json.NewEncoder(w).Encode(user)
	fmt.Fprintln(w, user)
}

// 修改信息
func updateInfo(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	username := sessions[token]
	var date User
	json.NewDecoder(r.Body).Decode(&date)
	user := users[username]
	user = date
	users[username] = user
	fmt.Fprintln(w, "修改成功")
}

// 查询历史
func infoHistory(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	username := sessions[token]
	user := users[username]
	history := make([]GameRecord, 0)
	for _, rec := range gameRecords {
		if rec.BlackID == user.ID || rec.WhiteID == user.ID {
			history = append(history, rec)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
	fmt.Fprintln(w, history)

}

// 登出
func logout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	sessions[token] = ""
	fmt.Fprintln(w, "登出成功")
}

func main() {
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/info", getInfo)
	http.HandleFunc("/Info/update", updateInfo)
	http.HandleFunc("/info/history", infoHistory)
	http.HandleFunc("/logout", logout)
	http.ListenAndServe(":8080", nil)
}
