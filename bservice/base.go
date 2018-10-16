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

var apiURL = map[string]string{
	"login":           "https://passport.bilibili.com/api/v2/oauth2/login",
	"getKey":          "https://passport.bilibili.com/api/oauth2/getKey",
	"caseObtain":      "http://api.bilibili.com/x/credit/jury/caseObtain",
	"share":           "https://app.bilibili.com/x/v2/view/share/add",
	"following":       "https://api.bilibili.com/x/relation/followings",
	"getSubmitVideos": "https://space.bilibili.com/ajax/member/getSubmitVideos",
	"getPageList":     "https://www.bilibili.com/widget/getPageList",
	"watchAv":         "https://api.bilibili.com/x/report/web/heartbeat",
	"reward":          "https://account.bilibili.com/home/reward",
	"giveCoin":        "https://api.bilibili.com/x/web-interface/coin/add",
	"view":            "https://api.bilibili.com/x/web-interface/view",
	"getInfo":         "https://space.bilibili.com/ajax/member/GetInfo",
}

// BService 基础的服务
type BService struct {
	client    *http.Client
	loginInfo LoginInfo
	videoList []float64
	logger    *log.Logger
	user      UserInfo
}

// Init 初始化服务
func (b *BService) Init() {
	b.client = &http.Client{}
	b.logger = log.New(os.Stderr, "[BiliBili-Tools] ", log.LstdFlags)
}

// LoadVideoInfo 读取视频列表
func (b *BService) LoadVideoInfo(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		b.loadVideoList()
		WaitHours(12)
	}
}
