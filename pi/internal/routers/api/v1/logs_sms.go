package v1

import (
	"VGO/pi/internal/cons"
	"VGO/pi/internal/core"
	"VGO/pi/internal/orm"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type LogsSms struct {
	*orm.LogsSms
	DatelineStr string `form:"dateline_str" xml:"dateline_str" json:"dateline_str"`
}

func (l *LogsSms) List(c *gin.Context) {
	l.AddOrUpdate(c)
	page := c.Query("page")
	p, _ := strconv.Atoi(page)
	listOrm := orm.LogsSms{}.All(&LogsSms{}, 0)
	var text = make(map[int]string, len(listOrm))
	var sms = make(map[int]*orm.LogsSms, len(listOrm))
	var sortedKeys = make([]int, 0)
	for _, v := range listOrm {
		text[v.Dateline] += v.Text
		sms[v.Dateline] = v
	}
	for k, _ := range sms {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(sortedKeys)))
	var limit = 20
	var res []*LogsSms
	var list = make([]*LogsSms, 0, len(sms))
	if p > 0 {
		p = p - 1
	}
	for _, k := range sortedKeys {
		v := sms[k]
		v.Text = text[k]
		list = append(list, &LogsSms{v, time.Unix(int64(v.Dateline), 0).Format("2006-01-02 15:04:05")})
	}
	count := len(list)
	start := limit * p
	end := limit*p + limit
	if end <= count {
		res = list[start:end]
	} else if start < count {
		res = list[start:count]
	} else {
		res = nil
	}
	jsonResult(c, http.StatusOK, map[string]interface{}{
		"count": count,
		"list":  res,
	})
}

func (l *LogsSms) AddOrUpdate(c *gin.Context) {
	for _, v := range core.SmsList {
		orm.LogsSms{}.InsertOrUpdate(v)
	}
}

func (l *LogsSms) Del(c *gin.Context) {
	//id := c.PostForm("id")
	dateline, _ := strconv.Atoi(c.PostForm("dateline"))
	if dateline > 0 {
		smsList := orm.LogsSms{}.All(&orm.LogsSms{Dateline: dateline}, 0)
		go func() {
			for _, v := range smsList {
				err := orm.LogsSms{}.Delete(v.SmsId)
				if err != nil {
					jsonErr(c, http.StatusBadRequest, cons.JsonErrDefault, "无法删除LogsSms！")
					return
				}
				LogsPi := new(LogsPi)
				LogsPi.port()
				if core.SerialAT.Port != nil {
					core.SerialAT.AT("AT+CMGD=" + v.SmsId)
				}
				delete(core.SmsList, v.SmsId)
				time.Sleep(1 * time.Second)
			}
		}()
	}
	jsonResult(c, http.StatusOK, true)
}
