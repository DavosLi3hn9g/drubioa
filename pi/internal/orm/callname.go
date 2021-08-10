package orm

type CallName struct {
	Id   int    `form:"id" xml:"id" json:"id" gorm:"PRIMARY_KEY;AUTO_INCREMENT"` //唯一ID
	Uid  int    `form:"uid" xml:"uid" json:"uid" `                               //用户UID
	Call string `form:"call" xml:"call" json:"call"`                             //称呼
}

func (_ CallName) All(wh *CallName) []CallName {
	var iList []CallName
	err = db.Where(wh).Order("uid asc,id asc").Find(&iList).Error
	if ErrDB(err) {
		return nil
	} else {
		return iList
	}
}
func (c CallName) Get(wh *CallName) *CallName {
	err = db.Where(wh).First(&c).Error
	if ErrDB(err) {
		return &c
	} else {
		return &c
	}
}
func (_ CallName) Add(data *CallName) *CallName {
	err = db.Create(data).Error
	if ErrDB(err) {
		return data
	} else {
		return data
	}
}

func (c CallName) Delete(k string, v int) error {
	return db.Where(k, v).Delete(&c).Error
}
func (c CallName) CreatTable() error {
	return db.Migrator().CreateTable(&c)
}
