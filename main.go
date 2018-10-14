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
	err := bservice.Login(false)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Println("登陆成功!")
	// 启动服务
	fmt.Println("启动服务: 分享 观看 投币")
	wg := sync.WaitGroup{}
	wg.Add(4)
	go bservice.LoadVideoInfo(&wg) // 半天读取一次视频列表
	go bservice.ShareService(&wg)  // 分享
	go bservice.WatchService(&wg)  // 观看视频
	go bservice.CoinService(&wg)   // 投币
	wg.Wait()
	fmt.Println("退出")
}
