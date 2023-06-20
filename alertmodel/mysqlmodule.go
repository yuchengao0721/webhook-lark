package alertmodel

type Labels struct {
	Cluster  string `toml:"cluster"`
	Instance string `toml:"instance"`
}
type Instance struct {
	Address  string `toml:"address"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Labels   Labels `toml:"labels"`
}

type MysqlConfig struct {
	Instances []Instance `toml:"instances"`
}
