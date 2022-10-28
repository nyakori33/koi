package gocqhttp

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

// API返回的内容
type ws_body struct {
	Body     any     `json:"data"`    // 数据主体
	Code     *string `json:"msg"`     // 错误代码
	HTTPCode int     `json:"retcode"` // HTTP状态码
	Status   string  `json:"status"`  // 响应状态
	Message  *string `json:"wording"` // 错误信息
}

// 符合gocqhttp规范的数据
type ws_data struct {
	Action string `json:"action"`           // API终结点
	Params any    `json:"params,omitempty"` // 参数
}

var values = make(map[string]any)

func (message ws_data) do(ws *websocket.Conn) ([]byte, error) {
	var data *ws_body

	err := ws.WriteJSON(message)
	if err != nil {
		return nil, err
	}

	err = ws.ReadJSON(&data)
	if err != nil {
		return nil, err
	}

	if data.Code != nil {
		return nil, fmt.Errorf("%s: %s", *data.Code, *data.Message)
	}

	return json.Marshal(data.Body)
}

// 消息ID
type msg_id struct {
	ID int `json:"message_id"`
}

// 发送私聊消息
func SendPrivateMessage(ws *websocket.Conn, user_id uint, text string) (message_id int, err error) {
	return SendTemporaryMessage(ws, user_id, 0, text)
}

// 发送群消息
func SendGroupMessage(ws *websocket.Conn, group_id uint, text string) (message_id int, err error) {
	return SendTemporaryMessage(ws, 0, group_id, text)
}

// 发送临时会话消息
func SendTemporaryMessage(ws *websocket.Conn, user_id, group_id uint, text string) (message_id int, err error) {
	type params struct {
		UserID  uint   `json:"user_id"`
		GroupID uint   `json:"group_id"`
		Message string `json:"message"`
	}

	var msg *msg_id

	var message = ws_data{
		Action: "send_msg",
		Params: params{user_id, group_id, text},
	}

	data, err := message.do(ws)
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal(data, &msg)
	if err != nil {
		return 0, err
	}

	return msg.ID, nil
}

type forward_msg_id struct {
	ID int `json:"id"`
}

// 转发消息ID
type forward_message_id struct {
	Type string         `json:"type"`
	Data forward_msg_id `json:"data"`
}

// 自定义转发消息
type forward_message_custom struct {
	Type string         `json:"type"`
	Data ForwardMessage `json:"data"`
}

// 自定义转发消息
type ForwardMessage struct {
	Name    string `json:"name"`
	Uin     int    `json:"uin"`
	Content string `json:"content"`
	Seq     string `json:"seq"`
}

// 发送转发消息数据
func sendForwardData(ws *websocket.Conn, user_id, group_id uint, v any) (message_id int, err error) {
	type params struct {
		UserID   uint `json:"user_id"`
		GroupID  uint `json:"group_id"`
		Messages any  `json:"messages"`
	}

	var msg *msg_id

	var message = ws_data{
		Action: "send_forward_msg",
		Params: params{user_id, group_id, v},
	}

	data, err := message.do(ws)
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal(data, &msg)
	if err != nil {
		return 0, err
	}

	return msg.ID, err
}

// 发送转发消息ID (私聊)
func SendPrivateForwardMessageID(ws *websocket.Conn, user_id uint, message_ids ...int) (message_id int, err error) {
	var contents []forward_message_id

	for _, mid := range message_ids {
		contents = append(contents, forward_message_id{
			Type: "node",
			Data: forward_msg_id{mid},
		})
	}

	return sendForwardData(ws, user_id, 0, contents)
}

// 发送自定义转发消息 (私聊)
func SendPrivateForwardMessageCustom(ws *websocket.Conn, user_id uint, messages []ForwardMessage) (message_id int, err error) {
	var contents []forward_message_custom

	for _, message := range messages {
		contents = append(contents, forward_message_custom{
			Type: "node",
			Data: message,
		})
	}

	return sendForwardData(ws, user_id, 0, contents)
}

// 发送转发消息ID (群)
func SendGroupForwardMessageID(ws *websocket.Conn, group_id uint, message_ids ...int) (message_id int, err error) {
	var messages []forward_message_id

	for _, mid := range message_ids {
		messages = append(messages, forward_message_id{
			Type: "node",
			Data: forward_msg_id{mid},
		})
	}

	return sendForwardData(ws, 0, group_id, messages)
}

// 发送自定义转发消息 (群)
func SendGroupForwardMessageCustom(ws *websocket.Conn, group_id uint, messages []ForwardMessage) (message_id int, err error) {
	var contents []forward_message_custom

	for _, message := range messages {
		contents = append(contents, forward_message_custom{
			Type: "node",
			Data: message,
		})
	}

	return sendForwardData(ws, 0, group_id, contents)
}

// 标记消息已读
func MarkMessageRead(ws *websocket.Conn, message_id int) error {
	var message = ws_data{
		Action: "send_forward_msg",
		Params: msg_id{message_id},
	}

	_, err := message.do(ws)
	return err
}

// 撤回消息
func DeleteMessage(ws *websocket.Conn, message_id int) error {
	var message = ws_data{
		Action: "delete_msg",
		Params: msg_id{message_id},
	}

	_, err := message.do(ws)
	return err
}

