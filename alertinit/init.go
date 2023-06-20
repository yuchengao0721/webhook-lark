package alertinit

import (
	"edge-alert/alertmodel"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Conf alertmodel.Config
var MysqlConf alertmodel.MysqlConfig

// 请勿修改这里的文件地址，是为了保持windows开发环境与docker容器内运行环境配置文件路径一致
const mysqlConfigPath = "./etc/edge-alert/conf/mysql.toml"
const configPath = "./etc/edge-alert/conf/conf.toml"
const logPath = "./etc/edge-alert/log"

func Init() {
	initZeroLog()
	loadConfToml()
	loadMysqlConfToml()
}
func loadConfToml() {
	log.Info().Msgf("加载配置文件开始")
	if _, err := toml.DecodeFile(configPath, &Conf); err != nil {
		log.Error().Msgf("加载配置文件错误: %v", err)
		return
	}
	log.Info().Any("conf", Conf).Msg("加载配置文件成功")
}
func loadMysqlConfToml() {
	log.Info().Msgf("加载categraf配置文件开始")
	if _, err := toml.DecodeFile(mysqlConfigPath, &MysqlConf); err != nil {
		log.Error().Msgf("加载categraf配置文件错误: %v", err)
		return
	}
	log.Info().Any("categraf", MysqlConf).Msg("加载categraf配置文件成功")
}
func initZeroLog() {
	fmt.Println("加载日志文件开始")
	logFilePath := fmt.Sprintf("%s/log.log", logPath)
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		os.MkdirAll(logPath, 0666)
		os.Create(logFilePath)
	}
	logfile, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Error().Msgf("日志文件错误: %v", err)
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = zerolog.New(logfile).With().
		Caller().
		Timestamp().
		Logger()
	// defer logfile.Close()
	fmt.Println("加载日志文件结束")
}
