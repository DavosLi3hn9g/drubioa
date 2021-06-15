package orm

import (
	"VGO/pkg/fun"
	"strings"
)

type Settings struct {
	Key     string `form:"key" xml:"key" json:"key" gorm:"PRIMARY_KEY"` //名称
	Value   string `form:"value" xml:"value" json:"value"`              //数值
	History string `form:"history" xml:"history" json:"history"`        //设置历史
}

func (_ Settings) All() []Settings {
	var all []Settings
	err = db.Find(&all).Error
	if ErrDB(err) {
		return all
	} else {
		return all
	}
}
func (s Settings) Get(key string) *Settings {
	err = db.Where("key = ?", key).First(&s).Error
	if ErrDB(err) {
		return &s
	} else {
		return &s
	}
}
func (s Settings) Set(key, value string, addHistory bool) *Settings {
	exist := s.Get(key)
	var set Settings
	var history = exist.History
	set.Key = key
	set.Value = value
	if exist.Key == "" {
		err = db.Create(&set).Error
	} else {

		if set.Value != exist.Value && addHistory {
			if history != "" {
				hArr := strings.Split(history, "|")
				if exist.Value != "" && !fun.InSliceString(exist.Value, hArr) {
					if len(hArr) >= 5 {
						history = strings.Join(hArr[0:5], "|")
					}
					history = exist.Value + "|" + history
				}
			} else {
				history = exist.Value
			}
		}
		set.History = history
		err = db.Model(&exist).Updates(map[string]interface{}{"value": value, "history": history}).Error
	}
	if ErrDB(err) {
		return &set
	} else {
		return &set
	}
}

func (s *Settings) CreatTable() {
	db.Table(pre + "settings").CreateTable(&s)
}
