package bservice

import (
	"log"
	"net/http"
)

// BURL urls
type BURL struct {
	Login        string
	Share        string
	WatchAv      string
	GiveCoin     string
	EncryptKey   string
	Following    string
	SubmitVideos string
	Reward       string
	VideoView    string
	UserInfo     string
	UnreadCount  string
	Dynamic      string
	Replay       string
}

// BService 基础的服务
type BService struct {
	client    *BClient
	loginInfo LoginInfo
	videoList []float64
	logger    *log.Logger
	config    *Config
	user      UserInfo
	urls      BURL
}

// BClient 处理请求的对象
type BClient struct {
	http.Client
}

// Config 配置文件
type Config struct {
	ShareServerEnable     bool   `json:"ShareServerEnable"`
	WatchServerEnable     bool   `json:"WatchServerEnable"`
	CoinServerEnable      bool   `json:"CoinServerEnable"`
	DynamicServerEnable   bool   `json:"DynamicServerEnable"`
	BarkKey               string `json:"BarkKey"`
	DynamicCheckTime      int    `json:"DynamicCheckTime"`
	DefaultReplay         string `json:"DefaultReplay"`
	OnlySpecialAttentions bool   `json:"OnlySpecialAttentions"`
	SpecialAttentions     []struct {
		MID    int    `json:"mid"`
		Replay string `json:"replay"`
	} `json:"SpecialAttentions"`
}

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

// **********
// Responses
// **********

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

// VideoView 视频信息
type VideoView struct {
	Code int `json:"code"`
	Data struct {
		Aid   int    `json:"aid"`
		Cid   int    `json:"cid"`
		Desc  string `json:"desc"`
		Owner struct {
			Name string `json:"name"`
			Mid  int    `json:"mid"`
		} `json:"owner"`
		Stat struct {
			Coin    int `json:"coin"`
			Danmuku int `json:"danmuku"`
			Dislike int `json:"dislike"`
			Like    int `json:"like"`
			Share   int `json:"share"`
			View    int `json:"view"`
		} `json:"stat"`
		Tid   int    `json:"tid"`
		Title string `json:"title"`
		TName string `json:"tname"`
	}
}

// UserInfo 用户信息
type UserInfo struct {
	Birthday  string `json:"birthday"`
	Im9Sign   string `json:"im9_sign"`
	LevelInfo struct {
		CurrentLevel int `json:"current_level"`
	} `json:"level_info"`
	MID     int    `json:"mid"`
	Name    string `json:"name"`
	Rank    int    `json:"rank"`
	RegTime int    `json:"regtime"`
	Sex     string `json:"sex"`
	Sign    string `json:"sign"`
	Vip     struct {
		VipStatus int `json:"vipStatus"`
		VipType   int `json:"vipType"`
	} `json:"vip"`
}
