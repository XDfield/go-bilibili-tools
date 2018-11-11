### Go-Bilibili-Tools

参考[这个项目](https://github.com/Dawnnnnnn/bilibili-tools)写的 golang 版本

#### 功能

- 每日投币
- 每日分享
- 每日观看
- 关注推送并评论
- ~~每日风纪委员投票~~ (这个是要风纪委员才能搞的吗??找不到入口索性先不做了)

#### 使用

先在同目录创建配置文件 config.json , 内容如下:

```json
{
  "ShareServerEnable": true,
  "WatchServerEnable": true,
  "CoinServerEnable": true,
  "DynamicServerEnable": true,
  "BarkKey": "",
  "DynamicCheckTime": 15,
  "DefaultReplay": "(=・ω・=)",
  "OnlySpecialAttentions": false,
  "SpecialAttentions": [
    {
      "mid": 11111,
      "replay": "指定评论内容"
    }
  ]
}
```

> ShareServerEnable ~ DynamicServerEnable: 是否开启服务
> BarkKey: Bark 推送 key 不填则不推送
> DynamicCheckTime: 关注更新检查间隔(单位秒)
> DefaultReplay: 默认评论内容
> OnlySpecialAttentions: 是否只评论特殊关注列表
> SpecialAttentions: 特殊关注列表
> mid: up 主的 mid
> replay: 评论内容

```bash
# 安装
go get github.com/XDfield/go-bilibili-tools
# 编译
cd go-bilibili-tools/
go install
# 运行(windows: go-bilibili-tools.exe)
./go-bilibili-tools
```

> 第一次使用会要求输入账号密码， 登陆成功后会保存 cookie 到同目录下
