package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"koi/pkg/gocqhttp"
	"koi/pkg/gocqhttp/event"
	"koi/pkg/log"

	"github.com/gorilla/websocket"
)

var (
	URL = flag.String("url", "ws://localhost:3020", "gocqhttp websocket url")
)

func init() {}

func main() {
	var websocket websocket.Dialer

	log.Info("当前版本: dev-0.0.1")
	log.Info("正在连接到go-cqhttp...")

	ws_url := *URL + "/api"
	ws_api, _, err := websocket.Dial(ws_url, nil)
	if err != nil {
		panic(err)
	}

	log.Info("WebSocket 连接:", ws_url)

	ws_url = *URL + "/event"
	ws_event, _, err := websocket.Dial(ws_url, nil)
	if err != nil {
		panic(err)
	}

	log.Info("WebSocket 连接:", ws_url)

	// 输出登陆号信息
	login, err := gocqhttp.GetLoginInfo(ws_api)
	if err != nil {
		log.Warn("未获取到登录信息, 请检查go-cqhttp控制台输出")
	} else {
		log.Info("欢迎使用:", login.Nickname)
	}

	// 读取事件
	for {
		Listen(ws_api, ws_event)
	}
}

// 读取事件
func Listen(ws_api, ws_event *websocket.Conn) {
	var event map[string]any

	_, data, err := ws_event.ReadMessage()
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &event)
	if err != nil {
		panic(err)
	}

	// 启用协程, 避免阻塞
	go Switch(ws_api, data, event)
}

// 事件分支
func Switch(ws_api *websocket.Conn, data []byte, event map[string]any) {
	types := event["post_type"].(string)

	// 通过"post_type"判断事件类型
	switch types {
	case "meta_event":
		types = event["meta_event_type"].(string)
		meta(ws_api, types, data)
	case "message":
		types = event["message_type"].(string)
		message(ws_api, types, data)
	case "request":
		types = event["request_type"].(string)
		request(ws_api, types, data)
	case "notice":
		notice(ws_api, data, event)
	default:
		HandlerUnknown(ws_api, data)
	}
}

// 元事件
func meta(ws_api *websocket.Conn, types string, data []byte) {
	switch types {
	case "lifecycle":
		HandlerLifecycle(ws_api, data)
	case "heartbeat":
		HandlerHeartbeat(ws_api, data)
	default:
		HandlerUnknown(ws_api, data)
	}
}

// 消息事件
func message(ws_api *websocket.Conn, types string, data []byte) {
	switch types {
	case "private":
		HandlerPrivateMessage(ws_api, data)
	case "group":
		HandlerGroupMessage(ws_api, data)
	default:
		HandlerUnknown(ws_api, data)
	}
}

// 请求事件
func request(ws_api *websocket.Conn, types string, data []byte) {
	switch types {
	case "friend":
		HandlerFriendRequest(ws_api, data)
	case "group":
		HandlerGroupRequest(ws_api, data)
	default:
		HandlerUnknown(ws_api, data)
	}
}

// 通知事件
func notice(ws_api *websocket.Conn, data []byte, event map[string]any) {
	types := event["post_type"].(string)

	switch types {
	case "group_admin":
		HandlerGroupAdmin(ws_api, data)
	case "group_ban":
		HandlerGroupBan(ws_api, data)
	case "group_card":
		HandlerGroupCard(ws_api, data)
	case "group_decrease":
		HandlerGroupDecrease(ws_api, data)
	case "group_increase":
		HandlerGroupIncrease(ws_api, data)
	case "group_recall":
		HandlerGroupRecall(ws_api, data)
	case "group_upload":
		HandlerGroupUpload(ws_api, data)
	case "friend_add":
		HandlerFriendAdd(ws_api, data)
	case "friend_recall":
		HandlerFriendRecall(ws_api, data)
	case "notify":
		types = event["sub_type"].(string)
		notify(ws_api, types, data)
	case "essence":
		HandlerEssence(ws_api, data)
	case "offline_file":
		HandlerOfflineFile(ws_api, data)
	case "client_status":
		HandlerClientStatus(ws_api, data)
	default:
		HandlerUnknown(ws_api, data)
	}
}

// 通知事件
func notify(ws_api *websocket.Conn, types string, data []byte) {
	switch types {
	case "poke":
		HandlerPoke(ws_api, data)
	case "lucky_king":
		HandlerLuckyKing(ws_api, data)
	case "honor":
		HandlerHonor(ws_api, data)
	}
}

// 未知事件
func HandlerUnknown(_ *websocket.Conn, data []byte) {
	fmt.Println(string(data))
}

// * 元事件
//
// 生命周期
func HandlerLifecycle(ws_api *websocket.Conn, data []byte) {}

// * 元事件
//
// 心跳
func HandlerHeartbeat(ws_api *websocket.Conn, data []byte) {}

// * 消息事件
//
// 私聊消息
func HandlerPrivateMessage(ws_api *websocket.Conn, data []byte) {
	// 私聊消息复读示例
	var message *event.PrivateMessage

	err := json.Unmarshal(data, &message)
	if err != nil {
		panic(err)
	}

	_, err = gocqhttp.SendPrivateMessage(ws_api, message.UserID, message.RawMessage)
	if err != nil {
		panic(err)
	}
}

// * 消息事件
//
// 群消息
func HandlerGroupMessage(ws_api *websocket.Conn, data []byte) {}

// * 请求事件
//
// 加好友请求
func HandlerFriendRequest(ws_api *websocket.Conn, data []byte) {}

// * 请求事件
//
// 加群请求/邀请
func HandlerGroupRequest(ws_api *websocket.Conn, data []byte) {}

// * 通知事件
//
// 群管理员变动
func HandlerGroupAdmin(ws_api *websocket.Conn, data []byte) {}

// * 通知事件
//
// 群禁言
func HandlerGroupBan(ws_api *websocket.Conn, data []byte) {}

// * 通知事件
//
// 群成员名片更新
func HandlerGroupCard(ws_api *websocket.Conn, data []byte) {}

// * 通知事件
//
// 群成员减少
func HandlerGroupDecrease(ws_api *websocket.Conn, data []byte) {}

// * 通知事件
//
// 群成员增加
func HandlerGroupIncrease(ws_api *websocket.Conn, data []byte) {}

// * 通知事件
//
// 群消息撤回
func HandlerGroupRecall(ws_api *websocket.Conn, data []byte) {}

// * 通知事件
//
// 群文件上传
func HandlerGroupUpload(ws_api *websocket.Conn, data []byte) {}

// * 通知事件
//
// 好友添加
func HandlerFriendAdd(ws_api *websocket.Conn, data []byte) {}

// * 通知事件
//
// 好友消息撤回
func HandlerFriendRecall(ws_api *websocket.Conn, data []byte) {}

// * 通知事件
//
// 戳一戳
func HandlerPoke(ws_api *websocket.Conn, data []byte) {}

// * 通知事件
//
// 群红包运气王
func HandlerLuckyKing(ws_api *websocket.Conn, data []byte) {}

// * 通知事件
//
// 群成员荣誉变更
func HandlerHonor(ws_api *websocket.Conn, data []byte) {}

// * 通知事件
//
// 精华消息
func HandlerEssence(ws_api *websocket.Conn, data []byte) {}

// * 通知事件
//
// 接收到离线文件
func HandlerOfflineFile(ws_api *websocket.Conn, data []byte) {}

// * 通知事件
//
// 其他客户端在线状态变更
func HandlerClientStatus(ws_api *websocket.Conn, data []byte) {}
