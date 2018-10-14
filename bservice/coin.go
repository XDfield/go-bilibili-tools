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
			for err := b.giveCoin(); err == nil; {
				fmt.Println("投币失败...")
			}
		}
		fmt.Println("投币任务完成")
		WaitHours(24)
	}
}

func (b *BService) giveCoin() error {
	aid, err := b.getRandAid()
	if err != nil {
		return err
	}
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
	var bresp struct {
		Code int `json:"code"`
	}
	if err := JSONProc(resp, &bresp); err != nil {
		return err
	}
	if bresp.Code != 0 {
		return b.giveCoin()
	}

	fmt.Printf("投币成功: aid: %v\n", aid)
	return nil
}