type sender struct { // 发送人信息
	Age      int    `json:"age"`      // 年龄
	Area     string `json:"area"`     // 地区
	Card     string `json:"card"`     // 群名片／备注
	Level    string `json:"level"`    // 成员等级
	Nickname string `json:"nickname"` // 昵称
	Role     string `json:"role"`     // 角色, owner/admin/member
	Sex      string `json:"sex"`      // 性别, male/female/unknown
	Title    string `json:"title"`    // 专属头衔
	UserID   uint   `json:"user_id"`  // 发送者QQ号
}

type sender_lite struct {
	Nickname string `json:"nickname"` // 昵称
	UserID   uint   `json:"user_id"`  // QQ号
}

// 消息数据
type message_content struct {
	Group       bool        `json:"group"`         // 群聊
	GroupID     uint        `json:"group_id"`      // 群号
	Message     string      `json:"message"`       // 内容
	MessageID   int         `json:"message_id"`    // ID
	MessageIDV2 string      `json:"message_id_v2"` // ID v2
	MessageSeq  int         `json:"message_seq"`   // 序列
	MessageType string      `json:"message_type"`  // 类型
	RealID      int         `json:"real_id"`       // 真实ID
	Sender      sender_lite `json:"sender"`        // 发送者信息
	Time        int         `json:"time"`          // 发送消息时的时间戳
}

// 获取消息
func GetMessage(ws *websocket.Conn, message_id int) (content *message_content, err error) {
	var message = ws_data{
		Action: "get_msg",
		Params: msg_id{message_id},
	}

	data, err := message.do(ws)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &content)
	if err != nil {
		return nil, err
	}

	return content, nil
}

// 图片信息
type image struct {
	File     string `json:"file"`     // 图片缓存路径
	Filename string `json:"filename"` // 图片文件原名
	Size     int    `json:"size"`     // 图片源文件大小
	URL      string `json:"url"`      // 图片下载地址
}

// 获取图片信息
func GetImage(ws *websocket.Conn, file string) (image *image, err error) {
	type params struct {
		File string `json:"file"`
	}

	var message = ws_data{
		Action: "get_image",
		Params: params{file},
	}

	data, err := message.do(ws)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &image)
	if err != nil {
		return nil, err
	}

	return image, nil
}

// 群组踢人
func SetGroupKick(ws *websocket.Conn, group_id, user_id uint, reject_add_request bool) error {
	type params struct {
		GroupID          uint `json:"group_id"`
		UserID           uint `json:"user_id"`
		RejectAddRequest bool `json:"reject_add_request"`
	}

	var message = ws_data{
		Action: "set_group_kick",
		Params: params{group_id, user_id, reject_add_request},
	}

	_, err := message.do(ws)
	return err
}

// 群组单人禁言
func SetGroupBan(ws *websocket.Conn, group_id, user_id uint, duration uint) error {
	type params struct {
		GroupID  uint `json:"group_id"`
		UserID   uint `json:"user_id"`
		Duration uint `json:"duration"`
	}

	var message = ws_data{
		Action: "set_group_ban",
		Params: params{group_id, user_id, duration},
	}

	_, err := message.do(ws)
	return err
}

// 群组匿名用户禁言
func SetGroupAnonymousBan(ws *websocket.Conn, group_id uint, flag string, duration uint) error {
	type params struct {
		GroupID  uint   `json:"group_id"`
		Flag     string `json:"flag"`
		Duration uint   `json:"duration"`
	}

	var message = ws_data{
		Action: "set_group_anonymous_ban",
		Params: params{group_id, flag, duration},
	}

	_, err := message.do(ws)
	return err
}

// 群组全员禁言
func SetGroupWholeBan(ws *websocket.Conn, group_id uint, enable bool) error {
	type params struct {
		GroupID uint `json:"group_id"`
		Enable  bool `json:"enable"`
	}

	var message = ws_data{
		Action: "set_group_whole_ban",
		Params: params{group_id, enable},
	}

	_, err := message.do(ws)
	return err
}

// 群组设置管理员
func SetGroupAdmin(ws *websocket.Conn, group_id, user_id uint, enable bool) error {
	type params struct {
		GroupID uint `json:"group_id"`
		UserID  uint `json:"user_id"`
		Enable  bool `json:"enable"`
	}

	var message = ws_data{
		Action: "set_group_admin",
		Params: params{group_id, user_id, enable},
	}

	_, err := message.do(ws)
	return err
}

// 该API暂未被go-cqhttp支持
// TODO: 群组匿名

// 设置群名片(群备注)
func SetGroupCard(ws *websocket.Conn, group_id, user_id uint, card string) error {
	type params struct {
		GroupID uint   `json:"group_id"`
		UserID  uint   `json:"user_id"`
		Card    string `json:"card"`
	}

	var message = ws_data{
		Action: "set_group_card",
		Params: params{group_id, user_id, card},
	}

	_, err := message.do(ws)
	return err
}

// 设置群名
func SetGroupName(ws *websocket.Conn, group_id uint, name string) error {
	type params struct {
		GroupID uint   `json:"group_id"`
		Name    string `json:"group_name"`
	}

	var message = ws_data{
		Action: "set_group_name",
		Params: params{group_id, name},
	}

	_, err := message.do(ws)
	return err
}

// 退出群组
func SetGroupLeave(ws *websocket.Conn, group_id uint, dismiss bool) error {
	type params struct {
		GroupID uint `json:"group_id"`
		Dismiss bool `json:"is_dismiss"`
	}

	var message = ws_data{
		Action: "set_group_leave",
		Params: params{group_id, dismiss},
	}

	_, err := message.do(ws)
	return err
}

