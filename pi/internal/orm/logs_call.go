package orm

type LogsCall struct {
	Id        int    `form:"id" xml:"id" json:"id" gorm:"PRIMARY_KEY;AUTO_INCREMENT"` //唯一ID
	Text      string `form:"text" xml:"text" json:"text"`                             //内容Text
	Content   string `form:"content" xml:"content" json:"content"`                    //内容Json
	Recording string `form:"recording" xml:"recording" json:"recording"`              //录音文件地址
	TelFrom   string `form:"tel_from" xml:"tel_from" json:"tel_from"`                 //主叫电话号码
	TelTo     string `form:"tel_to" xml:"tel_to" json:"tel_to"`                       //被叫电话号码
	Intention string `form:"intention" xml:"intention" json:"intention"`              //已触发意图，可能包含多个触发
	Policy    string `form:"policy" xml:"policy" json:"policy"`                       //已触发挂机策略，仅包含最后一个触发
	TimeStart int    `form:"time_start" xml:"time_start" json:"time_start"`           //通话开始时间
	TimeEnd   int    `form:"time_end" xml:"time_end" json:"time_end"`                 //通话结束时间
	Minute    int    `form:"minute" xml:"minute" json:"minute"`                       //通话分钟
}

func (_ LogsCall) All(wh interface{}, page int) []*LogsCall {
	var pList []*LogsCall
	var limit = 20
	d := db.Where(wh).Order("id desc")
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
func (_ LogsCall) Count(wh interface{}) int64 {
	var count int64
	var pList []*LogsCall
	err = db.Where(wh).Find(&pList).Count(&count).Error
	if ErrDB(err) {
		return 0
	} else {
		return count
	}
}
func (l LogsCall) Get(id int) *LogsCall {
	err = db.First(&l, id).Error
	if ErrDB(err) {
		return &l
	} else {
		return &l
	}
}
func (_ LogsCall) Add(data *LogsCall) *LogsCall {
	err = db.Create(data).Error
	if ErrDB(err) {
		return data
	} else {
		return data
	}
}

func (l LogsCall) Updates(data LogsCall) LogsCall {
	err = db.Model(&l).Updates(&data).Error
	if ErrDB(err) {
		return data
	} else {
		return data
	}
}
func (l LogsCall) Delete(id int) error {
	return db.Where("id = ?", id).Delete(&l).Error
}
func (l LogsCall) CreatTable() error {
	return db.Migrator().CreateTable(&l)
}
