package v1

import (
	"VGO/pi/internal/cache"
	"VGO/pi/internal/orm"
	"VGO/pi/internal/pkg/logfile"
	"VGO/pkg/fun"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

type Setting struct {
	*orm.Settings
	noSql bool
}

var history = []string{"talk_prologue"}

func (set Setting) List(c *gin.Context) {

	kv := orm.Settings{}.All()
	var data = make(map[string]string, len(kv))
	for _, v := range kv {
		data[v.Key] = v.Value
		if fun.InSliceString(v.Key, history) {
			data[v.Key+"_history"] = v.History
		}
	}
	jsonResult(c, http.StatusOK, data)
}

func (set Setting) Set(c *gin.Context) {
	err := c.ShouldBindBodyWith(&set, binding.JSON)
	var data = make(map[string]string, 2)
	if err == nil {
		if set.noSql {
			data[set.Key] = set.Value
		} else {
			addHistory := fun.InSliceString(set.Key, history)
			v := orm.Settings{}.Set(set.Key, set.Value, addHistory)
			data[v.Key] = v.Value
			if addHistory {
				data[v.Key+"_history"] = v.History
			}
		}

	} else {
		var setList []Setting
		err = c.ShouldBindBodyWith(&setList, binding.JSON)
		if err == nil {
			for _, v := range setList {
				addHistory := fun.InSliceString(v.Key, history)
				orm.Settings{}.Set(v.Key, v.Value, addHistory)
				data[v.Key] = v.Value
				if addHistory {
					data[v.Key+"_history"] = v.History
				}
			}
		}
	}
	go cache.SettingsCache.Update()

	if err != nil {
		logfile.Warning(err)
	}
	jsonResult(c, http.StatusOK, data)
}