// 设置群组专属头衔
func SetGroupSpecialTitle(ws *websocket.Conn, group_id, user_id uint, title string, duration uint) error {
	type params struct {
		GroupID      uint   `json:"group_id"`
		UserID       uint   `json:"user_id"`
		SpecialTitle string `json:"special_title"`
		Duration     uint   `json:"duration"`
	}

	var message = ws_data{
		Action: "set_group_special_title",
		Params: params{group_id, user_id, title, duration},
	}

	_, err := message.do(ws)
	return err
}

// 群打卡
func SendGroupSign(ws *websocket.Conn, group_id uint) error {
	type params struct {
		GroupID uint `json:"group_id"`
	}

	var message = ws_data{
		Action: "send_group_sign",
		Params: params{group_id},
	}

	_, err := message.do(ws)
	return err
}

// 处理加好友请求
func SetFriendAddRequest(ws *websocket.Conn, flag string, approve bool, remark string) error {
	type params struct {
		Flag    string `json:"flag"`
		Approve bool   `json:"approve"`
		Remark  string `json:"remark"`
	}

	var message = ws_data{
		Action: "set_friend_add_request",
		Params: params{flag, approve, remark},
	}

	_, err := message.do(ws)
	return err
}

// 处理加群请求／邀请
func SetGroupAddRequest(ws *websocket.Conn, flag, sub_type string, approve bool, remark string) error {
	type params struct {
		Flag    string `json:"flag"`
		SubType string `json:"sub_type"`
		Approve bool   `json:"approve"`
		Remark  string `json:"remark"`
	}

	var message = ws_data{
		Action: "set_group_add_request",
		Params: params{flag, sub_type, approve, remark},
	}

	_, err := message.do(ws)
	return err
}

// 登录信息
type login struct {
	Nickname string `json:"nickname"` // 昵称
	UserID   uint   `json:"user_id"`  // QQ号
}

// 获取登录号信息
func GetLoginInfo(ws *websocket.Conn) (login *login, err error) {
	var message = ws_data{Action: "get_login_info"}

	data, err := message.do(ws)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &login)
	if err != nil {
		return nil, err
	}

	return login, err
}

// 企点账号信息
type qidian struct {
	MasterID   uint   `json:"master_id"`   // 父账号ID
	ExtName    string `json:"ext_name"`    // 用户昵称
	CreateTime int    `json:"create_time"` // 账号创建时间
}

// 获取企点账号信息
//
// * 该API只有企点协议可用
func GetQidianAccountInfo(ws *websocket.Conn) (qidian *qidian, err error) {
	var message = ws_data{Action: "qidian_get_account_info"}

	data, err := message.do(ws)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &qidian)
	if err != nil {
		return nil, err
	}

	return qidian, err
}

// 设置个人资料
//
// * 该API缺少文档描述, 根据源码编写
func SetProfile(ws *websocket.Conn, nickname, company, email, college, personal_note string) error {
	type params struct {
		Nickname     string `json:"nickname"`
		Company      string `json:"company"`
		Email        string `json:"email"`
		College      string `json:"college"`
		PersonalNote string `json:"personal_note"`
	}

	var message = ws_data{
		Action: "set_qq_profile",
		Params: params{nickname, company, email, college, personal_note},
	}

	_, err := message.do(ws)
	return err
}

// 陌生人信息
type stranger struct {
	Age       int    `json:"age"`        // 年龄
	Level     int    `json:"level"`      // 等级
	LoginDays int    `json:"login_days"` // QQ达人
	Nickname  string `json:"nickname"`   // 昵称
	Qid       string `json:"qid"`        // QQ ID身份卡
	Sex       string `json:"sex"`        // 性别, male/female/unknown
	UserID    uint   `json:"user_id"`    // QQ号
}

// 获取陌生人信息
func GetStrangerInfo(ws *websocket.Conn, user_id uint) (stranger *stranger, err error) {
	type params struct {
		UserID uint `json:"user_id"`
	}

	var message = ws_data{
		Action: "get_stranger_info",
		Params: params{user_id},
	}

	data, err := message.do(ws)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &stranger)
	if err != nil {
		return nil, err
	}

	return stranger, err
}

// 好友信息
type friend struct {
	Nickname string `json:"nickname"` // 昵称
	Remark   string `json:"remark"`   // 备注名
	UserID   uint   `json:"user_id"`  // QQ号
}

// 获取好友列表
func GetFriendList(ws *websocket.Conn) (list *[]friend, err error) {
	var message = ws_data{
		Action: "get_friend_list",
	}

	data, err := message.do(ws)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &list)
	if err != nil {
		return nil, err
	}

	return list, err
}

// 单向好友信息
type unidirectional_friend struct {
	Nickname string `json:"nickname"` // 昵称
	UserID   uint   `json:"user_id"`  // QQ号
	Source   string `json:"source"`   // 添加途径
}

