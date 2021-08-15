package orm

type LogsSms struct {
	Id       int    `form:"id" xml:"id" json:"id" gorm:"PRIMARY_KEY;"` //唯一ID
	Text     string `form:"text" xml:"text" json:"text"`               //内容text
	TelFrom  string `form:"tel_from" xml:"tel_from" json:"tel_from"`   //发送方电话号码
	TelTo    string `form:"tel_to" xml:"tel_to" json:"tel_to"`         //接收方电话号码
	Dateline int    `form:"dateline" xml:"dateline" json:"dateline"`   //接收时间
	SmsId    string `form:"sms_id" xml:"sms_id" json:"sms_id"`         //SmsID
}

func (_ LogsSms) All(wh interface{}, page int) []*LogsSms {
	var pList []*LogsSms
	var limit = 20
	d := db.Where(wh).Order("dateline desc,id asc")
	if page > 0 {
		page = page - 1
		d = d.Offset(page * limit).Limit(limit)
	}
	err = d.Find(&pList).Error
	if ErrDB(err) {
		return nil
	} else {
		return pList
	}
}
func (_ LogsSms) Count(wh interface{}) int64 {
	var count int64
	err = db.Where(wh).Count(&count).Error
	if ErrDB(err) {
		return 0
	} else {
		return count
	}
}
func (l LogsSms) Get(dateline int) *LogsSms {
	err = db.Where("dateline = ?", dateline).First(&l).Error
	if ErrDB(err) {
		return &l
	} else {
		return &l
	}
}
func (l LogsSms) InsertOrUpdate(data *LogsSms) *LogsSms {
	db.Where("dateline = ? AND sms_id = ?", data.Dateline, data.SmsId).First(&l)
	if l.Id > 0 {
		err = db.Updates(data).Error
	} else {
		err = db.Create(data).Error
	}
	if ErrDB(err) {
		return data
	} else {
		return data
	}
}
func (l LogsSms) Delete(id int) error {
	return db.Where("id = ?", id).Delete(&l).Error
}
func (l LogsSms) CreatTable() error {
	return db.Migrator().CreateTable(&l)
}
