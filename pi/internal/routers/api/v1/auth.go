package v1

import (
	"VGO/pi/internal/cons"
	"VGO/pi/internal/orm"
	"VGO/pkg/fun"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type Auth struct {
	UserAvatar string `json:"user_avatar"`
	AiName     string `json:"ai_name"`
	IsLogged   bool   `json:"is_logged"`
	Token      Token  `json:"token"`
	Activated  int    `json:"activated"`
}
type Token struct {
	Exp time.Time `json:"exp"`
	Str string    `json:"str"`
}
type AuthCH struct {
	List []Auth
}

var CacheAuth = make(map[int]*Auth)

const tokenMaxAge time.Duration = 3600 * 24 * 7 //Second

const defaultUid = 1

func (a *Auth) Login(c *gin.Context) {
	var data Auth
	password := c.PostForm("password")
	old := orm.Auth{}.Get(defaultUid)
	if old.Password == ENC(password) {
		uid := old.AuthId
		token := fun.MD5(password + fmt.Sprintf("iqiar%d", time.Now().UnixNano()))
		data.IsLogged = true
		data.UserAvatar = ""
		data.AiName = old.AiName
		c.SetCookie("token", token, int(tokenMaxAge), "/", "", false, true)
		CacheAuth[uid] = new(Auth)
		data.Token = CacheAuth[uid].SetToken(TokenENC(uid, token))
		jsonResult(c, http.StatusOK, data)
		return
	} else {
		if old.Password != "" {
			jsonErr(c, http.StatusUnauthorized, cons.JsonErrDefault, "密码错误！")
		} else {
			jsonErr(c, http.StatusUnauthorized, cons.JsonErrDefault, "账号未激活！")
		}

		return
	}

}

func (a *Auth) Register(c *gin.Context) {

	password1 := c.PostForm("password")
	password2 := c.PostForm("password2")
	if password1 == password2 && password1 != "" {
		s := orm.Auth{}.Get(defaultUid)
		if s.AuthId == 0 {
			newPassword := ENC(password1)
			orm.Auth{}.Save(&orm.Auth{AuthId: defaultUid, AiName: "我的设备", Password: newPassword, Activated: 0})
			a.Login(c)
			return
		} else {
			jsonErr(c, http.StatusBadRequest, cons.JsonErrDefault, "账号已经激活过！")
			return
		}
	}
	jsonErr(c, http.StatusBadRequest, cons.JsonErrDefault, "激活出错！")
	return

}

func (a *Auth) Update(c *gin.Context) {
	var data bool
	password1 := c.PostForm("password")
	password2 := c.PostForm("password2")
	//aiName := c.PostForm("ai_name")
	if password1 == password2 && password1 != "" {
		newPassword := ENC(password1)
		orm.Auth{}.Save(&orm.Auth{AuthId: defaultUid, AiName: "我的设备", Password: newPassword})
		data = true
		jsonResult(c, http.StatusOK, data)
	} else {
		data = false
		jsonResult(c, http.StatusBadRequest, data)
	}

}

func (a *Auth) Info(c *gin.Context) {
	var data = a
	s := orm.Auth{}.Get(defaultUid)
	data.IsLogged = true
	data.UserAvatar = ""
	data.AiName = s.AiName
	data.Activated = s.Activated
	jsonResult(c, http.StatusOK, data)
}

func (a *Auth) First(c *gin.Context) {
	var data bool
	s := orm.Auth{}.Get(defaultUid)
	if s.AuthId == 0 {
		data = true
	} else {
		data = false
	}
	jsonResult(c, http.StatusOK, data)
}
func (a *Auth) Logout(c *gin.Context) {
	var data = a
	data = nil
	username, _, ok := c.Request.BasicAuth()
	if ok {
		uid, _ := strconv.Atoi(username)
		CacheAuth[uid].SetToken("")
		c.SetCookie("token", "", 1, "/", "", false, true)
	}
	jsonResult(c, http.StatusOK, data)
}

func (a *Auth) Reset() string {
	str := fun.SubStr(fun.MD5(strconv.Itoa(int(time.Now().UnixNano()))), 0, 10)
	password := ENC(fun.MD5(str))
	orm.Auth{}.Save(&orm.Auth{AuthId: defaultUid, AiName: "我的设备", Password: password})
	return str
}
func ENC(password string) string {
	return fun.MD5(password)
}

func TokenENC(uid int, token string) string {
	return fun.Base64Encode(fmt.Sprintf("%d:%s", uid, token))
}
func (a *Auth) GetToken() string {
	return a.Token.Str
}
func (a *Auth) SetToken(token string) Token {
	if a == nil {
		return Token{}
	}
	a.Token.Str = token
	if token != "" {
		a.Token.Exp = time.Now().Add(time.Second * tokenMaxAge)
	} else {
		a.Token.Exp = time.Time{}
	}

	return a.Token
}
