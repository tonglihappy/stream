package util

import (
	"bytes"
	"encoding/json"
	//"fmt"
	"net/http"
)

type DingText struct {
	Content string `json:"content"`
}

type DingMsg struct {
	MsgType string   `json:"msgtype"`
	Text    DingText `json:"text"`
}

func DingNotify(content string) {
	return
	url := "https://oapi.dingtalk.com/robot/send?access_token=33c0343f251f4e976362b7995867217c71a5d0fb4f79ac96e50f80dd0fd15190"
	var msg = DingMsg{}
	msg.MsgType = "text"
	msg.Text.Content = content

	b, _ := json.Marshal(msg)
	http.Post(url, "application/json; charset=utf-8", bytes.NewBuffer(b))
}
