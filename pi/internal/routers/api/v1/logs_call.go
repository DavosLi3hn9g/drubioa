package v1

import (
	"VGO/pi/internal/cons"
	"VGO/pi/internal/core"
	"VGO/pi/internal/orm"
	"VGO/pi/internal/sim"
	"VGO/pkg/fun"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strconv"
	"time"
)

type LogsCall struct {
	*orm.LogsCall
	DatelineStr   string `form:"dateline_str" xml:"dateline_str" json:"dateline_str"`
	RecordingPath string `form:"recording_path" xml:"recording_path" json:"recording_path"`
	RecordingUrl  string `form:"recording_url" xml:"recording_url" json:"recording_url"`
}

func (_ *LogsCall) List(c *gin.Context) {
	var list []*LogsCall
	var data []*orm.LogsCall
	var count int64
	isWAV := c.Query("is_wav") != ""
	page := c.Query("page")
	p, _ := strconv.Atoi(page)
	if isWAV {
		data = orm.LogsCall{}.All(`recording <> ""`, p)
		count = orm.LogsCall{}.Count(`recording <> ""`)
	} else {
		data = orm.LogsCall{}.All(&orm.LogsCall{}, p)
		count = orm.LogsCall{}.Count(&orm.LogsCall{})
	}

	for _, v := range data {
		path := configENV["wav_path"] + v.Recording + ".wav"
		url := "//" + c.Request.Host + "/" + path
		list = append(list, &LogsCall{v, time.Unix(int64(v.TimeStart), 0).In(sim.TimeLoc).Format("2006-01-02 15:04:05"), path, url})
	}
	jsonResult(c, http.StatusOK, map[string]interface{}{
		"count": count,
		"list":  list,
	})
}

func (i *LogsCall) AddOrUpdate(call orm.LogsCall) {
	var old *orm.LogsCall
	if call.Id > 0 {
		old = orm.LogsCall{}.Get(call.Id)
		if old.TimeStart > 0 && call.TimeEnd > old.TimeStart {
			call.Minute = int(math.Ceil(float64(old.TimeStart-call.TimeEnd) / 60))
		}
	}
	orm.LogsCall{}.Add(&call)
}

func (i *LogsCall) Del(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("id"))
	if id > 0 {
		err := orm.LogsCall{}.Delete(id)
		if err != nil {
			jsonErr(c, http.StatusBadRequest, cons.JsonErrDefault, "无法删除LogsCall！")
			return
		}
	}
	jsonResult(c, http.StatusOK, true)
}

func (i *LogsCall) Cost(c *gin.Context) {
	var m int
	data := orm.LogsCall{}.All(&orm.LogsCall{}, 0)
	for _, v := range data {
		m += v.Minute
	}
	cost := fun.Round(float64(m)*core.Setting.UnitPrice, 2)
	jsonResult(c, http.StatusOK, map[string]interface{}{
		"minute":     m,
		"cost":       cost,
		"unit_price": core.Setting.UnitPrice,
	})
}
