package v1

import (
	"VGO/pi/internal/cache"
	"VGO/pi/internal/cons"
	"VGO/pi/internal/orm"
	"VGO/pi/internal/pkg/logfile"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type User struct {
	*orm.Users
	Call []orm.CallName `form:"call" xml:"call" json:"call"`
}

func (_ *User) List(c *gin.Context) {
	data := cache.UsersCache.New(true)
	jsonResult(c, http.StatusOK, data.List)
}
func (u *User) AddOrUpdate(c *gin.Context) {
	var dbUser *orm.Users
	if err := c.ShouldBind(&u); err != nil {
		logfile.Warning(err)
	}
	uid := u.Uid
	tel := u.Tel
	username := u.UserName
	dbUser = orm.Users{}.InsertOrUpdate(&orm.Users{Uid: uid, Tel: tel, UserName: username})
	jsonResult(c, http.StatusOK, &User{
		dbUser, nil,
	})
	cache.UsersCache.Update()
}
func (u *User) Del(c *gin.Context) {
	uid, _ := strconv.Atoi(c.PostForm("uid"))
	if uid > 0 {
		err := orm.CallName{}.Delete("uid = ?", uid)
		if err != nil {
			jsonErr(c, http.StatusBadRequest, cons.JsonErrDefault, "无法删除call！")
			return
		}
		err = orm.Users{}.Delete(uid)
		if err != nil {
			jsonErr(c, http.StatusBadRequest, cons.JsonErrDefault, "无法删除user！")
			return
		}

	}
	jsonResult(c, http.StatusOK, true)
	cache.UsersCache.Update()
}

func (u *User) AddCall(c *gin.Context) {
	var dbCall *orm.CallName
	uid, _ := strconv.Atoi(c.PostForm("uid"))
	call := c.PostForm("call")
	if uid > 0 && call != "" {
		old := orm.CallName{}.Get(&orm.CallName{Call: call})
		if old.Id > 0 {
			jsonErr(c, http.StatusBadRequest, cons.JsonErrDefault, "这个称呼已经存在，请勿重复！")
			return
		} else {
			dbCall = orm.CallName{}.Add(&orm.CallName{Uid: uid, Call: call})
			//callList = append(callList, dbCall)
			jsonResult(c, http.StatusOK, dbCall)
			cache.UsersCache.Update()
		}
	} else {
		jsonErr(c, http.StatusBadRequest, cons.JsonErrDefault, "缺少参数，请检查！")
	}
}

func (u *User) DelCall(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("id"))
	if id > 0 {
		err := orm.CallName{}.Delete("id = ?", id)
		if err != nil {
			jsonErr(c, http.StatusBadRequest, cons.JsonErrDefault, "无法删除call！")
			return
		}
	}
	jsonResult(c, http.StatusOK, true)
	cache.UsersCache.Update()
}
