package orm

type Auth struct {
	AuthId    int    `form:"auth_id" xml:"auth_id" json:"auth_id" gorm:"PRIMARY_KEY;AUTO_INCREMENT"` //唯一ID
	AiName    string `form:"ai_name" xml:"ai_name" json:"ai_name"`                                   //设备名称
	Password  string `form:"password" xml:"password" json:"password"`                                //用户密码
	Activated int    `form:"activated" xml:"activated" json:"activated"`                             //是否已初始化设置
}

func (u Auth) Get(uid int) *Auth {
	err = db.First(&u, uid).Error
	if ErrDB(err) {
		return &u
	} else {
		return &u
	}
}
func (u Auth) GetByName(AiName string) *Auth {
	err = db.Where("ai_name = ?", AiName).First(&u).Error
	if ErrDB(err) {
		return &u
	} else {
		return &u
	}
}
func (_ Auth) Save(data *Auth) *Auth {
	err = db.Save(&data).Error
	if ErrDB(err) {
		return data
	} else {
		return data
	}
}
func (_ Auth) CreatTable() {
	var u *Auth
	db.Table(pre + "users").CreateTable(&u)
}
