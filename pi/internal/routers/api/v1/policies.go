package v1

import (
	"VGO/pi/internal/cons"
	"VGO/pi/internal/orm"
	"VGO/pi/internal/pkg/logfile"
	"VGO/pkg/fun"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type Policy struct {
	*orm.Policies
	Checked    *orm.PolicyChecked `form:"checked" xml:"checked" json:"checked"`
	Intentions []*orm.Intentions  `form:"intentions" xml:"intentions" json:"intentions"`
}

func (_ *Policy) List(c *gin.Context) {
	var list = make([]*Policy, 0)
	var intentions = make(map[int][]*orm.Intentions)

	ItList := orm.Intentions{}.All(&orm.Intentions{})
	for _, c := range ItList {
		if fun.IsNumeric(c.End) {
			end, _ := strconv.Atoi(c.End)
			intentions[end] = append(intentions[end], c)
		}
	}
	data := orm.Policies{}.All(&orm.Policies{})
	for _, v := range data {
		if intentions[v.Id] == nil {
			intentions[v.Id] = []*orm.Intentions{}
		}
		vStr := v
		checked := orm.Policies{}.CheckedToStruct(v.Checked)
		list = append(list, &Policy{&vStr, checked, intentions[v.Id]})
	}
	jsonResult(c, http.StatusOK, list)
}
func (u *Policy) AddOrUpdate(c *gin.Context) {
	var dbPo *orm.Policies
	var checkedSlice []string
	if err := c.ShouldBind(&u); err != nil {
		logfile.Warning(err)
	}
	id := u.Id
	title := u.Title
	checked := u.Checked
	silent := u.Silent
	checkedMap := fun.Struct2Map(checked, "json")
	for k, v := range checkedMap {
		if v.(bool) {
			checkedSlice = append(checkedSlice, k)
		}
	}

	dbPo = orm.Policies{}.InsertOrUpdate(&orm.Policies{Id: id, Title: title, Checked: strings.Join(checkedSlice, "|"), Silent: silent})
	jsonResult(c, http.StatusOK, &Policy{
		dbPo, checked, nil,
	})
}
func (u *Policy) Del(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("id"))
	if id > 0 {
		err := orm.Intentions{}.EmptyType("end", strconv.Itoa(id))
		if !err {
			jsonErr(c, http.StatusBadRequest, cons.JsonErrDefault, "无法删除Intention！")
			return
		}
		errP := orm.Policies{}.Delete(id)
		if errP != nil {
			jsonErr(c, http.StatusBadRequest, cons.JsonErrDefault, "无法删除Policy！")
			return
		}

	}
	jsonResult(c, http.StatusOK, true)
}
