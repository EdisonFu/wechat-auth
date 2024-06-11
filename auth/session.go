package auth

// 使用map保存和前端交互的session
var SessionMap map[string]string

func init() {
	SessionMap = make(map[string]string)
}

// 保存session
func SaveSession(sessionId string, openid string) {

}
