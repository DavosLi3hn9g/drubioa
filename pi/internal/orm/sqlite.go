package orm

import (
	"VGO/pi/internal/config"
	"VGO/pi/internal/pkg/logfile"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	db        *gorm.DB
	err       error
	pre       = "iq_"
	configENV = config.ENV
)

type DatelineMap struct {
	Str       string `json:"str"`
	HighLight bool   `json:"highlight"`
}

type Table struct {
	Id int
}

func OpenDB() {

	db, err = gorm.Open("sqlite3", configENV["home_path"]+configENV["db_path"])
	if err != nil {
		logfile.Fatal(err)
	}
	//设置表名前缀
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return pre + defaultTableName
	}
	if configENV["dev_mode"] == "true" {
		db.LogMode(true)
	} else {
		db.LogMode(false)
	}
	db.SingularTable(true) //约定模型名禁用表名复数
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	creatTable()
}

func creatTable() {
	var (
		auth       = Auth{}
		settings   = Settings{}
		intentions = Intentions{}
		queries    = Queries{}
		users      = Users{}
		callName   = CallName{}
		policies   = Policies{}
		logCall    = LogsCall{}
		logSms     = LogsSms{}
	)
	db.AutoMigrate(&auth, &settings, &intentions, &queries, &users, &callName, &policies, logCall, logSms)
}

func CloseDB() {
	err = db.Close()
	if err != nil {
		logfile.Error("数据库连接关闭出错了！")
	}
}
func ErrDB(err error) bool {
	if err != nil && err != gorm.ErrRecordNotFound {
		logfile.Error(err)
		return true
	} else if err == gorm.ErrRecordNotFound {
		return true
	}
	return false
}
