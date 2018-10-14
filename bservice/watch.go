package bservice

import (
	"fmt"
	"sync"
)

// WatchService 观看服务
func (b *BService) WatchService(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		aid, err := b.getRandAid()
		if err != nil {
			fmt.Printf("获取aid失败: %v\n", err)
			continue
		}
		cid, err := b.getCid(aid)
		if err != nil {
			fmt.Printf("获取cid失败: %v\n", err)
			continue
		}
		if err := b.watch(aid, cid); err != nil {
			fmt.Printf("观看视频失败: %v\n", err)
			continue
		}

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
	resp, err := b.POST(apiURL["watchAv"], data, headers)
	if err != nil {
		return err
	}
	var bresp struct {
		Code int `json:"code"`
	}
	if err := JSONProc(resp, &bresp); err != nil {
		fmt.Printf("%v", err)
		return err
	}
	if bresp.Code == 0 {
		fmt.Printf("观看视频完成 aid: %v cid: %v\n", aid, cid)
	} else {
		fmt.Printf("观看视频失败 aid: %v cid: %v\n", aid, cid)
	}
	return nil
}
