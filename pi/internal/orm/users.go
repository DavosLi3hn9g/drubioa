package orm

import "gorm.io/gorm"

type Users struct {
	Uid      int    `form:"uid" xml:"uid" json:"uid" gorm:"PRIMARY_KEY;AUTO_INCREMENT"` //唯一ID
	UserName string `form:"username" xml:"username" json:"username"`                    //用户名称
	Tel      string `form:"tel" xml:"tel" json:"tel"`                                   //手机号
	Email    string `form:"email" xml:"email" json:"email"`                             //邮箱地址
	Push     string `form:"push" xml:"push" json:"push"`                                //推送通道
	Lev      int    `form:"lev" xml:"lev" json:"lev"`                                   //优先级
	Hits     int    `form:"hits" xml:"hits" json:"hits"`                                //触发次数
}

func (_ Users) All(wh *Users) []Users {
	var iList []Users
	err = db.Where(wh).Order("uid asc").Find(&iList).Error
	if ErrDB(err) {
		return nil
	} else {
		return iList
	}
}
func (u Users) Get(uid int) *Users {
	err = db.First(&u, uid).Error
	if ErrDB(err) {
		return &u
	} else {
		return &u
	}
}
func (u Users) InsertOrUpdate(data *Users) *Users {
	if data.Uid > 0 {
		err = db.Omit("hits").Updates(data).Error
	} else {
		err = db.Omit("hits").Create(data).Error
		db.Order("uid desc").First(&u)
	}
	if ErrDB(err) {
		return data
	} else {
		return data
	}
}
func (u *Users) IncHit(num int) {
	db.Model(&u).UpdateColumn("hits", gorm.Expr("hits + ?", num))
}
func (u Users) Delete(uid int) error {
	return db.Where("uid = ?", uid).Delete(&u).Error
}
func (u Users) CreatTable() error {
	return db.Migrator().CreateTable(&u)
}
