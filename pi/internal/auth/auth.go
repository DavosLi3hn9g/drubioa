package auth

import (
	v1 "VGO/pi/internal/routers/api/v1"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

func CheckCookie(auth map[int]*v1.Auth) gin.HandlerFunc {
	return func(c *gin.Context) {

		token, err := c.Cookie("token")
		if err != nil {
			token = "NotSet"
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		if token == auth[1].Token.Str {
			log.Printf("Cookie value: %s \n", token)
		} else {
			log.Println("Cookie err")
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		// 请求前
		c.Next()
		// 请求后

	}
}
func CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, password, ok := c.Request.BasicAuth()
		if ok {
			uid, _ := strconv.Atoi(user)
			cacheAuth := v1.CacheAuth
			if token, ok := cacheAuth[uid]; ok {
				if v1.TokenENC(uid, password) != token.Token.Str || token.Token.Exp.Before(time.Now()) {
					log.Println("token error ", password)
					c.AbortWithStatus(http.StatusForbidden)
				}
			} else {
				log.Println("token empty")
				c.AbortWithStatus(http.StatusUnauthorized)
			}
		} else {
			log.Println("AuthUserKey error")
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		// 请求前
		c.Next()
		// 请求后

	}
}
