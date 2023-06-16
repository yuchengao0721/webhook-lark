package alertmodel

// ApplicationApplication type
type Application struct {
	Name string `toml:"name"`
	Port int    `toml:"port"`
}

// Alert type
type Alert struct {
	Type    []string `toml:"type"`
	Minutes int      `toml:"minutes"`
}

// Feishu type
type Feishu struct {
	Token []string `toml:"rebot_token"`
}

// Config type
type Config struct {
	Application   Application `toml:"application"`
	Alert         Alert       `toml:"alert"`
	Feishu        Feishu      `toml:"feishu"`
	LongQueryTime int         `toml:"long_query_time"`
}