// 获取单向好友列表
func GetUnidirectionalFriendList(ws *websocket.Conn) (list *[]unidirectional_friend, err error) {
	var message = ws_data{Action: "get_unidirectional_friend_list"}

	data, err := message.do(ws)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// 删除单向好友
func DeleteUnidirectionalFriend(ws *websocket.Conn, user_id uint) error {
	type params struct {
		UserID uint `json:"user_id"`
	}

	var message = ws_data{
		Action: "delete_unidirectional_friend",
		Params: params{user_id},
	}

	_, err := message.do(ws)
	return err
}

// 删除好友
func DeleteFriend(ws *websocket.Conn, friend_id uint) error {
	type params struct {
		FriendID uint `json:"friend_id"`
	}

	var message = ws_data{
		Action: "delete_friend",
		Params: params{friend_id},
	}

	_, err := message.do(ws)
	return err
}

// 群信息
type group struct {
	GroupCreateTime int    `json:"group_create_time"` // 群创建时间
	GroupID         uint   `json:"group_id"`          // 群号
	GroupLevel      int    `json:"group_level"`       // 群等级
	GroupMemo       string `json:"group_memo"`        // 群备注
	GroupName       string `json:"group_name"`        // 群名称
	MaxMemberCount  uint   `json:"max_member_count"`  // 最大成员数
	MemberCount     uint   `json:"member_count"`      // 成员数
}

// 获取群信息
//
// * 如果机器人尚未加入群, group_create_time, group_level, max_member_count 和 member_count 将会为0
func GetGroupInfo(ws *websocket.Conn, group_id uint) (info *group, err error) {
	type params struct {
		GroupID uint `json:"group_id"`
	}

	var message = ws_data{
		Action: "get_group_info",
		Params: params{group_id},
	}

	data, err := message.do(ws)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &info)
	if err != nil {
		return nil, err
	}

	return info, err
}

// 获取群列表
func GetGroupList(ws *websocket.Conn) (list *[]group, err error) {
	var message = ws_data{Action: "get_group_list"}

	data, err := message.do(ws)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &list)
	if err != nil {
		return nil, err
	}

	return list, err
}

// 获取群成员信息
type group_member struct {
	Age             uint   `json:"age"`               // 年龄
	Area            string `json:"area"`              // 地区
	Card            string `json:"card"`              // 群名片／备注
	CardChangeable  bool   `json:"card_changeable"`   // 是否允许修改群名片
	GroupID         uint   `json:"group_id"`          // 群号
	JoinTime        int    `json:"join_time"`         // 加群时间戳
	LastSentTime    int    `json:"last_sent_time"`    // 最后发言时间戳
	Level           string `json:"level"`             // 成员等级
	Nickname        string `json:"nickname"`          // 昵称
	Role            string `json:"role"`              // 角色, owner/admin/member
	Sex             string `json:"sex"`               // 性别
	ShutUpTimestamp uint   `json:"shut_up_timestamp"` // 禁言到期时间
	Title           string `json:"title"`             // 专属头衔
	TitleExpireTime uint   `json:"title_expire_time"` // 专属头衔过期时间戳
	Unfriendly      bool   `json:"unfriendly"`        // 是否不良记录成员
	UserID          uint   `json:"user_id"`           // QQ号
}

// 获取群成员信息
func GetGroupMemberInfo(ws *websocket.Conn, group_id, user_id uint) (info *group_member, err error) {
	type params struct {
		GroupID uint `json:"group_id"`
		UserID  uint `json:"user_id"`
	}

	var message = ws_data{
		Action: "get_group_member_info",
		Params: params{group_id, user_id},
	}

	data, err := message.do(ws)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &info)
	if err != nil {
		return nil, err
	}

	return info, err
}

// 获取群成员列表
func GetGroupMemberList(ws *websocket.Conn, group_id uint) (list *[]group_member, err error) {
	type params struct {
		GroupID uint `json:"group_id"`
	}

	var message = ws_data{
		Action: "get_group_member_list",
		Params: params{group_id},
	}

	data, err := message.do(ws)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &list)
	if err != nil {
		return nil, err
	}

	return list, err
}

// 群荣誉列表
type honor_list struct {
	GroupID          uint          `json:"group_id"`           // 群号
	CurrentTalkative group_honor   `json:"current_talkative"`  // 当前龙王
	EmotionList      []group_honor `json:"emotion_list"`       // 快乐之源
	LegendList       []group_honor `json:"legend_list"`        // 群聊炽焰
	PerformerLis     []group_honor `json:"performer_lis"`      // 群聊之火
	StrongNewbieList []group_honor `json:"strong_newbie_list"` // 冒尖小春笋
	TalkativeList    []group_honor `json:"talkative_list"`     // 历史龙王
}

// 群荣誉
type group_honor struct {
	Avatar      string `json:"avatar"`      // 头像URL
	Description string `json:"description"` // 荣誉描述
	Nickname    string `json:"nickname"`    // 昵称
	UserID      uint   `json:"user_id"`     // QQ号
}

// 获取群荣誉信息
//
// * type: 要获取的群荣誉类型, talkative, performer, legend, strong_newbie emotion, 以分别获取单个类型的群荣誉数据, 或传入all获取所有数据
func GetGroupHonorInfo(ws *websocket.Conn, group_id uint, types string) (honor *honor_list, err error) {
	type params struct {
		GroupID uint   `json:"group_id"`
		Type    string `json:"type"`
	}

	var message = ws_data{
		Action: "get_group_honor_info",
		Params: params{group_id, types},
	}

	data, err := message.do(ws)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &honor)
	if err != nil {
		return nil, err
	}

	return honor, err
}

