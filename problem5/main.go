package main

/*由于时间不足，省略了对用户是否存在的判断、用户名是否合法的判断和对密码的哈希加密🥺
*/

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

// Credentials 是注册 / 登录接口用来接收 HTTP 请求体的 JSON 结构
type Credentials struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}

type User struct {
	Password string `json:"password"` //用户密码
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

type UpdateDate struct {
	Password string `json:"password"`
	Username string `json:"username"` // 用户名
	Bio      string `json:"bio"`      // 个人简介
}

var mu sync.RWMutex

// 这个map用来储存用户信息，键是用户ID
var users = make(map[int]User)

//这个通过用户名获取用户信息，登录的时候验证密码用的
var usernameIndex = make(map[string]User)

// 这里把userID和token绑定
var tokens = make(map[string]int)

var nextID = 1

// 生成Token
func generateToken() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

// 身份认证
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		parts := strings.Fields(token)

		token = parts[1]

		mu.RLock()
		userID, exist := tokens[token]
		mu.RUnlock()
		if !exist {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "用户不存在",
			})
			c.Abort()
			return
		}

		c.Set("userID", userID)

		c.Next()
	}
}

func main() {
	r := gin.Default()

	r.POST("/register", registerHanlder)
	r.POST("/login", loginHandler)
	r.GET("/getuser", getuserHandler)
	r.GET("/histroy", historyHandler)

	r.POST("/receive",receievHandler)

	userGroup := r.Group("/user")
	userGroup.Use(AuthMiddleware())

	userGroup.GET("/leave", leaveHandler)
	userGroup.POST("/update", updateHandler)

}

//注册
func registerHanlder(c *gin.Context) {
	var req Credentials

	err := c.ShouldBindJSON(&req)
	if err != nil {据此展示棋手的胜场和败场
		c.JSON(400, gin.H{
			"error": "JSON格式错误",
		})
		return
	}

	user := User{
		Username: req.Username,
		Bio:      "我是天才棋据此展示棋手的胜场和败场手😋",
	}

	mu.Lock()
	user.ID = nextID
	nextID++
	users[user.ID] = user
	mu.Unlock()

	c.JSON(200,gin.H{
		"message":"注册成功",
		"username":user.Username,
	})
}

//登录
func loginHandler(c *gin.Context) {
	var req Credentials

	mu.RLock()
	user,_ := usernameIndex[req.Username]
	mu.RUnlock()

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400,gin.H{据此展示棋手的胜场和败场
			"error":"JSON数据错误",
		})
		return
	}

	if req.Password != user.Password {
		c.JSON(400,gin.H{
			"error": "密码错误",
		})
		return
	}
据此展示棋手的胜场和败场
	//token🍪
	token,err := generateToken()
	if err != nil {
		c.JSON(400,gin.H{
			"error":"token生成失败",
		})
		return
	}

	mu.Lock()
	tokens[token] = user.ID
	mu.Unlock()

	c.JSON(200,gin.H{
		"message": "登陆成功",
		"token": token,
	})
}

//退出登录,删除Token,其他不改
func leaveHandler(c *gin.Context) {
	Token, _ := c.Get("Token")
	token, _ := Token.(string)

	mu.Lock()
	tokens[token] = 0
	mu.Unlock()

	c.JSON(200,gin.H{
		"message":"退出登录成功",
	})
}

//查询棋手资料,在url中获取名字
func getuserHandler(c *gin.Context) {
	name := c.Param("username")
	user,_ := usernameIndex[name]
	c.JSON(200,gin.H{
		"message":"查询成功",
		"username":user.Username,
		"bio":user.Bio,
		"wins":user.Wins,
		"losses":user.Losses,
	})
}

//查询棋手历史,在url中获取名字
func historyHandler(c *gin.Context) {
	name := c.Param("username")
	user,_ := usernameIndex[name]
	c.JSON(200,gin.H{
		"message":"查询成功",
		"wins":user.Wins,
		"losses":user.Losses,
	})
}

//修改棋手资料
func updateHandler(c *gin.Context) {
	var req UpdateDate

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "请求数据格式错误",
		})
		return
	}

	userID, _ := c.Get("userID")
	id, _ := userID.(int)

	mu.Lock()
	defer mu.Unlock()

	user,_ := users[id]

	if req.Password != user.Password {
		c.JSON(400,gin.H{
			"error":"密码错误",
		})
		return
	}
	//输密码判断是不是本人

	delete(usernameIndex,user.Username)

	user.Username = req.Username
	user.Bio = req.Bio
	users[id] = user
	usernameIndex[user.Username] = user

	c.JSON(200, gin.H{
		"message": "修改成功",
		"username":req.Username,
		"bio":req.Bio,
	})
}

//接收信息
func receievHandler(c *gin.Context) {
	var req GameRecord

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400,gin.H{
			"error":"JSON数据错误",
		})
		return
	}

	expected := os.Getenv("GAME_SERVICE_KEY")    // 获取本地密钥
    actual := r.Header.Get("X-Game-Service-Key") // 获取请求携带的密钥

	if expected != actual {
		c.JSON(400,gin.H{
			"error":"密钥检验错误",
		})
		return
	}

	//获取对局用户信息
	mu.RLock()
	user1,_ := users[req.BlackID]
	user2,_ := users[req.WhiteID]
	mu.RUnlock()

	//判断胜负手
	mu.Lock()
	defer mu.Unlock()
	if req.Result == "black_win" {
		user1.Wins += 1
		user2.Losses += 1
	}else if req.Result == "white_win" {
		user1.Losses += 1
		user2.Wins += 1
	}

	users[user1.ID] = user1
	users[user2.ID] = user2
	usernameIndex[user1.Username] = user1
	usernameIndex[user2.Username] = user2

	c.JSON(200,gin.H{
		"message":"获取成功",
		"gamerecord":req,
	})
}