package alertmodel

// Config type
type Config struct {
	Application    Application    `toml:"application"`
	Alert          AlertToken     `toml:"alert_tokens"`
	MysqlSlowQuery MysqlSlowQuery `toml:"mysql_slow_query"`
}

// ApplicationApplication type
type Application struct {
	Name string `toml:"name"`
	Port int    `toml:"port"`
}

// 慢查询配置
type MysqlSlowQuery struct {
	Tag            string `toml:"tag"`
	LongQueryTime  int    `toml:"long_query_time"`
	RetrospectTime int    `toml:"retrospect_time"`
}

// Alert type
type AlertToken struct {
	FSToken string `toml:"feishu_rebot_token"`
}
