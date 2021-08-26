package orm

const (
	QueryModeKeywords int = 1
	QueryModeRegexp   int = 2
	QueryModeSql      int = 3
)

type Queries struct {
	Id     int    `form:"id" xml:"id" json:"id" gorm:"PRIMARY_KEY;AUTO_INCREMENT"` //唯一ID
	Sid    int    `form:"sid" xml:"sid" json:"sid"`                                //分组id
	Scores int    `form:"scores" xml:"scores" json:"scores"`                       //权重分数
	Query  string `form:"query" xml:"query" json:"query"`                          //查询规则
	Answer string `form:"answer" xml:"answer" json:"answer"`                       //应答话术
	Mode   int    `form:"mode" xml:"mode" json:"mode"`                             //匹配模式：1 中文字符串 2 正则表达式
}

func (_ Queries) All(wh *Queries) []*Queries {
	var qList []*Queries
	err = db.Where(wh).Order("scores desc,id asc").Find(&qList).Error
	if ErrDB(err) {
		return nil
	} else {
		return qList
	}
}
func (q Queries) Get(wh *Queries) *Queries {
	err = db.Where(wh).First(&q).Error
	if ErrDB(err) {
		return &q
	} else {
		return &q
	}
}
func (_ Queries) AllByKeywords(keywords []string) []*Queries {
	var wh = "mode = 1 AND ("
	var like []interface{}
	for k, v := range keywords {
		if k == 0 {
			wh += " query LIKE  ? "
		} else {
			wh += " OR query LIKE  ? "
		}
		like = append(like, "%"+v+"%")
	}
	wh += ")"
	var qList []*Queries
	err = db.Where(wh, like...).Find(&qList).Error
	if ErrDB(err) {
		return nil
	} else {
		return qList
	}
}
func (q Queries) InsertOrUpdate(data *Queries) *Queries {
	if data.Id > 0 {
		err = db.Model(&q).Where("id = ?", data.Id).Updates(data).Error
	} else {
		err = db.Create(data).Error
		db.Order("id desc").First(data)
	}
	if ErrDB(err) {
		return data
	} else {
		return data
	}
}
func (q Queries) Delete(k string, v int) error {
	return db.Where(k, v).Delete(&q).Error
}
func (q Queries) CreatTable() error {
	return db.Migrator().CreateTable(&q)
}
