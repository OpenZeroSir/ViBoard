package Core

import (
	"MinView/Lib"
	"encoding/hex"
	"errors"

	jwt "github.com/golang-jwt/jwt/v5"
)

// 远程登陆信息
type RSession struct {
	User   string
	Passwd string
	Token  string
}

// 新建远程登陆信息
func NewRSession() *RSession {
	var rs RSession
	rs.User = "nni"
	rs.ResetRSession()
	return &rs
}

type RemoteSessionInfo struct {
	Token  string `json:"Token"`
	Passwd string `json:"Passwd"`
}

// 刷新登陆会话信息，应该个加载一个文件更新一次
func (r *RSession) ResetRSession() RemoteSessionInfo {
	r.Token = hex.EncodeToString([]byte(Lib.GenerateRandomStringN(16)))
	r.Passwd = Lib.GenerateRandomPasswd(12)
	// r.Passwd = "123"
	rsi := RemoteSessionInfo{
		Token:  r.Token,
		Passwd: r.Passwd,
	}
	return rsi
}

// 刷新远程登陆密码
func (r *RSession) ResetRSessionPasswd() string {
	r.Passwd = Lib.GenerateRandomPasswd(12)
	return r.Passwd
}

// 刷新远程token
func (r *RSession) ResetRSessionToken() string {
	r.Token = hex.EncodeToString([]byte(Lib.GenerateRandomStringN(16)))
	return r.Token
}

// 检查登陆账号密码
func (r RSession) CheckLogin(user, passwd string) bool {
	return (r.User == user) && (r.Passwd == passwd)
}

// 检查token
func (r RSession) CheckToken(rsToken string) bool {
	return r.Token == rsToken
}

// 获取密码
func (r RSession) GetPasswd() string {
	return r.Passwd
}

// 获取Token
func (r RSession) GetToken() string {
	return r.Token
}

// 登陆会话cookie
func (r RSession) GetCookie(username, passwd string) (string, error) {
	if !(username == r.User && passwd == r.Passwd) {
		return "", errors.New("username or passwd failed")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": r.User,
		"pass": r.Passwd,
	})
	tokenString, err := token.SignedString([]byte(Lib.VI_APP_KEYS[Lib.VI_VERSION]))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// 验证登陆会话cookie
func (r RSession) CheckCookie(c string) (bool, error) {
	token, err := jwt.Parse(c, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("check cookie failed")
		}
		return []byte(Lib.VI_APP_KEYS[Lib.VI_VERSION]), nil
	})
	if err != nil {
		return false, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["user"] == r.User && claims["pass"] == r.Passwd {
			return true, nil
		}
	}
	return false, errors.New("check cookie failed")
}
