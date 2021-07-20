package orm

import (
	"VGO/pi/internal/config"
	"VGO/pi/internal/pkg/logfile"
	"database/sql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
)

var (
	db        *gorm.DB
	sqlDB     *sql.DB
	err       error
	pre       = "iq_"
	configENV = config.ENV
	logLevel  = logger.Error
)

type DatelineMap struct {
	Str       string `json:"str"`
	HighLight bool   `json:"highlight"`
}

type Table struct {
	Id int
}

func OpenDB() {
	if configENV["dev_mode"] == "true" {
		logLevel = logger.Silent
	}
	db, err = gorm.Open(sqlite.Open(configENV["home_path"]+configENV["db_path"]), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			TablePrefix:   pre,
		},
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		logfile.Fatal(err)
	}
	sqlDB, err = db.DB()
	if err != nil {
		log.Fatal(err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
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
	err = sqlDB.Close()
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
