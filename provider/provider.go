package provider

import (
	"azmesh-gateway/config"
	"azmesh-gateway/models"
	"net/http"
)

type Provider interface {
	Proxy(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
}

func NewProvider(provider config.Provider) Provider {
	switch provider.Name {
	case "wechat":
		return &Wechat{
			Name:       provider.Name,
			AuthUrl:    provider.AuthUrl,
			SessionMap: make(map[string]models.WechatSessionResp),
		}
	default:
		return nil
	}
}
