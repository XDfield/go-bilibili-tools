package bservice

import (
	"fmt"
	"sync"
)

// CoinService 投币服务
func (b *BService) CoinService(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		_, coinExp, err := b.queryReward()
		if err != nil {
			fmt.Printf("%v", err)
			continue
		}
		for ; coinExp < 50; coinExp += 10 {
			for {
				aid, err := b.getRandAid()
				if err != nil {
					fmt.Printf("获取aid失败: %v\n", err)
					continue
				}
				if err := b.giveCoin(aid); err != nil {
					fmt.Printf("投币失败: %v\n", err)
				}
				view, err := b.getView(aid)
				if err != nil {
					fmt.Printf("获取视频信息失败: av%v\n", aid)
				} else {
					fmt.Printf("成功投币: (av%v) %v\n", aid, view.Data.Title)
					break
				}
			}
		}
		fmt.Println("今日投币任务完成, 二十四小时后继续")

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
	resp, err := b.POST(apiURL["giveCoin"], data, headers)
	if err != nil {
		return err
	}
	return CheckCode(resp)
}