// 该API暂未被go-cqhttp支持
// TODO: 获取Cookies
// TODO: 获取CSRF Token
// TODO: 获取QQ相关接口凭证
// TODO: 获取语音

// 是否能发送图片或语音
type can_send_image_or_record struct {
	Bool bool `json:"yes"`
}

// 检查是否可以发送图片
func CanSendImage(ws *websocket.Conn) (status bool, err error) {
	var send *can_send_image_or_record

	message := ws_data{
		Action: "can_send_image",
	}

	data, err := message.do(ws)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal(data, &send)
	if err != nil {
		return false, err
	}

	return send.Bool, err
}

// 检查是否可以发送语音
func CanSendRecord(ws *websocket.Conn) (status bool, err error) {
	var send *can_send_image_or_record

	message := ws_data{
		Action: "can_send_record",
	}

	data, err := message.do(ws)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal(data, &send)
	if err != nil {
		return false, err
	}

	return send.Bool, err
}

// 版本
type version struct {
	AppFullName     string `json:"app_full_name"`    // 应用完整名称
	AppName         string `json:"app_name"`         // 应用标识, 固定值
	AppVersion      string `json:"app_version"`      // 应用版本
	CoolqDirectory  string `json:"coolq_directory"`  // 原CoolQ运行目录
	Protocol        int    `json:"protocol"`         // 登陆使用协议类型
	ProtocolVersion string `json:"protocol_version"` // OneBot标准版本, 固定值
	RuntimeOS       string `json:"runtime_os"`       // 运行时操作系统
	RuntimeVersion  string `json:"runtime_version"`  // 运行时版本
	Version         string `json:"version"`          // 应用版本
}

// 获取版本信息
func GetVersionInfo(ws *websocket.Conn) (ver *version, err error) {
	var message = ws_data{Action: "get_version_info"}

	data, err := message.do(ws)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &ver)
	if err != nil {
		return nil, err
	}

	return ver, err
}

// 该API暂未被go-cqhttp支持
// TODO: 清理缓存

// 获取群头像URL
func GetGroupAvatarURL(group_id uint) (url string) {
	return fmt.Sprint("https://p.qlogo.cn/gh/", group_id, "/", group_id, "/100")
}

// 设置群头像
//
// * 目前这个API在登录一段时间后因cookie失效而失效, 请考虑后使用
func SetGroupPortrait(ws *websocket.Conn, group_id uint, file string) error {
	type params struct {
		GroupID uint   `json:"group_id"`
		File    string `json:"file"`
	}

	var message = ws_data{
		Action: "set_group_portrait",
		Params: params{group_id, file},
	}

	_, err := message.do(ws)
	return err
}

// 分词
type word struct {
	Slices []string `json:"slices"`
}

// 获取中文分词
func GetWordSlices(ws *websocket.Conn, text string) (slice []string, err error) {
	type params struct {
		Content string `json:"content"`
	}

	var word *word

	message := ws_data{
		Action: ".get_word_slices",
		Params: params{text},
	}

	data, err := message.do(ws)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &word)
	if err != nil {
		return nil, err
	}

	return word.Slices, err
}

// OCR
type ocr struct {
	Texts    []text `json:"texts"`    // 文本
	Language string `json:"language"` // 语言
}

type text struct {
	Text        string        `json:"text"` // 置信度
	Confidence  int           `json:"confidence"`
	Coordinates []coordinates `json:"coordinates"` // 坐标
}

type coordinates struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// 图片OCR
func OcrImage(ws *websocket.Conn, image string) (ocr *ocr, err error) {
	type params struct {
		Image string `json:"image"`
	}

	var message = ws_data{
		Action: "ocr_image",
		Params: params{image},
	}

	data, err := message.do(ws)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &ocr)
	if err != nil {
		return nil, err
	}

	return ocr, err
}

// 群系統消息
type group_system_msg struct {
	InvitedRequest *[]invited_request `json:"invited_requests,omitempty"` // 邀请消息列表
	JoinRequest    *[]join_request    `json:"join_requests,omitempty"`    // 进群消息列表
}

// 邀请消息列表
type invited_request struct {
	Actor         uint   `json:"actor"`          // 处理者, 未处理为0
	Checked       bool   `json:"checked"`        // 是否已被处理
	GroupID       uint   `json:"group_id"`       // 群号
	GroupName     string `json:"group_name"`     // 群名
	RequestID     int    `json:"request_id"`     // 请求ID
	RequesterNick string `json:"requester_nick"` // 请求者昵称
	RequesterUin  uint   `json:"requester_uin"`  // 请求者ID
}

// 进群消息列表
type join_request struct {
	Actor         uint   `json:"actor"`          // 处理者, 未处理为0
	Checked       bool   `json:"checked"`        // 是否已被处理
	GroupID       uint   `json:"group_id"`       // 群号
	GroupName     string `json:"group_name"`     // 群名
	Message       string `json:"message"`        // 验证消息
	RequestID     int    `json:"request_id"`     // 请求ID
	RequesterNick string `json:"requester_nick"` // 请求者昵称
	RequesterUin  uint   `json:"requester_uin"`  // 请求者ID
}

// 获取群系统消息
func GetGroupSystemMessage(ws *websocket.Conn) (system_message *group_system_msg, err error) {
	var message = ws_data{Action: "get_group_system_msg"}

	data, err := message.do(ws)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &system_message)
	if err != nil {
		return nil, err
	}

	return system_message, err
}

