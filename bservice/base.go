package bservice

import (
	"log"
	"net/http"
	"os"
	"sync"
)

const (
	appKey         = "1d8b6e7d45233436"
	cookieFileName = "cookie"
)

// var apiURL = map[string]string{
// 	"login":           "https://passport.bilibili.com/api/v2/oauth2/login",
// 	"getKey":          "https://passport.bilibili.com/api/oauth2/getKey",
// 	"share":           "https://app.bilibili.com/x/v2/view/share/add",
// 	"following":       "https://api.bilibili.com/x/relation/followings",
// 	"getSubmitVideos": "https://space.bilibili.com/ajax/member/getSubmitVideos",
// 	"watchAv":         "https://api.bilibili.com/x/report/web/heartbeat",
// 	"reward":          "https://account.bilibili.com/home/reward",
// 	"giveCoin":        "https://api.bilibili.com/x/web-interface/coin/add",
// 	"view":            "https://api.bilibili.com/x/web-interface/view",
// 	"getInfo":         "https://space.bilibili.com/ajax/member/GetInfo",
// }

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
}

// BService 基础的服务
type BService struct {
	client    *http.Client
	loginInfo LoginInfo
	videoList []float64
	logger    *log.Logger
	user      UserInfo
	urls      BURL
}

// Init 初始化服务
func (b *BService) Init() {
	b.client = &http.Client{}
	b.logger = log.New(os.Stderr, "[BiliBili-Tools] ", log.LstdFlags)
	b.setURLs()
}

// LoadVideoInfo 读取视频列表
func (b *BService) LoadVideoInfo(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		b.loadVideoList()
		WaitHours(12)
	}
}

func (b *BService) setURLs() {
	b.urls = BURL{
		Login:        "https://passport.bilibili.com/api/v2/oauth2/login",
		Share:        "https://app.bilibili.com/x/v2/view/share/add",
		WatchAv:      "https://api.bilibili.com/x/report/web/heartbeat",
		GiveCoin:     "https://api.bilibili.com/x/web-interface/coin/add",
		EncryptKey:   "https://passport.bilibili.com/api/oauth2/getKey",
		Following:    "https://api.bilibili.com/x/relation/followings",
		SubmitVideos: "https://space.bilibili.com/ajax/member/getSubmitVideos",
		Reward:       "https://account.bilibili.com/home/reward",
		VideoView:    "https://api.bilibili.com/x/web-interface/view",
		UserInfo:     "https://space.bilibili.com/ajax/member/GetInfo",
		UnreadCount:  "https://api.bilibili.com/x/feed/unread/count?type=0",
		Dynamic:      "https://api.bilibili.com/x/feed/pull?ps=1&type=0",
	}
}
