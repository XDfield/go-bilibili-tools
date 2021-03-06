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
	if !b.config.ShareServerEnable {
		return
	}
	b.logger.Println("启动视频分享服务")
	for {
		aid, err := b.getRandAid()
		if err != nil {
			b.logger.Printf("<ShareService>: %v\n", err)
			continue
		}
		if err := b.share(aid); err != nil {
			b.logger.Printf("<ShareService>: %v\n", err)
			continue
		} else {
			view, err := b.getView(aid)
			if err != nil {
				b.logger.Printf("<ShareService>: %v\n", err)
			} else {
				b.logger.Printf("成功分享视频: (av%v) %v\n", aid, view.Data.Title)
			}

		}
		b.logger.Println("分享任务完成, 六小时后继续")

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
	state := stateStruct{}
	if err := b.client.PostAndDecode(b.urls.Share, data, headers, &state); err != nil {
		return fmt.Errorf("<share>: %v", err)
	}
	if state.Code != 0 {
		return fmt.Errorf("<share>: %s", state.Message)
	}
	return nil
}
