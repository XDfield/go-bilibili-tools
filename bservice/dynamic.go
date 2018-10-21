package bservice

import (
	"strconv"
	"sync"
)

// DynamicService 关注推送服务
func (b *BService) DynamicService(wg *sync.WaitGroup) {
	b.logger.Println("启动关注推送服务")
	defer wg.Done()
	for {
		if err := b.showDynamic(); err != nil {
			b.logger.Printf("%v\n", err)
			continue
		}
		WaitSeconds(60)
	}
}

func (b *BService) showDynamic() error {
	unread, err := b.getUnreadCount()
	if err != nil {
		return err
	}
	if unread <= 0 {
		return nil
	}
	var bresp struct {
		Code int `json:"code"`
		Data struct {
			Feeds []struct {
				AddID    int `json:"add_id"`
				Addition struct {
					AID       int    `json:"aid"`
					MID       int    `json:"mid"`
					Author    string `json:"author"`
					Coins     int    `json:"coins"`
					Create    string `json:"create"`
					Desc      string `json:"description"`
					Duration  string `json:"duration"`
					Favorites int    `json:"favorites"`
					Link      string `json:"link"`
					Play      int    `json:"play"`
					Review    int    `json:"review"`
					Status    int    `json:"status"`
					Title     string `json:"title"`
					SubTitle  string `json:"subtitle"`
					TypeID    int    `json:"typeid"`
					TypeName  string `json:"typename"`
				} `json:"addition"`
				CTime  int `json:"ctime"`
				ID     int `json:"id"`
				MCID   int `json:"mcid"`
				Source struct {
					MID   string `json:"mid"`
					UName string `json:"uname"`
				} `json:"source"`
				SrcID int `json:"src_id"`
				Type  int `json:"type"`
			} `json:"feeds"`
			Page struct {
				Count int `json:"count"`
				Num   int `json:"num"`
				Size  int `json:"size"`
			} `json:"page"`
		} `json:"data"`
	}
	if err := b.client.GetAndDecode(b.urls.Dynamic, nil, b.loginInfo.Headers, &bresp); err != nil {
		return err
	}
	if len(bresp.Data.Feeds) > 0 {
		content := bresp.Data.Feeds[0]
		message := content.Addition.Author + " 在" + content.Addition.Create + "更新了《" + content.Addition.Title + "》"
		b.logger.Println(message)
		if err := b.replay("(=・ω・=)", strconv.Itoa(content.Addition.AID)); err != nil {
			return err
		}
		b.logger.Println("评论发送成功")
	}

	return nil
}
