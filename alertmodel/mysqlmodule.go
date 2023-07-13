package alertmodel

type Labels struct {
	Namespace string `toml:"namespace"`
	Pod       string `toml:"pod"`
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
