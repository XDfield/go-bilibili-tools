package bservice

import (
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type stateStruct struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// GetCurrentTime 返回时间戳
func GetCurrentTime() string {
	return strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
}

// ParseCookies 解析cookie响应
func ParseCookies(cookieData *CookieData) LoginInfo {
	var loginInfo LoginInfo
	var cookieFormat string
	for _, cookie := range cookieData.CookieInfo.Cookies {
		cookieFormat += cookie.Name + "=" + cookie.Value + ";"
		if cookie.Name == "bili_jct" {
			loginInfo.Csrf = cookie.Value
		}
		if cookie.Name == "DedeUserID" {
			loginInfo.UID = cookie.Value
		}
	}
	loginInfo.Cookies = cookieFormat
	loginInfo.Headers = map[string]string{
		"Host":            "api.bilibili.com",
		"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
		"Cookie":          cookieFormat,
	}
	loginInfo.AccessKey = cookieData.TokenInfo.AccessToken
	return loginInfo
}

// DeParseCookies 将本地读取的cookie内容转为loginInfo
func DeParseCookies(cookieInfo map[string]string) LoginInfo {
	var loginInfo LoginInfo
	loginInfo.Username = cookieInfo["username"]
	loginInfo.Password = cookieInfo["password"]
	loginInfo.Cookies = cookieInfo["cookies"]
	loginInfo.Headers = map[string]string{
		"Host":            "api.bilibili.com",
		"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
		"Cookie":          cookieInfo["cookies"],
	}
	loginInfo.AccessKey = cookieInfo["accessKey"]
	loginInfo.Csrf = regexp.MustCompile(`bili_jct=(.*?);`).FindAllStringSubmatch(cookieInfo["cookies"], 1)[0][1]
	loginInfo.UID = regexp.MustCompile(`DedeUserID=(.*?);`).FindAllStringSubmatch(cookieInfo["cookies"], 1)[0][1]
	return loginInfo
}

// CheckCode 检查状态响应
// func CheckCode(resp *http.Response) error {
// 	state := stateStruct{}
// 	if err := JSONProc(resp, &state); err != nil {
// 		return err
// 	}
// 	if state.Code != 0 {
// 		return errors.New(state.Message)
// 	}
// 	return nil
// }

// WaitHours 等待 h 小时
func WaitHours(h int) {
	time.Sleep(time.Duration(h) * time.Hour)
}

// WaitSeconds 等待 s 秒
func WaitSeconds(s int) {
	time.Sleep(time.Duration(s) * time.Second)
}

// Float64ToString float64转字符串
func Float64ToString(value float64) string {
	return strconv.FormatInt(int64(value), 10)
}

// RsaEncrypt RSA加密
func RsaEncrypt(data []byte, publickey string) (encrypt []byte) {
	block, _ := pem.Decode([]byte(publickey))
	pubInterface, _ := x509.ParsePKIXPublicKey(block.Bytes)
	encrypt, _ = rsa.EncryptPKCS1v15(crand.Reader, pubInterface.(*rsa.PublicKey), data)
	return
}

// SaveCookieToFile 保存cookie到本地文件
func SaveCookieToFile(loginInfo *LoginInfo, filename string) error {
	buffer := make([]string, 0, 4)
	buffer = append(buffer, "username: "+loginInfo.Username)
	buffer = append(buffer, "password: "+loginInfo.Password)
	buffer = append(buffer, "cookies: "+loginInfo.Cookies)
	buffer = append(buffer, "accessKey: "+loginInfo.AccessKey)

	contents := strings.Join(buffer, "\n")
	fileObj, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer fileObj.Close()
	if err != nil {
		return fmt.Errorf("<SaveCookieToFile>: %v", err)
	}
	if _, err := io.WriteString(fileObj, contents); err != nil {
		return fmt.Errorf("<SaveCookieToFile>: %v", err)
	}
	return nil
}

// LoadCookieFromFile 读取本地cookie文件
func LoadCookieFromFile(filename string) (map[string]string, error) {
	fileObj, err := os.Open(filename)
	defer fileObj.Close()
	if err != nil {
		return nil, fmt.Errorf("<LoadCookieFromFile>: %v", err)
	}
	contents, err := ioutil.ReadAll(fileObj)
	if err != nil {
		return nil, fmt.Errorf("<LoadCookieFromFile>: %v", err)
	}
	values := strings.Split(string(contents), "\n")
	cookieInfo := make(map[string]string)
	for _, val := range values {
		t := strings.Split(val, ": ")
		cookieInfo[t[0]] = t[1]
	}
	return cookieInfo, nil
}
