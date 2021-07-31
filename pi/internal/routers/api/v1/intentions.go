package v1

import (
	"VGO/pi/internal/cache"
	"VGO/pi/internal/cons"
	"VGO/pi/internal/orm"
	"VGO/pi/internal/pkg/logfile"
	"VGO/pkg/fun"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type Intention struct {
	*orm.Intentions
	Queries []*orm.Queries `form:"queries" xml:"queries" json:"queries"`
}

func (_ *Intention) List(c *gin.Context) {
	data := cache.IntentionsCache.New(true)
	jsonResult(c, http.StatusOK, data.List)
}

func (i *Intention) AddOrUpdateSid(c *gin.Context) {
	var dbIn *orm.Intentions
	if err := c.ShouldBind(&i); err != nil {
		logfile.Warning(err)
	}
	sid := i.Sid
	title := i.Title
	end := i.End
	level := i.Level
	dbIn = orm.Intentions{}.InsertOrUpdate(&orm.Intentions{Sid: sid, Title: title, End: end, Level: level})
	go cache.IntentionsCache.Update()
	jsonResult(c, http.StatusOK, &Intention{
		dbIn, []*orm.Queries{},
	})
}

func (i *Intention) Del(c *gin.Context) {
	sid, _ := strconv.Atoi(c.PostForm("sid"))
	if sid > 0 {
		err := orm.Queries{}.Delete("sid = ?", sid)
		if err != nil {
			jsonErr(c, http.StatusBadRequest, cons.JsonErrDefault, "无法删除Query！")
			return
		}
		err = orm.Intentions{}.Delete(sid)
		if err != nil {
			jsonErr(c, http.StatusBadRequest, cons.JsonErrDefault, "无法删除Intention！")
			return
		}
	}
	jsonResult(c, http.StatusOK, true)
	go cache.IntentionsCache.Update()
}

func (i *Intention) AddOrUpdateQuery(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("id"))
	sid, _ := strconv.Atoi(c.PostForm("sid"))
	query := c.PostForm("q")
	answer := c.PostForm("answer")
	scores, _ := strconv.Atoi(c.PostForm("scores"))
	mode, _ := strconv.Atoi(c.PostForm("mode"))
	if query == "" {
		jsonErr(c, http.StatusBadRequest, cons.JsonErrDefault, "缺少参数，请检查！")
		return
	}
	query = strings.Replace(query, "，", ",", -1)
	query = strings.Replace(query, "|", ",", -1)
	query = strings.Replace(query, "｜", ",", -1)
	answer = strings.Replace(answer, "|", ",", -1)
	answer = strings.Replace(answer, "｜", ",", -1)

	switch mode {
	case orm.QueryModeKeywords:
		keywords := strings.Split(query, ",")
		oldList := orm.Queries{}.AllByKeywords(keywords)
		for _, v := range oldList {
			if v.Mode == orm.QueryModeKeywords {
				vQuery := strings.Split(v.Query, ",")
				for _, q := range vQuery {
					if fun.InSliceString(q, keywords) {
						if v.Id == id {
							continue
						} else if v.Sid != sid {
							intention := orm.Intentions{}.Get(v.Sid)
							jsonErr(c, http.StatusBadRequest, cons.JsonErrDefault, "关键词["+q+"]已经存在于["+intention.Title+"]分类，请勿重复！")
						} else if v.Sid == sid {
							intention := orm.Intentions{}.Get(sid)
							jsonErr(c, http.StatusBadRequest, cons.JsonErrDefault, "关键词["+q+"]已经存在于["+intention.Title+"]分类，请勿重复！")
						} else {
							jsonErr(c, http.StatusBadRequest, cons.JsonErrDefault, "关键词["+q+"]已经存在，请勿重复！")
						}
						return
					}
				}
			}
		}
	}

	dbQuery := orm.Queries{}.InsertOrUpdate(&orm.Queries{Id: id, Sid: sid, Query: query, Answer: answer, Scores: scores, Mode: mode})
	jsonResult(c, http.StatusOK, dbQuery)
	go cache.IntentionsCache.Update()
}

func (i *Intention) DelQuery(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("id"))
	if id > 0 {
		err := orm.Queries{}.Delete("id = ?", id)
		if err != nil {
			jsonErr(c, http.StatusBadRequest, cons.JsonErrDefault, "无法删除Query！")
			return
		}
	}
	jsonResult(c, http.StatusOK, true)
	go cache.IntentionsCache.Update()
}
