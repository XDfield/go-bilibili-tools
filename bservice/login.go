package bservice

import (
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
)

// Login 登陆
func (b *BService) Login(relogin bool) error {
	if !relogin {
		if err := b.loadCookie(); err == nil {
			fmt.Println("读取本地cookie成功")
			return nil
		}
	}

	if b.loginInfo.Username == "" {
		fmt.Print("输入账号: ")
		fmt.Scan(&b.loginInfo.Username)
		fmt.Print("输入密码: ")
		fmt.Scan(&b.loginInfo.Password)
	}

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
	bresp := BResponse{}
	if err := JSONProc(resp, &bresp); err != nil {
		return err
	}
	if bresp.Data == nil {
		return errors.New("登陆失败, 请检查账号密码是否输入正确")
	}
	err = b.saveCookie(bresp.Data)
	if err != nil {
		fmt.Println("本地cookie保存失败")
	}
	return nil
}

func (b *BService) saveCookie(data map[string]interface{}) error {
	cookies := data["cookie_info"].(map[string]interface{})["cookies"].([]interface{})
	var cookieFormat, name, value string
	var cookie map[string]interface{}
	for i := 0; i < len(cookies); i++ {
		cookie = cookies[i].(map[string]interface{})
		name = cookie["name"].(string)
		value = cookie["value"].(string)
		cookieFormat += name + "=" + value + ";"
		if name == "bili_jct" {
			b.loginInfo.Csrf = value
		}
		if name == "DedeUserID" {
			b.loginInfo.UID = value
		}
	}
	b.loginInfo.Cookies = cookieFormat
	b.loginInfo.Headers = map[string]string{
		"Host":            "api.bilibili.com",
		"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
		"Cookie":          cookieFormat,
	}
	b.loginInfo.AccessKey = data["token_info"].(map[string]interface{})["access_token"].(string)
	if err := SaveCookieToFile(b.loginInfo, cookieFileName); err != nil {
		return err
	}
	return nil
}

func (b *BService) loadCookie() error {
	cookieInfo, err := LoadCookieFromFile(cookieFileName)
	if err != nil {
		return err
	}
	b.loginInfo.Username = cookieInfo["username"]
	b.loginInfo.Password = cookieInfo["password"]
	b.loginInfo.Cookies = cookieInfo["cookies"]
	b.loginInfo.Headers = map[string]string{
		"Host":            "api.bilibili.com",
		"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
		"Cookie":          cookieInfo["cookies"],
	}
	b.loginInfo.AccessKey = cookieInfo["accessKey"]
	b.loginInfo.Csrf = regexp.MustCompile(`bili_jct=(.*?);`).FindAllStringSubmatch(cookieInfo["cookies"], 1)[0][1]
	b.loginInfo.UID = regexp.MustCompile(`DedeUserID=(.*?);`).FindAllStringSubmatch(cookieInfo["cookies"], 1)[0][1]
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
	bresp := BResponse{}
	if err := JSONProc(resp, &bresp); err != nil {
		return "", err
	}
	hash := bresp.Data["hash"].(string)
	key := bresp.Data["key"].(string)
	encrypt := RsaEncrypt(append([]byte(hash), data...), key)
	return base64.URLEncoding.EncodeToString(encrypt), nil
}
