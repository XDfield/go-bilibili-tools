package bservice

import (
	"fmt"
	"sync"
)

type failState struct {
	State bool   `json:"state"`
	Data  string `json:"data"`
}

// ShareService 视频分享服务
func (b *BService) ShareService(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		aid, err := b.getRandAid()
		if err != nil {
			fmt.Printf("%v", err)
			continue
		}
		if err := b.share(aid); err != nil {
			fmt.Printf("分享视频失败: %v\n", err)
			continue
		} else {
			view, err := b.getView(aid)
			if err != nil {
				fmt.Printf("获取视频信息失败: av%v\n", aid)
			} else {
				fmt.Printf("成功分享视频: (av%v) %v\n", aid, view.Data.Title)
			}

		}
		fmt.Println("分享任务完成, 六小时后继续")

		WaitHours(6)
	}
}

func (b *BService) share(aid string) error {
	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 BiliDroid/5.26.3 (bbcallen@gmail.com)",
		"Host":       "app.bilibili.com",
		"Cookie":     "sid=8wfvu7i7",
	}
	data := map[string]string{
		"access_key": b.loginInfo.AccessKey,
		"aid":        aid,
		"appkey":     appKey,
		"build":      "5260003",
		"from":       "7",
		"mobi_app":   "android",
		"platform":   "android",
		"ts":         GetCurrentTime(),
	}
	resp, err := b.POST(apiURL["share"], data, headers)
	if err != nil {
		return err
	}
	return CheckCode(resp)
}
