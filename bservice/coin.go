package bservice

import (
	"fmt"
	"sync"
)

// CoinService 投币服务
func (b *BService) CoinService(wg *sync.WaitGroup) {
	defer wg.Done()
	if !b.config.CoinServerEnable {
		return
	}
	b.logger.Println("启动投币服务")
	for {
		_, coinExp, err := b.queryReward()
		if err != nil {
			b.logger.Printf("<CoinService>: %v", err)
			continue
		}
		for ; coinExp < 50; coinExp += 10 {
			for {
				aid, err := b.getRandAid()
				if err != nil {
					b.logger.Printf("<CoinService>: %v\n", err)
					continue
				}
				if err := b.giveCoin(aid); err != nil {
					b.logger.Printf("<CoinService>: %v\n", err)
					continue
				}
				if view, err := b.getView(aid); err == nil {
					b.logger.Printf("成功投币: (av%v) %v\n", aid, view.Data.Title)
				} else {
					b.logger.Printf("获取视频信息失败: av%v\n", aid)
				}

				break
			}
		}
		b.logger.Println("今日投币任务完成, 二十四小时后继续")

		WaitHours(24)
	}
}

func (b *BService) giveCoin(aid string) error {
	headers := b.loginInfo.Headers
	headers["Referer"] = "https://www.bilibili.com/video/av" + aid
	data := map[string]string{
		"aid":          aid,
		"multiply":     "1",
		"cross_domain": "true",
		"csrf":         b.loginInfo.Csrf,
	}
	state := stateStruct{}
	if err := b.client.PostAndDecode(b.urls.GiveCoin, data, headers, &state); err != nil {
		return fmt.Errorf("<giveCoin>: %v", err)
	}
	if state.Code != 0 {
		return fmt.Errorf("<giveCoin>: %s", state.Message)
	}
	return nil
}
