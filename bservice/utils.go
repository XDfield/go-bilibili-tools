package bservice

import (
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

// GetCurrentTime 返回时间戳
func GetCurrentTime() string {
	return strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
}

// WaitHours 等待 h 小时
func WaitHours(h int) {
	time.Sleep(time.Duration(h*3600) * time.Second)
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
func SaveCookieToFile(loginInfo LoginInfo, filename string) error {
	buffer := make([]string, 0, 4)
	buffer = append(buffer, "username: "+loginInfo.Username)
	buffer = append(buffer, "password: "+loginInfo.Password)
	buffer = append(buffer, "cookies: "+loginInfo.Cookies)
	buffer = append(buffer, "accessKey: "+loginInfo.AccessKey)

	contents := strings.Join(buffer, "\n")
	fileObj, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer fileObj.Close()
	if err != nil {
		return err
	}
	if _, err := io.WriteString(fileObj, contents); err != nil {
		return err
	}
	return nil
}

// LoadCookieFromFile 读取本地cookie文件
func LoadCookieFromFile(filename string) (map[string]string, error) {
	fileObj, err := os.Open(filename)
	defer fileObj.Close()
	if err != nil {
		return nil, err
	}
	contents, err := ioutil.ReadAll(fileObj)
	if err != nil {
		return nil, err
	}
	values := strings.Split(string(contents), "\n")
	cookieInfo := make(map[string]string)
	for _, val := range values {
		t := strings.Split(val, ": ")
		cookieInfo[t[0]] = t[1]
	}
	return cookieInfo, nil
}