// 上传私聊文件
//
// * 只能上传本地文件, 需要上传 http 文件的话请先调用DownloadFile()下载
func UploadPrivateFile(ws *websocket.Conn, user_id uint, file, name string) error {
	type params struct {
		UserID uint   `json:"user_id"`
		File   string `json:"file"`
		Name   string `json:"name"`
	}

	var message = ws_data{
		Action: "upload_private_file",
		Params: params{user_id, file, name},
	}

	_, err := message.do(ws)
	return err
}

// 上传群文件
//
// * 在不提供 folder 参数的情况下默认上传到根目录
//
// * 只能上传本地文件, 需要上传 http 文件的话请先调用DownloadFile()下载
func UploadGroupFile(ws *websocket.Conn, group_id uint, file, name, folder string) error {
	type params struct {
		GroupID uint   `json:"group_id"`
		File    string `json:"file"`
		Name    string `json:"name"`
		Folder  string `json:"folder"`
	}

	var message = ws_data{
		Action: "upload_group_file",
		Params: params{group_id, file, name, folder},
	}

	_, err := message.do(ws)
	return err
}

// 群文件系统信息
type group_file_system struct {
	FileCount  uint `json:"file_count"`  // 文件总数
	LimitCount uint `json:"limit_count"` // 文件上限
	UsedSpace  uint `json:"used_space"`  // 已使用空间, 单位byte
	TotalSpace uint `json:"total_space"` // 空间上限, 单位byte
	GroupID    uint `json:"group_id"`    // 群号
}

// 获取群文件系统信息
func GetGroupFileSystemInfo(ws *websocket.Conn, group_id uint) (file *group_file_system, err error) {
	type params struct {
		GroupID uint `json:"group_id"`
	}

	var message = ws_data{
		Action: "get_group_file_system_info",
		Params: params{group_id},
	}

	data, err := message.do(ws)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &file)
	if err != nil {
		return nil, err
	}

	return file, err
}

// 群根目录文件列表
type group_files struct {
	Files   *[]file   `json:"files"`   // 文件列表
	Folders *[]folder `json:"folders"` // 文件夹列表
}

// 文件
type file struct {
	GroupID       uint   `json:"group_id"`       // 群号
	FileID        string `json:"file_id"`        // 文件ID
	FileName      string `json:"file_name"`      // 文件名
	Busid         int    `json:"busid"`          // 文件类型
	FileSize      uint   `json:"file_size"`      // 文件大小, 单位byte
	UploadTime    uint   `json:"upload_time"`    // 上传时间
	DeadTime      uint   `json:"dead_time"`      // 过期时间, 永久文件为0
	ModifyTime    int    `json:"modify_time"`    // 最后修改时间
	DownloadTimes int    `json:"download_times"` // 下载次数
	Uploader      uint   `json:"uploader"`       // 上传者ID
	UploaderName  string `json:"uploader_name"`  // 上传者名字
}

// 文件夹
type folder struct {
	GroupID        uint   `json:"group_id"`         // 群号
	FolderID       string `json:"folder_id"`        // 文件夹ID
	FolderName     string `json:"folder_name"`      // 文件名
	CreateTime     int    `json:"create_time"`      // 创建时间
	Creator        uint   `json:"creator"`          // 创建者
	CreatorName    string `json:"creator_name"`     // 创建者名字
	TotalFileCount uint   `json:"total_file_count"` // 子文件数量
}

// 获取群根目录文件列表
func GetGroupRootFiles(ws *websocket.Conn, group_id uint) (file *group_files, err error) {
	type params struct {
		GroupID uint `json:"group_id"`
	}

	var message = ws_data{
		Action: "get_group_root_files",
		Params: params{group_id},
	}

	data, err := message.do(ws)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &file)
	if err != nil {
		return nil, err
	}

	return file, err
}

// 获取群子目录文件列表
func GetGroupFilesByFolder(ws *websocket.Conn, group_id uint, folder_id string) (file *group_files, err error) {
	type params struct {
		GroupID  uint   `json:"group_id"`
		FolderID string `json:"folder_id"`
	}

	var message = ws_data{
		Action: "get_group_files_by_folder",
		Params: params{group_id, folder_id},
	}

	data, err := message.do(ws)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &file)
	if err != nil {
		return nil, err
	}

	return file, err
}

// 创建群文件文件夹
func CreateGroupFileFolder(ws *websocket.Conn, group_id uint, name string) error {
	type params struct {
		GroupID  uint   `json:"group_id"`
		Name     string `json:"name"`
		ParentID string `json:"parent_id"`
	}

	var message = ws_data{
		Action: "create_group_file_folder",
		Params: params{group_id, name, "/"},
	}

	_, err := message.do(ws)
	return err
}

// 删除群文件文件夹
func DeleteGroupFolder(ws *websocket.Conn, group_id uint, folder_id string) error {
	type params struct {
		GroupID  uint   `json:"group_id"`
		FolderID string `json:"folder_id"`
	}

	var message = ws_data{
		Action: "delete_group_folder",
		Params: params{group_id, folder_id},
	}

	_, err := message.do(ws)
	return err
}

// 删除群文件
func DeleteGroupFile(ws *websocket.Conn, group_id uint, file_id string, busid int) error {
	type params struct {
		GroupID uint   `json:"group_id"`
		FileID  string `json:"file_id"`
		BusID   int    `json:"busid"`
	}

	var message = ws_data{
		Action: "delete_group_file",
		Params: params{group_id, file_id, busid},
	}

	_, err := message.do(ws)
	return err
}

