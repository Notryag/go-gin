package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Dsn       string
	SecretKey string
}

var Cfg *Config

func Init() {
	//	安装读取yml文件的包 viper
	// go get github.com/spf13/viper
	viper.SetConfigFile("config/config.yml")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置文件失败", err))
	}
	host := viper.GetString("mysql.host")
	port := viper.GetString("mysql.port")
	user := viper.GetString("mysql.user")
	password := viper.GetString("mysql.password")
	dbname := viper.GetString("mysql.dbname")
	Cfg = &Config{
		Dsn:       fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbname),
		SecretKey: viper.GetString("jwt.secretKey"),
	}
}
