package bservice

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
)

func (b *BService) getRandAid() (string, error) {
	videoList, err := b.getSubmitVideo()
	if err != nil {
		return "", err
	}
	return Float64ToString(videoList[rand.Intn(len(videoList))]), nil
}

func (b *BService) queryReward() ([]bool, int, error) {
	headers := b.loginInfo.Headers
	headers["Referer"] = "https://account.bilibili.com/account/home"
	resp, err := b.GET(apiURL["reward"], nil, headers)
	if err != nil {
		return nil, 0, err
	}
	var bresp struct {
		Data struct {
			Login   bool `json:"login"`
			WatchAv bool `json:"watch_av"`
			ShareAv bool `json:"share_av"`
			CoinsAv int  `json:"coins_av"`
		} `json:"data"`
	}
	if err := JSONProc(resp, &bresp); err != nil {
		return nil, 0, err
	}
	return []bool{bresp.Data.Login, bresp.Data.WatchAv, bresp.Data.ShareAv}, bresp.Data.CoinsAv, nil
}

func (b *BService) getCid(aid string) (string, error) {
	params := map[string]string{
		"aid": aid,
	}
	resp, err := b.GET(apiURL["getPageList"], params, nil)
	if err != nil {
		return "", err
	}
	var bresp []struct {
		Cid float64 `json:"cid"`
	}
	if err := JSONProc(resp, &bresp); err != nil {
		return "", err
	}
	return Float64ToString(bresp[0].Cid), nil
}

func (b *BService) getAttention() ([]float64, error) {
	attentionList := make([]float64, 0, 50)
	params := map[string]string{
		"vmid":  b.loginInfo.UID,
		"ps":    "50",
		"order": "desc",
	}
	resp, err := b.GET(apiURL["following"], params, b.loginInfo.Headers)
	if err != nil {
		return nil, err
	}
	var bresp struct {
		Data struct {
			List []struct {
				Mid float64 `json:"mid"`
			} `json:"list"`
		} `json:"data"`
	}
	if err := JSONProc(resp, &bresp); err != nil {
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
		resp, err := b.GET(apiURL["getSubmitVideos"], params, nil)
		if err != nil {
			fmt.Printf("%v", err)
			continue
		}
		var bresp struct {
			Data struct {
				Vlist []struct {
					Aid float64 `json:"aid"`
				} `json:"vlist"`
			} `json:"data"`
		}
		if err := JSONProc(resp, &bresp); err != nil {
			fmt.Printf("%v", err)
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