// 群文件URL
type group_file struct {
	URL string `json:"url"`
}

// 获取群文件资源链接
func GetGroupFileURL(ws *websocket.Conn, group_id uint, file_id string, busid int) (URL string, err error) {
	type params struct {
		GroupID uint   `json:"group_id"`
		FileID  string `json:"file_id"`
		BusID   int    `json:"busid"`
	}

	var file *group_file

	message := ws_data{
		Action: "get_group_file_url",
		Params: params{group_id, file_id, busid},
	}

	data, err := message.do(ws)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(data, &file)
	if err != nil {
		return "", err
	}

	return file.URL, err
}

// 状态
type status struct {
	Online bool     `json:"online"` // 表示BOT是否在线
	More   struct { // 运行统计
		DisconnectTime  int `json:"DisconnectTimes"` // TCP链接断开次数
		LastMessageTime int `json:"LastMessageTime"` // 最后一次发送消息的时间戳
		LostTime        int `json:"LostTimes"`       // 账号掉线次数
		MessageReceived int `json:"MessageReceived"` // 接受信息总数
		MessageSent     int `json:"MessageSent"`     // 发送信息总数
		PacketLost      int `json:"PacketLost"`      // 数据包丢失总数
		PacketReceived  int `json:"PacketReceived"`  // 收到的数据包总数
		PacketSent      int `json:"PacketSent"`      // 发送的数据包总数
	} `json:"stat"`
}

// 获取状态
func GetStatus(ws *websocket.Conn) (status *status, err error) {
	var message = ws_data{
		Action: "get_status",
	}

	data, err := message.do(ws)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &status)
	if err != nil {
		return nil, err
	}

	return status, err
}

// @全体成员
type at_all struct {
	CanAtAll                 bool `json:"can_at_all"`                    // 是否可以@全体成员
	RemainAtAllCountForGroup int  `json:"remain_at_all_count_for_group"` // 群内所有管理当天剩余@全体成员次数
	RemainAtAllCountForUin   int  `json:"remain_at_all_count_for_uin"`   // Bot 当天剩余@全体成员次数
}

// 获取群@全体成员剩余次数
func GetGroupAtAllRemain(ws *websocket.Conn, group_id uint) (at *at_all, err error) {
	type params struct {
		GroupID uint `json:"group_id"`
	}

	var message = ws_data{
		Action: "get_group_at_all_remain",
		Params: params{group_id},
	}

	data, err := message.do(ws)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &at)
	if err != nil {
		return nil, err
	}

	return at, err
}

// 群公告
type group_notice struct {
	Message struct {
		Images *[]group_notice_image `json:"images"` // 图片
		Text   string                `json:"text"`   // 公告内容
	} `json:"message"`
	PublishTime int  `json:"publish_time"` // 发送时间
	SenderID    uint `json:"sender_id"`    // 操作者
}

// 群公告图片
type group_notice_image struct {
	ID     string `json:"id"`     // 图片ID
	Height string `json:"height"` // 高
	Width  string `json:"width"`  // 宽
}

// 发送群公告
func SendGroupNotice(ws *websocket.Conn, group_id uint, content, image string) error {
	type params struct {
		GroupID uint   `json:"group_id"`
		Content string `json:"content"`
		Image   string `json:"image"`
	}

	var message = ws_data{
		Action: "_send_group_notice",
		Params: params{group_id, content, image},
	}

	_, err := message.do(ws)
	return err
}

// 获取群公告
func GetGroupNotice(ws *websocket.Conn, group_id uint) (notice *[]group_notice, err error) {
	type params struct {
		GroupID uint `json:"group_id"`
	}

	var message = ws_data{
		Action: "_get_group_notice",
		Params: params{group_id},
	}

	data, err := message.do(ws)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &notice)
	if err != nil {
		return nil, err
	}

	return notice, err
}

// 重载事件过滤器
func ReloadEventFilter(ws *websocket.Conn, file string) error {
	type params struct {
		File string `json:"file"`
	}

	var message = ws_data{
		Action: "reload_event_filter",
		Params: params{file},
	}

	_, err := message.do(ws)
	return err
}

// 资源URL
type filepath struct {
	Path string `json:"file"`
}

// 下载文件到缓存目录
//
// * headers格式: User-Agent=YOUR_UA[\r\n]Referer=https://www.example.com
//
// * [\r\n] 为换行符, 使用http请求时请注意编码
//
// * 调用后会阻塞直到下载完成后才会返回数据，请注意下载大文件时的超时
func DownloadFile(ws *websocket.Conn, url string, header []string, thread_count int) (file string, err error) {
	type params struct {
		URL         string   `json:"url"`
		ThreadCount int      `json:"thread_count"`
		Header      []string `json:"headers"`
	}

	var download *filepath

	message := ws_data{
		Action: "download_file",
		Params: params{url, thread_count, header},
	}

	data, err := message.do(ws)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(data, &download)
	if err != nil {
		return "", err
	}

	return download.Path, err
}

// 在线设备
type online struct {
	Clients []client `json:"clients"` // 在线客户端列表
}

type client struct {
	AppID      int    `json:"app_id"`      // 客户端ID
	DeviceKind string `json:"device_kind"` // 设备类型
	DeviceName string `json:"device_name"` // 设备名称
}

