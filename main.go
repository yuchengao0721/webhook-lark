package main

import (
	"edge-alert/alertinit"
	"edge-alert/alertmodel"
	"edge-alert/alertsender"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func main() {
	fmt.Println("程序启动")
	// 全局设置
	alertinit.Init()
	alertsender.InitializeConnectionPools()
	app := fiber.New(fiber.Config{
		Prefork:       false, //docker环境下千万别开，会导致程序执行闪退
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       alertinit.Conf.Application.Name})
	// 测试是否成功的接口
	app.Get("/healthy", func(c *fiber.Ctx) error {
		return c.SendString("true")
	})
	// 接收logstash传输过来的日志信息
	app.Post("/log/alert", func(c *fiber.Ctx) error {
		log.Info().Msg("接收到一起慢查询报警")
		// 判断是否是合格的告警日志格式
		data := new(alertmodel.N9eAlert)
		if err := c.BodyParser(data); err != nil {
			log.Err(err).Msg("接收到的日志格式好像不大对")
			return c.SendString("true")
		}
		if data.IsRecovered {
			return c.SendString("true")
		}
		sender := new(alertsender.FeishuSender)
		sender.SendMsg(*data)
		return c.SendString("true")
	})
	port := alertinit.Conf.Application.Port
	fmt.Println("程序启动成功")
	app.Listen(":" + strconv.Itoa(port))

}
