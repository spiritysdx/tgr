package config

type Register struct {
	Name        string // 用户名
	AuthorityId uint   // 权限ID
	// TgBotToken  string // telegram的bot的token
	// ChannelId   string // 上面telegram的bot所在的频道的chat_id
}
