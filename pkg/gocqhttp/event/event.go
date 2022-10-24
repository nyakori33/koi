package event

// 生命周期
type Lifecycle struct {
	// 通信方式
	/*
		1: HTTP通信
		2: 正向 Websocket 通信
		3: 反向 Websocket 通信
		4: pprof 性能分析服务器
		5: 云函数服务
	*/
	PostMethod    int    `json:"_post_method"`
	MetaEventType string `json:"meta_event_type"` // 元事件类型
	SelfID        uint   `json:"self_id"`         // 收到事件的机器人QQ号
	PostType      string `json:"post_type"`       // 上报类型
	SubType       string `json:"sub_type"`        // 事件子类型，分别表示 go-cqhttp 启用、停用、WebSocket 连接成功
	Time          int    `json:"time"`            // 事件发生的时间戳
}

// 心跳
type Heartbeat struct {
	Interval      int    `json:"interval"`        // 到下次心跳的间隔，单位毫秒
	MetaEventType string `json:"meta_event_type"` // 元事件类型
	PostType      string `json:"post_type"`       // 上报类型
	SelfID        uint   `json:"self_id"`         // 收到事件的机器人QQ号
	Status        status `json:"status"`          // 状态信息
	Time          int    `json:"time"`            // 事件发生的时间戳
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

// 私聊消息
type PrivateMessage struct {
	Font        int      `json:"font"`         // 字体
	Message     string   `json:"message"`      // 消息内容
	MessageID   int      `json:"message_id"`   // 消息ID
	MessageType string   `json:"message_type"` // 消息类型
	PostType    string   `json:"post_type"`    // 上报类型
	RawMessage  string   `json:"raw_message"`  // 原始消息内容
	SelfID      uint     `json:"self_id"`      // 收到事件的机器人QQ号
	Sender      struct { // 发送人信息
		Age      uint   `json:"age"`      // 年龄
		Nickname string `json:"nickname"` // 昵称
		Sex      string `json:"sex"`      // 性别, male/female/unknown
		UserID   uint   `json:"user_id"`  // 发送者QQ号
	} `json:"sender"`
	SubType string `json:"sub_type"` // 消息子类型, 如果是好友则是friend, 如果是群临时会话则是group, 如果是在群中自身发送则是group_self
	// 临时会话来源
	// 0: 群聊
	// 1: QQ咨询
	// 2: 查找
	// 3: QQ电影
	// 4: 热聊
	// 6: 验证消息
	// 7: 多人聊天
	// 8: 约会
	// 9: 通讯录
	TempSource int  `json:"temp_source"`
	TargetID   uint `json:"target_id"` // 接收者QQ号
	Time       int  `json:"time"`      // 事件发生的时间戳
	UserID     uint `json:"user_id"`   // 发送者QQ号
}

// 群消息
type GroupMessage struct {
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
	Sender      struct {  // 发送人信息
		Age      uint   `json:"age"`      // 年龄
		Area     string `json:"area"`     // 地区
		Card     string `json:"card"`     // 群名片／备注
		Level    string `json:"level"`    // 成员等级
		Nickname string `json:"nickname"` // 昵称
		Role     string `json:"role"`     // 角色, owner/admin/member
		Sex      string `json:"sex"`      // 性别, male/female/unknown
		Title    string `json:"title"`    // 专属头衔
		UserID   uint   `json:"user_id"`  // 发送者QQ号
	} `json:"sender"`
}

// 匿名信息
type anonymous struct {
	ID   int    `json:"id"`   // 匿名用户ID
	Name string `json:"name"` // 匿名用户名称
	Flag string `json:"flag"` // 匿名用户flag, 在调用禁言API时需要传入
}

// 加好友请求
type FriendRequest struct {
	Time        int    `json:"time"`         // 事件发生的时间戳
	SelfID      uint   `json:"self_id"`      // 收到事件的机器人QQ号
	PostType    string `json:"post_type"`    // 上报类型
	RequestType string `json:"request_type"` // 请求类型
	UserID      uint   `json:"user_id"`      // 发送请求的QQ号
	Comment     string `json:"comment"`      // 验证信息
	Flag        string `json:"flag"`         // 请求flag, 在调用处理请求的API时需要传入
}

// 加群请求/邀请
type GroupRequest struct {
	Time        int    `json:"time"`         // 事件发生的时间戳
	SelfID      uint   `json:"self_id"`      // 收到事件的机器人QQ号
	PostType    string `json:"post_type"`    // 上报类型
	RequestType string `json:"request_type"` // 请求类型
	SubType     string `json:"sub_type"`     // 请求子类型, 分别表示加群请求、邀请登录号入群
	GroupID     uint   `json:"group_id"`     // 群号
	UserID      uint   `json:"user_id"`      // 发送请求的QQ号
	Comment     string `json:"comment"`      // 验证信息
	Flag        string `json:"flag"`         // 请求flag, 在调用处理请求的API时需要传入
}

// 群文件上传
type GroupUpload struct {
	Time       int       `json:"time"`        // 事件发生的时间戳
	SelfID     uint      `json:"self_id"`     // 收到事件的机器人QQ号
	PostType   string    `json:"post_type"`   // 上报类型
	NoticeType string    `json:"notice_type"` // 通知类型
	GroupID    uint      `json:"group_id"`    // 群号
	UserID     uint      `json:"user_id"`     // 发送者QQ号
	File       file_info `json:"file"`        // 文件信息
}

// 文件信息
type file_info struct {
	ID    string `json:"id"`    // 文件ID
	Name  string `json:"name"`  // 文件名
	Size  uint   `json:"size"`  // 文件大小 (字节数)
	BusID int    `json:"busid"` // busid(目前不清楚有什么作用)
}

// 群管理员变动
type GroupAdmin struct {
	Time       int    `json:"time"`        // 事件发生的时间戳
	SelfID     uint   `json:"self_id"`     // 收到事件的机器人QQ号
	PostType   string `json:"post_type"`   // 上报类型
	NoticeType string `json:"notice_type"` // 通知类型
	SubType    string `json:"sub_type"`    // 事件子类型, 分别表示设置和取消管理员
	GroupID    uint   `json:"group_id"`    // 群号
	UserID     uint   `json:"user_id"`     // 管理员QQ号
}

// 群成员减少
type GroupDecrease struct {
	Time       int    `json:"time"`        // 事件发生的时间戳
	SelfID     uint   `json:"self_id"`     // 收到事件的机器人QQ号
	PostType   string `json:"post_type"`   // 上报类型
	NoticeType string `json:"notice_type"` // 通知类型
	SubType    string `json:"sub_type"`    // 事件子类型, 分别表示主动退群、成员被踢、登录号被踢
	GroupID    uint   `json:"group_id"`    // 群号
	OperatorID uint   `json:"operator_id"` // 操作者QQ号(如果是主动退群, 则和user_id相同)
	UserID     uint   `json:"user_id"`     // 离开者QQ号
}

// 群成员增加
type GroupIncrease struct {
	Time       int    `json:"time"`        // 事件发生的时间戳
	SelfID     uint   `json:"self_id"`     // 收到事件的机器人QQ号
	PostType   string `json:"post_type"`   // 上报类型
	NoticeType string `json:"notice_type"` // 通知类型
	SubType    string `json:"sub_type"`    // 事件子类型, 分别表示管理员已同意入群、管理员邀请入群
	GroupID    uint   `json:"group_id"`    // 群号
	OperatorID uint   `json:"operator_id"` // 操作者QQ号
	UserID     uint   `json:"user_id"`     // 加入者QQ号
}

// 群禁言
type GroupBan struct {
	Time       int    `json:"time"`        // 事件发生的时间戳
	SelfID     uint   `json:"self_id"`     // 收到事件的机器人QQ号
	PostType   string `json:"post_type"`   // 上报类型
	NoticeType string `json:"notice_type"` // 通知类型
	SubType    string `json:"sub_type"`    // 事件子类型, 分别表示禁言、解除禁言
	GroupID    uint   `json:"group_id"`    // 群号
	OperatorID uint   `json:"operator_id"` // 操作者QQ号
	UserID     uint   `json:"user_id"`     // 被禁言QQ号
	Duration   uint   `json:"duration"`    // 禁言时长, 单位秒
}

// 好友添加
type FriendAdd struct {
	Time       int    `json:"time"`        // 事件发生的时间戳
	SelfID     uint   `json:"self_id"`     // 收到事件的机器人QQ号
	PostType   string `json:"post_type"`   // 上报类型
	NoticeType string `json:"notice_type"` // 通知类型
	UserID     uint   `json:"user_id"`     // 新添加好友QQ号
}

// 群消息撤回
type GroupRecall struct {
	Time       int    `json:"time"`        // 事件发生的时间戳
	SelfID     uint   `json:"self_id"`     // 收到事件的机器人QQ号
	PostType   string `json:"post_type"`   // 上报类型
	NoticeType string `json:"notice_type"` // 通知类型
	GroupID    uint   `json:"group_id"`    // 群号
	UserID     uint   `json:"user_id"`     // 消息发送者QQ号
	OperatorID uint   `json:"operator_id"` // 操作者QQ号
	MessageID  int    `json:"message_id"`  // 被撤回的消息ID
}

// 好友消息撤回
type FriendRecall struct {
	Time       int    `json:"time"`        // 事件发生的时间戳
	SelfID     uint   `json:"self_id"`     // 收到事件的机器人QQ号
	PostType   string `json:"post_type"`   // 上报类型
	NoticeType string `json:"notice_type"` // 通知类型
	UserID     uint   `json:"user_id"`     // 好友QQ号
	MessageID  int    `json:"message_id"`  // 被撤回的消息ID
}

// 戳一戳
//
// * 此事件无法在手表协议上触发
type Poke struct {
	Time       int    `json:"time"`        // 事件发生的时间戳
	SelfID     uint   `json:"self_id"`     // 收到事件的机器人QQ号
	PostType   string `json:"post_type"`   // 上报类型
	NoticeType string `json:"notice_type"` // 通知类型
	SubType    string `json:"sub_type"`    // 提示类型
	GroupID    uint   `json:"group_id"`    // 群号
	SenderID   uint   `json:"sender_id"`   // 发送者QQ号
	UserID     uint   `json:"user_id"`     // 发送者QQ号
	TargetID   uint   `json:"target_id"`   // 被戳者QQ号
}

// 群红包运气王
//
// * 此事件无法在手表协议上触发
type LuckyKing system_notice

// 群成员荣誉变更
//
// * 此事件无法在手表协议上触发
type Honor system_notice

// 系统通知
type system_notice struct {
	Time       int    `json:"time"`
	SelfID     uint   `json:"self_id"`
	PostType   string `json:"post_type"`
	NoticeType string `json:"notice_type"`
	SubType    string `json:"sub_type"`
	SenderID   uint   `json:"sender_id"`
	UserID     uint   `json:"user_id"`
	TargetID   uint   `json:"target_id"`
}

// 群成员名片更新
//
// * 此事件不保证时效性, 仅在收到消息时校验卡片
type GroupCard struct {
	Time       int    `json:"time"`        // 事件发生的时间戳
	SelfID     uint   `json:"self_id"`     // 收到事件的机器人QQ号
	PostType   string `json:"post_type"`   // 上报类型
	NoticeType string `json:"notice_type"` // 通知类型
	GroupID    uint   `json:"group_id"`    // 群号
	UserID     uint   `json:"user_id"`     // 成员id
	NewCard    string `json:"card_new"`    // 新名片
	OldCard    string `json:"card_old"`    // 旧名片
}

// 接收到离线文件
type OfflineFile struct {
	Time       int       `json:"time"`        // 事件发生的时间戳
	SelfID     uint      `json:"self_id"`     // 收到事件的机器人QQ号
	PostType   string    `json:"post_type"`   // 上报类型
	NoticeType string    `json:"notice_type"` // 通知类型
	UserID     uint      `json:"user_id"`     // 发送者QQ号
	File       file_info `json:"file"`        // 文件信息, 离线文件没有ID
}

// 其他客户端在线状态变更
type ClientStatus struct {
	Time       int    `json:"time"`        // 事件发生的时间戳
	SelfID     uint   `json:"self_id"`     // 收到事件的机器人QQ号
	PostType   string `json:"post_type"`   // 上报类型
	NoticeType string `json:"notice_type"` // 通知类型
	Client     client `json:"client"`      // 客户端信息
	Online     bool   `json:"online"`      // 当前是否在线
}

// 客户端信息
type client struct {
	AppID      int    `json:"app_id"`      // 客户端ID
	DeviceName string `json:"device_name"` // 设备名称
	DeviceKind string `json:"device_kind"` // 设备类型
}

// 精华消息
type EssenceMessage struct {
	GroupID    uint   `json:"group_id"`    // 群号
	MessageID  int    `json:"message_id"`  // 消息ID
	NoticeType string `json:"notice_type"` // 消息类型
	OperatorID uint   `json:"operator_id"` // 操作者ID
	PostType   string `json:"post_type"`   // 上报类型
	SelfID     uint   `json:"self_id"`     // BOT QQ号
	SenderID   uint   `json:"sender_id"`   // 消息发送者ID
	SubType    string `json:"sub_type"`    // 添加为add,移出为delete
	Time       int    `json:"time"`        // 事件发生的时间戳
}
