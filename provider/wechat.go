package provider

import (
	"azmesh-gateway/config"
	"azmesh-gateway/middleware"
	"azmesh-gateway/models"
	"encoding/json"
	"fmt"
	l4g "github.com/alecthomas/log4go"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

type Wechat struct {
	Name       string                              `json:"name"`
	AuthUrl    string                              `json:"auth_url"`
	SessionMap map[string]models.WechatSessionResp `json:"session_map"`
}

func (p *Wechat) Proxy(w http.ResponseWriter, r *http.Request) {
	l4g.Info("wechat proxy:%s", r.URL.String())
	cookie, err := r.Cookie(models.CookieName) // 替换为实际的 Cookie 名称
	if err != nil {
		http.Error(w, "cookie is empty", http.StatusForbidden)
		return
	}
	l4g.Info("wechat proxy get cookie:%s", cookie.String())

	valid, err := validateCookie(*cookie, []byte(models.SecretKey))
	if err != nil {
		l4g.Error("validateCookie err:%v", err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	if !valid {
		http.Error(w, "cookie is invalid", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("access-token", "Bearer "+cookie.Value)
	w.WriteHeader(http.StatusOK)
}

func (p *Wechat) Login(w http.ResponseWriter, r *http.Request) {
	l4g.Info("wechat login:%s", r.URL.String())
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "code is empty", http.StatusForbidden)
		return
	}

	// get wechat session from wechat server
	appName := r.URL.Query().Get("app_name")
	app := config.GetAppInfo(appName)
	url := p.AuthUrl + "?appid=" + app.Appid + "&secret=" + app.Appsecret + "&js_code=" + code + "&grant_type=authorization_code"
	data, err := middleware.SendGetRequest(url)
	if err != nil {
		l4g.Error("get wechat from [%s] token err:%v", url, err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	var sessionResp models.WechatSessionResp
	err = json.Unmarshal(data, &sessionResp)
	if err != nil {
		l4g.Error("wechat login json unmarshal wechatSessionResp err:%v", err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	if sessionResp.Errcode != 0 {
		l4g.Error("wechat login get session err:%v", sessionResp.Errmsg)
		http.Error(w, sessionResp.Errmsg, http.StatusForbidden)
		return
	}

	token := getToken(sessionResp)
	cookie := &http.Cookie{
		Name:    models.CookieName,
		Value:   token,
		Path:    "/",
		Expires: time.Now().Add(time.Minute * 120),
	}
	p.SessionMap[cookie.String()] = sessionResp
	http.SetCookie(w, cookie)
	l4g.Info("wechat login code:%s, set cookie:%s", code, cookie.String())
	w.WriteHeader(http.StatusOK)
}

func getToken(wechatSession models.WechatSessionResp) string {
	//genarate jwt Token
	token, err := middleware.GenerateJWTToken(wechatSession.Openid, time.Now().Add(time.Minute*120), []byte(models.SecretKey))
	if err != nil {
		l4g.Error("genarate jwt Token err:%v", err)
		return ""
	}

	return token
}

func validateCookie(cookie http.Cookie, secretKey []byte) (bool, error) {
	// 解析和验证 JWT token
	tokenString := cookie.Value
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 指定用于验证签名的密钥
		return secretKey, nil
	})
	if err != nil {
		return false, fmt.Errorf("解析 JWT token 失败: %v", err)
	}

	// 检查 token 是否有效和过期
	if !token.Valid {
		l4g.Error("validateCookie token is invalid:%s", tokenString)
		return false, nil // token 无效
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, fmt.Errorf("无效的 token 声明")
	}

	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
	currentTime := time.Now()
	if currentTime.After(expirationTime) {
		l4g.Error("validateCookie token is expired:%s", tokenString)
		return false, nil // token 已过期
	}

	return true, nil // Cookie 合法且未过期
}
