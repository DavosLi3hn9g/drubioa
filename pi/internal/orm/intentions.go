package orm

type Intentions struct {
	Sid   int    `form:"sid" xml:"sid" json:"sid" gorm:"PRIMARY_KEY;AUTO_INCREMENT"` //唯一ID
	Title string `form:"title" xml:"title" json:"title"`                             //分类名称
	End   string `form:"end" xml:"end" json:"end"`                                   //终结策略
	Level int    `form:"level" xml:"level" json:"level"`                             //优先级
	Hello bool   `form:"hello" xml:"hello" json:"hello"`                             //是否闲聊
	Hits  int    `form:"hits" xml:"hits" json:"hits"`                                //触发次数
}

func (_ Intentions) All(wh *Intentions) []*Intentions {
	var iList []*Intentions
	err = db.Where(wh).Order("Level desc").Find(&iList).Error
	if ErrDB(err) {
		return nil
	} else {
		return iList
	}
}
func (i Intentions) Get(sid int) *Intentions {
	err = db.First(&i, sid).Error
	if ErrDB(err) {
		return &i
	} else {
		return &i
	}
}
func (i Intentions) InsertOrUpdate(data *Intentions) *Intentions {
	if data.Sid > 0 {
		err = db.Omit("hits").Updates(data).Error
	} else {
		err = db.Omit("hits").Create(data).Error
		db.Order("sid desc").First(&i)
	}
	if ErrDB(err) {
		return data
	} else {
		return data
	}
}
func (i Intentions) EmptyType(typename string, value string) bool {
	err = db.Model(&i).Where(typename+" = ?", value).Update(typename, "").Error
	if ErrDB(err) {
		return false
	} else {
		return true
	}
}
func (i Intentions) Delete(sid int) error {
	return db.Where("sid = ?", sid).Delete(&i).Error
}
func (i Intentions) CreatTable() error {
	return db.Migrator().CreateTable(&i)
}
