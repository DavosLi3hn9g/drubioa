package orm

import (
	"github.com/jinzhu/gorm"
	"strings"
)

type Policies struct {
	Id      int    `form:"id" xml:"id" json:"id" gorm:"PRIMARY_KEY;AUTO_INCREMENT"` //唯一ID
	Title   string `form:"title" xml:"title" json:"title"`                          //策略名称
	Checked string `form:"checked" xml:"checked" json:"checked"`                    //策略选项
	Silent  bool   `form:"silent" xml:"silent" json:"silent"`                       //是否静默
	Hits    int    `form:"hits" xml:"hits" json:"hits"`                             //触发次数
}

type PolicyChecked struct {
	Sms       bool `form:"sms" xml:"sms" json:"sms"`
	Call      bool `form:"call" xml:"call" json:"call"`
	Push      bool `form:"push" xml:"push" json:"push"`
	Blacklist bool `form:"blacklist" xml:"blacklist" json:"blacklist"`
}

func (_ Policies) CheckedToStruct(str string) *PolicyChecked {
	checkedArr := strings.Split(str, "|")
	checkedMap := map[string]bool{}
	for _, c := range checkedArr {
		checkedMap[c] = true
	}
	checked := &PolicyChecked{
		Sms:       checkedMap["sms"],
		Call:      checkedMap["call"],
		Push:      checkedMap["push"],
		Blacklist: checkedMap["blacklist"],
	}
	return checked
}

func (_ Policies) All(wh *Policies) []Policies {
	var pList []Policies
	err = db.Where(wh).Order("id asc").Find(&pList).Error
	if ErrDB(err) {
		return nil
	} else {
		return pList
	}
}
func (p *Policies) Get(id int) *Policies {
	err = db.First(&p, id).Error
	if ErrDB(err) {
		return p
	} else {
		return p
	}
}
func (_ Policies) Save(data *Policies) *Policies {
	err = db.Omit("hits").Save(&data).Error
	if ErrDB(err) {
		return data
	} else {
		return data
	}
}
func (p *Policies) IncHit(num int) {
	db.Model(&p).UpdateColumn("hits", gorm.Expr("hits + ?", num))
}
func (p Policies) Delete(id int) error {
	return db.Where("id = ?", id).Delete(&p).Error
}
func (_ Policies) CreatTable() {
	var p *Policies
	db.Table(pre + "policies").CreateTable(&p)
}
