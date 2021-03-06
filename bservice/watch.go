package bservice

import (
	"fmt"
	"sync"
)

// WatchService 观看服务
func (b *BService) WatchService(wg *sync.WaitGroup) {
	defer wg.Done()
	if !b.config.WatchServerEnable {
		return
	}
	b.logger.Println("启动视频观看服务")
	for {
		aid, err := b.getRandAid()
		if err != nil {
			b.logger.Printf("<WatchService>: %v\n", err)
			continue
		}
		view, err := b.getView(aid)
		if err != nil {
			b.logger.Printf("<WatchService>: %v\n", err)
			continue
		}
		if err := b.watch(aid, string(view.Data.Cid)); err != nil {
			b.logger.Printf("<WatchService>: %v\n", err)
			continue
		}
		b.logger.Printf("成功观看视频: (av%v) %v\n", aid, view.Data.Title)

		b.logger.Println("观看任务完成, 六小时后继续")

		WaitHours(6)
	}
}

func (b *BService) watch(aid, cid string) error {
	headers := b.loginInfo.Headers
	headers["Referer"] = "https://www.bilibili.com/video/av" + aid
	data := map[string]string{
		"aid":         aid,
		"cid":         cid,
		"mid":         b.loginInfo.UID,
		"csrf":        b.loginInfo.Csrf,
		"played_time": "0",
		"realtime":    "0",
		"start_ts":    GetCurrentTime(),
		"type":        "3",
		"dt":          "2",
		"play_type":   "1",
	}
	state := stateStruct{}
	if err := b.client.PostAndDecode(b.urls.WatchAv, data, headers, &state); err != nil {
		return fmt.Errorf("<watch>: %v", err)
	}
	if state.Code != 0 {
		return fmt.Errorf("<watch>: %s", state.Message)
	}
	return nil
}
