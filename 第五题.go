package main

import (
	//"fmt"
	"fmt"
	"net/http"
	// "time"
	//"io"
	//"bytes"
	"encoding/json"
	"os"
	"sync"
)
var users = make(map[string]User)
//var usersMu sync.Mutex
var sessions = make(map[string]string)
var userdetails = make(map[string]UserDetail)
var gamerecords = make(map[string]GameRecord)
//var sessionsMu sync.Mutex
type Rep struct {
    Username string `json:"username"` 
    Password string `json:"password"` 
}
type Check struct {
	ID int `json:"id"`
}
type User struct {
    Username string `json:"username"` 
    Password string `json:"password"`
	ID       int    `json:"id"`
    Bio      string `json:"bio"`
}
type UserDetail struct {
    ID       int    `json:"id"`       
    Username string `json:"username"` 
    Bio      string `json:"bio"`      
    Wins     int    `json:"wins"`     
    Losses   int    `json:"losses"`   
}	
type ChangeDetail struct {
    Oldpass  string `json:"old"` 
    NewID    int    `json:"newid"`       
    Newpass  string `json:"new"`
    Bio      string `json:"bio"`
}
type GameRecord struct {
    GameID  string  `json:"game_id"` 
    BlackID int     `json:"black_id"` 
    WhiteID int     `json:"white_id"` 
    Result  string  `json:"result"`   
    Margin  float64 `json:"margin"`   
}
//注册
func register(w http.ResponseWriter,r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"message" : "请使用POST请求",
		})
		return
	}
    var rep User
	err := json.NewDecoder(r.Body).Decode(&rep)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message" : "JSON解析错误",
		})
		return
	}
	_,exists := users[rep.username]
	if exists {
       w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message" : "用户名已经存在",
		})
		return
	}
	if !exists {
		users[rep.Username] = rep
	}

	if rep.Username == "" || rep.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message" : "用户名和密码不能为空",
		})
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
			"用户名" : rep.Username,
			"用户ID" : "rep.ID",
		})
}

//登录
func LoginRequest(w http.ResponseWriter,r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"message" : "请使用POST请求",
		})
		return
	}
    var rep Rep
	err := json.NewDecoder(r.Body).Decode(&rep)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message" : "JSON解析错误",
		})
		return
	}
	//usersMu.Lock()
	user,exists := users[rep.Username]
	//usersMu.Lock()
	if !exists || rep.Password != user.Password {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"message" : "用户名或密码错误",
		})
		return
  
	}
	sessionId := user.Username + "session_id"
	sessionsMu.Lock()
	sessions[sessionId] = user.Username
	sessionsMu.Unlock()
	http.SetCookie(w,&http.Cookie{
		Name : "sessionid",
		Value : sessionId,
		Path : "/",
	})
	w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message" : "登录成功",
			"name"    : user.Username,
		})
}

//修改资料
func Change(w http.ResponseWriter,r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"message" : "请使用POST请求",
		})
		return
	}
    var rep ChangeDetail
	err := json.NewDecoder(r.Body).Decode(&rep)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message" : "JSON解析错误",
		})
		return
	}

	cookie,err := r.Cookie("sessionid")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"message" : "登录状态失效",
		})
		return
	}
	//sessionsMu.Lock()
	username,exists := sessions[cookie.Value]
	//sessionsMu.Unlock()
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"message" : "登录状态失效",
		})
		return
	}

	//usersMu.Lock()
	user,exists := users[username]
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"message" : "用户不存在",
		})
		return
	}

	if user.Password != rep.Oldpass {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message" : "密码输入错误",
		})
		return
	}
	user.Bio = rep.Bio
	user.Password = rep.Newpass
	user.ID = rep.NewID
	users[user.Username] = user
    w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message" : "修改成功",
		})
}

//查询资料
func CheckDetail(w http.ResponseWriter,r *http.Request)  {
	w.Header().Set("Content-Type","application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"message" : "请使用POST请求",
		})
		return
	}
    var rep Check
	err := json.NewDecoder(r.Body).Decode(&rep)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message" : "JSON解析错误",
		})
		return
	}
	var detail UserDetail
	_,exists := userdetails["rep.ID"]
	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message" : "用户名不存在",
		})
		return
	}
	if exists {
		userdetails["rep.ID"] = detail
	}
    w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]UserDetail{
			"message" : userdetails["rep.ID"],
})
}

//接收对局服务上报的赛果
func accept(w http.ResponseWriter,r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"message" : "请使用POST请求",
		})
		return
	}
   
    var rep GameRecord
	_,exists := userdetails["rep.GameID"]
	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message" : "用户名不存在",
		})
		return
	}
	if exists {
		gamerecords["rep.GameID"] = rep
	}
    w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]UserDetail{
			"message" : userdetails["rep.GameID"],
})

}

func main() {
	http.HandleFunc("/register",register)
	http.HandleFunc("/loginrequest",LoginRequest)
	http.HandleFunc("/check",CheckDetail)
	http.HandleFunc("/change",Change)
	http.HandleFunc("/accept",accept)
	err := http.ListenAndServe(":8081",nil)
	if err != nil {
		fmt.Println("发送失败")
	}

}