package config

import (
	"fmt"
	"gopkg.in/ini.v1"
	"log"
	"os"
)

var (
	iniFile *ini.File
	ENV     = make(map[string]string)
)

func init() {
	var err error
	if iniFile == nil {
		iniFile, err = ini.Load("./configs/.env")
		if err != nil {
			fmt.Printf("找不到配置文件: %v", err)
			os.Exit(1)
		}
	}
	if len(ENV) == 0 {
		confENV()
	}
}

func confENV() {
	path, err := iniFile.GetSection("path")
	if err != nil {
		log.Fatalf("未找到配置 'path': %v", err)
	}
	env := make(map[string]string)
	env["log_path"] = path.Key("log_path").MustString("data/logs/")
	env["wav_path"] = path.Key("wav_path").MustString("data/wav/")
	env["tmp_path"] = path.Key("tmp_path").MustString("data/tmp/")
	env["cache_path"] = path.Key("cache_path").MustString("data/cache/")
	env["pcm_path"] = path.Key("pcm_path").MustString("")
	env["home_path"] = path.Key("home_path").MustString("./")

	core, err := iniFile.GetSection("core")
	if err != nil {
		log.Fatalf("未找到配置 'core': %v", err)
	}
	env["bcm_ring"] = core.Key("bcm_ring").MustString("27")
	env["bcm_power"] = core.Key("bcm_power").MustString("")
	env["port"] = core.Key("port").MustString("")
	env["template"] = core.Key("template").MustString("default")

	dev, err := iniFile.GetSection("dev")
	if err != nil {
		log.Fatalf("未找到配置 'dev': %v", err)
	}
	env["dev_mode"] = dev.Key("dev_mode").MustString("false")
	env["allow_origin"] = dev.Key("allow_origin").MustString("http://localhost")

	database, err := iniFile.GetSection("database")
	if err != nil {
		log.Fatalf("未找到配置 'database': %v", err)
	}
	env["db_path"] = database.Key("db_path").MustString("data/sql/iqiar.db")

	ENV = env
}