// 获取当前账号在线客户端列表
func GetOnlineClients(ws *websocket.Conn) (client *online, err error) {
	var message = ws_data{Action: "get_online_clients"}

	data, err := message.do(ws)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &client)
	if err != nil {
		return nil, err
	}

	return client, err
}

// 群消息历史记录
type history_message struct {
	Messages []group_message `json:"messages"` // 从起始序号开始的前19条消息
}

// 群消息
type group_message struct {
	// 匿名信息, 如果不是匿名消息则为null
	Anonymous   anonymous `json:"anonymous"`
	Font        int       `json:"font"`         // 字体
	Time        int       `json:"time"`         // 事件发生的时间戳
	SelfID      uint      `json:"self_id"`      // 收到事件的机器人QQ号
	PostType    string    `json:"post_type"`    // 上报类型
	MessageType string    `json:"message_type"` // 消息类型
	SubType     string    `json:"sub_type"`     // 消息子类型, 正常消息是normal, 匿名消息是anonymous, 系统提示(如「管理员已禁止群内匿名聊天」)是notice
	MessageID   int       `json:"message_id"`   // 消息ID
	GroupID     uint      `json:"group_id"`     // 群号
	UserID      uint      `json:"user_id"`      // 发送者QQ号
	Message     string    `json:"message"`      // 消息内容
	RawMessage  string    `json:"raw_message"`  // 原始消息内容
	MessageSeq  int       `json:"message_seq"`  // 消息序列
	Sender      sender    `json:"sender"`
}

// 匿名信息
type anonymous struct {
	ID   int    `json:"id"`   // 匿名用户ID
	Name string `json:"name"` // 匿名用户名称
	Flag string `json:"flag"` // 匿名用户flag, 在调用禁言API时需要传入
}

// 获取群消息历史记录
func GetGroupMessageHistory(ws *websocket.Conn, message_seq, group_id uint) (history *history_message, err error) {
	type params struct {
		Seq     uint `json:"message_seq"`
		GroupID uint `json:"group_id"`
	}

	var message = ws_data{
		Action: "get_group_msg_history",
		Params: params{message_seq, group_id},
	}

	data, err := message.do(ws)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &history)
	if err != nil {
		return nil, err
	}

	return history, err
}

// 精华消息
type essence struct {
	MessageID    int    `json:"message_id"`    // 消息ID
	OperatorID   uint   `json:"operator_id"`   // 操作者QQ号
	OperatorNick string `json:"operator_nick"` // 操作者昵称
	OperatorTime int    `json:"operator_time"` // 精华设置时间
	SenderID     uint   `json:"sender_id"`     // 发送者QQ号
	SenderNick   string `json:"sender_nick"`   // 发送者昵称
	SenderTime   int    `json:"sender_time"`   // 消息发送时间
}

// 设置精华消息
func SetEssenceMessage(ws *websocket.Conn, message_id int) error {
	type params struct {
		MessageID int `json:"message_id"`
	}

	var message = ws_data{
		Action: "set_essence_msg",
		Params: params{message_id},
	}

	_, err := message.do(ws)
	return err
}

// 移出精华消息
func DeleteEssenceMessage(ws *websocket.Conn, message_id int) error {
	type params struct {
		MessageID int `json:"message_id"`
	}

	var message = ws_data{
		Action: "delete_essence_msg",
		Params: params{message_id},
	}

	_, err := message.do(ws)
	return err
}

// 获取精华消息列表
func GetEssenceMessageList(ws *websocket.Conn, group_id uint) (essence *[]essence, err error) {
	type params struct {
		GroupID uint `json:"group_id"`
	}

	var message = ws_data{
		Action: "get_essence_msg_list",
		Params: params{group_id},
	}

	data, err := message.do(ws)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &essence)
	if err != nil {
		return nil, err
	}

	return essence, err
}

type url_safely struct {
	Level int `json:"level"`
}

// 检查链接安全性
//
// * level: 安全等级, 1.安全 2.未知 3.危险
func CheckURLSafely(ws *websocket.Conn, url string) (level int, err error) {
	type params struct {
		URL string `json:"url"`
	}

	var url_safely *url_safely

	var message = ws_data{
		Action: "check_url_safely",
		Params: params{url},
	}

	data, err := message.do(ws)
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal(data, &url_safely)
	if err != nil {
		return 0, err
	}

	return url_safely.Level, err
}

type model struct {
	Variants []variant `json:"variants"`
}

type variant struct {
	ModelShow string `json:"model_show"` // 显示机型
	NeedPay   bool   `json:"need_pay"`   // 是否需要会员
}

// 获取在线机型
func GetModelShow(ws *websocket.Conn, content string) (model *model, err error) {
	type params struct {
		Model string `json:"model"`
	}

	var message = ws_data{
		Action: "_get_model_show",
		Params: params{content},
	}

	data, err := message.do(ws)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &model)
	if err != nil {
		return nil, err
	}

	return model, err
}

// 设置在线机型
func SetModelShow(ws *websocket.Conn, content, model_show string) error {
	type params struct {
		Model     string `json:"model"`
		ModelShow string `json:"model_show"`
	}

	var message = ws_data{
		Action: "_set_model_show",
		Params: params{content, model_show},
	}

	_, err := message.do(ws)
	return err
}
