package bservice

import (
	"encoding/json"
	"io/ioutil"
)

// Config 配置文件
type Config struct {
	ShareServerEnable     bool   `json:"ShareServerEnable"`
	WatchServerEnable     bool   `json:"WatchServerEnable"`
	CoinServerEnable      bool   `json:"CoinServerEnable"`
	DynamicServerEnable   bool   `json:"DynamicServerEnable"`
	BarkKey               string `json:"BarkKey"`
	DynamicCheckTime      int    `json:"DynamicCheckTime"`
	DefaultReplay         string `json:"DefaultReplay"`
	OnlySpecialAttentions bool   `json:"OnlySpecialAttentions"`
	SpecialAttentions     []struct {
		MID    int    `json:"mid"`
		Replay string `json:"replay"`
	} `json:"SpecialAttentions"`
}

func (b *BService) parseConfig(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, b.config); err != nil {
		return err
	}
	return nil
}
