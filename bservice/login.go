package bservice

import (
	"encoding/base64"
	"errors"
	"fmt"
)

// LoginInfo 登陆信息
type LoginInfo struct {
	Username  string
	Password  string
	Csrf      string
	UID       string
	Cookies   string
	Headers   map[string]string
	AccessKey string
}

// CookieData cookie信息
type CookieData struct {
	TokenInfo struct {
		AccessToken string `json:"access_token"`
	} `json:"token_info"`
	CookieInfo struct {
		Cookies []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"cookies"`
	} `json:"cookie_info"`
}

// Login 登陆
func (b *BService) Login(relogin bool) error {
	if err := b.getLoginInfo(relogin); err != nil {
		return err
	}

	if err := SaveCookieToFile(&b.loginInfo, cookieFileName); err != nil {
		b.logger.Println("本地cookie保存失败")
	}

	if err := b.getCurrentUser(); err != nil {
		return errors.New("获取用户信息失败, 请检查账号密码是否输入正确")
	}

	b.logger.Printf("你好呀 %s\n", b.user.Name)

	return nil
}

func (b *BService) getLoginInfo(relogin bool) error {
	if relogin {
		fmt.Print("输入账号: ")
		fmt.Scan(&b.loginInfo.Username)
		fmt.Print("输入密码: ")
		fmt.Scan(&b.loginInfo.Password)

		encryptPw, err := b.getEncryptPw([]byte(b.loginInfo.Password))
		if err != nil {
			return err
		}
		params := map[string]string{
			"appkey":   appKey,
			"password": encryptPw,
			"username": b.loginInfo.Username,
		}
		resp, err := b.POST(apiURL["login"], params, nil)
		if err != nil {
			return err
		}
		var bresp struct {
			Data CookieData `json:"data"`
		}
		if err := JSONProc(resp, &bresp); err != nil {
			return errors.New("登陆失败, 请检查账号密码是否输入正确")
		}
		loginInfo := ParseCookies(&bresp.Data)
		b.loginInfo.AccessKey = loginInfo.AccessKey
		b.loginInfo.Csrf = loginInfo.Csrf
		b.loginInfo.Cookies = loginInfo.Cookies
		b.loginInfo.Headers = loginInfo.Headers
		b.loginInfo.UID = loginInfo.UID
	} else {
		cookieInfo, err := LoadCookieFromFile(cookieFileName)
		if err != nil {
			return b.getLoginInfo(true)
		}
		b.loginInfo = DeParseCookies(cookieInfo)
		b.logger.Println("读取本地cookie成功")
	}
	return nil
}

func (b *BService) getEncryptPw(data []byte) (string, error) {
	params := map[string]string{
		"appkey": appKey,
	}
	resp, err := b.POST(apiURL["getKey"], params, nil)
	if err != nil {
		return "", err
	}
	var bresp struct {
		Data struct {
			Hash string `json:"hash"`
			Key  string `json:"key"`
		} `json:"data"`
	}
	if err := JSONProc(resp, &bresp); err != nil {
		return "", err
	}
	hash := bresp.Data.Hash
	key := bresp.Data.Key
	encrypt := RsaEncrypt(append([]byte(hash), data...), key)
	return base64.URLEncoding.EncodeToString(encrypt), nil
}
