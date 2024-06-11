package models

const (
	SecretKey  = "wechat_session_secret"
	CookieName = "wechat_session"
)

type WechatSessionResp struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	Unionid    string `json:"unionid"`
	Errcode    int    `json:"errcode"`
	Errmsg     string `json:"errmsg"`
}
