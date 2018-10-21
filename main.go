package main

import (
	"fmt"
	"sync"

	"github.com/XDfield/go-bilibili-tools/bservice"
)

func main() {
	bservice := bservice.BService{}
	bservice.Init()
	// 登陆
	if err := bservice.Login(false); err != nil {
		fmt.Printf("%v", err)
		return
	}
	// 默认评论
	bservice.Replays = []string{"第一??", "(=・ω・=)", "emmmm"}
	// 启动服务
	wg := sync.WaitGroup{}
	wg.Add(5)
	go bservice.LoadVideoInfo(&wg)  // 半天读取一次视频列表
	go bservice.ShareService(&wg)   // 分享
	go bservice.WatchService(&wg)   // 观看视频
	go bservice.CoinService(&wg)    // 投币
	go bservice.DynamicService(&wg) // 关注推送
	wg.Wait()
	fmt.Println("退出")
}
