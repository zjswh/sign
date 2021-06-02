package conf

import (
	"fmt"
	"github.com/go-ini/ini"
	"io/ioutil"
	"log"
	"os"
)

var (
	Email struct {
		Host      string
		Port      int
		User      string
		Pass      string
		AdminUser string
	}
)

var cfg *ini.File

func Setup() {
	//判断配置是否存在
	configFile := "conf/app.ini"
	_, err := os.Stat(configFile)
	if os.IsNotExist(err) {
		//读取环境变量并生成配置文件
		conf := os.Getenv("config")
		if conf == "" {
			fmt.Println("环境变量缺失")
		}
		ioutil.WriteFile(configFile, []byte(conf), 0644)
	}
	cfg, err = ini.Load(configFile)
	if err != nil {
		log.Fatalf("conf.Setup, fail: %v", err)
	}
	mapTo("email", &Email)
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
