package main

import (
	"fmt"
	"sync"

	"github.com/XDfield/go-bilibili-tools/bservice"
)

const (
	version = "1.0.5"
)

func main() {
	fmt.Println("版本: " + version)

	bservice := bservice.BService{}
	bservice.Init()
	// 登陆
	if err := bservice.Login(false); err != nil {
		fmt.Printf("%v", err)
		return
	}
	// 启动服务
	wg := sync.WaitGroup{}
	wg.Add(4)
	go bservice.LoadVideoInfo()     // 半天读取一次视频列表
	go bservice.ShareService(&wg)   // 分享
	go bservice.WatchService(&wg)   // 观看视频
	go bservice.CoinService(&wg)    // 投币
	go bservice.DynamicService(&wg) // 关注推送
	wg.Wait()
	fmt.Println("退出")
}
