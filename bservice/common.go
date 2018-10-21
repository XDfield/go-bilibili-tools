package bservice

import (
	"errors"
	"math/rand"
	"strconv"
)

func (b *BService) getRandAid() (string, error) {
	videoList := b.videoList
	for ; len(videoList) == 0; videoList = b.videoList {
		WaitSeconds(2)
	}
	return Float64ToString(videoList[rand.Intn(len(videoList))]), nil
}

func (b *BService) replay(message string, mid string) error {
	data := map[string]string{
		"oid":     mid,
		"type":    "1",
		"message": message,
		"plat":    "1",
		"jsonp":   "jsonp",
		"csrf":    b.loginInfo.Csrf,
	}
	headers := b.loginInfo.Headers
	headers["Referer"] = "https://www.bilibili.com/video/av" + mid
	var bresp struct {
		Code int `json:"code"`
		Data struct {
			RPID    int    `json:"rpid"`
			RPIDStr string `json:"rpid_str"`
		} `json:"data"`
		Message string `json:"message"`
	}
	if err := b.client.PostAndDecode(b.urls.Replay, data, headers, &bresp); err != nil {
		return err
	}
	if bresp.Code != 0 {
		return errors.New("评论发送失败")
	}
	return nil
}

func (b *BService) loadVideoList() {
	videoList, err := b.getSubmitVideo()
	if err != nil {
		b.loadVideoList()
	} else {
		b.videoList = videoList
	}
}

func (b *BService) getView(aid string) (*videoView, error) {
	params := map[string]string{
		"aid": aid,
	}
	view := videoView{}
	if err := b.client.GetAndDecode(b.urls.VideoView, params, nil, &view); err != nil {
		return nil, err
	}
	return &view, nil
}

func (b *BService) getCurrentUser() error {
	data := map[string]string{
		"mid":  b.loginInfo.UID,
		"csrf": b.loginInfo.Csrf,
	}
	headers := b.loginInfo.Headers
	headers["Referer"] = "https://space.bilibili.com/3213445" + b.loginInfo.UID
	var bresp struct {
		Status bool     `json:"status"`
		Data   UserInfo `json:"data"`
	}
	if err := b.client.PostAndDecode(b.urls.UserInfo, data, headers, &bresp); err != nil {
		return err
	}
	b.user = bresp.Data
	return nil
}

func (b *BService) queryReward() ([]bool, int, error) {
	headers := b.loginInfo.Headers
	headers["Referer"] = "https://account.bilibili.com/account/home"
	var bresp struct {
		Data struct {
			Login   bool `json:"login"`
			WatchAv bool `json:"watch_av"`
			ShareAv bool `json:"share_av"`
			CoinsAv int  `json:"coins_av"`
		} `json:"data"`
	}
	if err := b.client.GetAndDecode(b.urls.Reward, nil, headers, &bresp); err != nil {
		return nil, 0, err
	}
	return []bool{bresp.Data.Login, bresp.Data.WatchAv, bresp.Data.ShareAv}, bresp.Data.CoinsAv, nil
}

func (b *BService) getUnreadCount() (int, error) {
	var bresp struct {
		Code int `json:"code"`
		Data struct {
			All int `json:"all"`
		} `json:"data"`
	}
	if err := b.client.GetAndDecode(b.urls.UnreadCount, nil, b.loginInfo.Headers, &bresp); err != nil {
		return 0, err
	}
	return bresp.Data.All, nil
}

func (b *BService) getAttention() ([]float64, error) {
	attentionList := make([]float64, 0, 50)
	params := map[string]string{
		"vmid":  b.loginInfo.UID,
		"ps":    "50",
		"order": "desc",
	}
	var bresp struct {
		Data struct {
			List []struct {
				Mid float64 `json:"mid"`
			} `json:"list"`
		} `json:"data"`
	}
	if err := b.client.GetAndDecode(b.urls.Following, params, b.loginInfo.Headers, &bresp); err != nil {
		return nil, err
	}
	for _, val := range bresp.Data.List {
		attentionList = append(attentionList, val.Mid)
	}
	return attentionList, nil
}

func (b *BService) getSubmitVideo() ([]float64, error) {
	attentionList, err := b.getAttention()
	if err != nil {
		return nil, err
	}
	videoList := make([]float64, 0, 50)
	for _, mid := range attentionList {
		params := map[string]string{
			"mid":      strconv.FormatInt(int64(mid), 10),
			"pagesize": "100",
			"tid":      "0",
		}
		var bresp struct {
			Data struct {
				Vlist []struct {
					Aid float64 `json:"aid"`
				} `json:"vlist"`
			} `json:"data"`
		}
		if err := b.client.GetAndDecode(b.urls.SubmitVideos, params, nil, &bresp); err != nil {
			continue
		}
		for _, val := range bresp.Data.Vlist {
			videoList = append(videoList, val.Aid)
		}
	}
	if len(videoList) == 0 {
		return nil, errors.New("获取视频信息失败")
	}
	return videoList, nil
}
