package config

import (
	"encoding/json"
	l4g "github.com/alecthomas/log4go"
	"os"
)

type Config struct {
	Providers []Provider `json:"providers"`
	Apps      []App      `json:"apps"`
}

type Provider struct {
	Name    string `json:"name"`
	AuthUrl string `json:"auth_url"`
}
type App struct {
	Name      string `json:"name"`
	Provider  string `json:"provider"`
	Appid     string `json:"appid"`
	Appsecret string `json:"appsecret"`
}

var configInstance *Config

func InitConfig(cfgPath string) {
	configInstance = &Config{}
	//把config的json文件解析到configInstance中
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, configInstance)
	if err != nil {
		panic(err)
	}
	l4g.Info("config:%+v", configInstance)

	initAppInfo()
	return
}

func GetConfig() *Config {
	return configInstance
}

var appMap map[string]App

func initAppInfo() {
	appMap = make(map[string]App)
	for _, app := range configInstance.Apps {
		appMap[app.Name] = app
	}
}

func GetAppInfo(appName string) App {
	return appMap[appName]
}
