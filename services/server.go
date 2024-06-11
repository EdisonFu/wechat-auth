package services

import (
	"azmesh-gateway/config"
	"azmesh-gateway/provider"
	l4g "github.com/alecthomas/log4go"
	"net/http"
)

var ProviderMap map[string]provider.Provider

func initProvider() {
	cfg := config.GetConfig()
	if len(cfg.Providers) == 0 {
		panic("no provider")
	}

	ProviderMap = make(map[string]provider.Provider)
	for _, p := range cfg.Providers {
		ProviderMap[p.Name] = provider.NewProvider(p)
	}
}

// 使用
func StartServer(port string) {
	initProvider()

	http.HandleFunc("/login", Login)
	http.HandleFunc("/", Proxy)

	http.ListenAndServe(":"+port, nil)
}

func Proxy(w http.ResponseWriter, r *http.Request) {
	appName := r.Host
	cfg := config.GetConfig()
	for _, app := range cfg.Apps {
		if app.Name == appName {
			l4g.Info("start proxy:%v, url:%s", app.Provider, r.URL.String())
			ProviderMap[app.Provider].Proxy(w, r)
			return
		}
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	appName := r.URL.Query().Get("app_name")
	cfg := config.GetConfig()
	for _, app := range cfg.Apps {
		if app.Name == appName {
			l4g.Info("start login:%v, url:%s", app.Provider, r.URL.String())
			ProviderMap[app.Provider].Login(w, r)
			return
		}
	}
}
