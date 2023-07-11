package alertsender

import (
	"bytes"
	"edge-alert/alertinit"
	"edge-alert/alertmodel"
	"encoding/json"
	"html/template"
	"io/ioutil"

	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/imroc/req/v3"
	"github.com/rs/zerolog/log"
)

const (
	MsgType      = "interactive"
	TitleTag     = "plain_text"
	TitleContent = "🔔  MySQL慢查询告警"
	Template     = "red"
	ElementTag   = "div"
	TextTag      = "lark_md"
)
const slow_alert_tplPath = "./etc/edge-alert/conf/slow_alert.tpl"
const feishu_tplPath = "./etc/edge-alert/conf/feishu.tpl"

type FeishuSender struct{}

// 飞书发送消息
func (s *FeishuSender) SendMsg(alertData alertmodel.GrafanaAlert) bool {
	var message alertmodel.Message
	err := json.Unmarshal([]byte(alertData.Message), &message)
	if err != nil {
		fmt.Println("Error:", err)
		log.Error().Msgf("Grafana通知媒介配置错误了？: %v", err)
	}
	alertData.MessageObj = message
	//配置项里面的飞书token必填,或者填写grafana里面通知媒介里面Message内的fs_rebot_token值
	fs_tokens := make(alertmodel.Set)
	if strings.TrimSpace(alertData.MessageObj.FSRebotToken) != "" {
		var arr = strings.Split(alertData.MessageObj.FSRebotToken, ",")
		fs_tokens.AddArr(arr)
	}
	if strings.TrimSpace(alertinit.Conf.Alert.FSToken) != "" {
		var arr = strings.Split(alertinit.Conf.Alert.FSToken, ",")
		fs_tokens.AddArr(arr)
	}
	if len(fs_tokens) == 0 {
		return true
	}
	var alerts = alertmodel.Convert(alertData)
	// 常规的告警
	client := req.C().DevMode()
	for _, al := range alerts {
		// 慢查询告警，走另一个通道
		var feishu_card alertmodel.FeishuCard
		if al.Labels.AlertTag == alertinit.Conf.MysqlSlowQuery.Tag {
			slowList := GetSlowList(al)
			if len(slowList) == 0 {
				return true
			}
			content, _ := create_slow_query_alert_content(slowList)
			feishu_card = alertmodel.CreateFsCard("🔔  MySQL慢查询告警", content, "S1")
		} else {
			content, _ := create_common_alert_content(*al)
			if len(content) == 0 {
				return true
			}
			feishu_card = alertmodel.CreateFsCard("⚠️  告警通知", content, al.Labels.Level)
		}
		for token, _ := range fs_tokens {
			if len(token) > 0 {
				feishu_url := fmt.Sprintf("https://open.feishu.cn/open-apis/bot/v2/hook/%s", token)
				resp, err := client.R().
					SetHeader(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8).
					SetHeader(fiber.HeaderHost, "open.feishu.cn").
					SetBody(feishu_card).
					Post(feishu_url)
				if err != nil {
					log.Error().Msgf("飞书通知异常了: %v", err)
				}
				if !resp.IsSuccessState() {
					log.Error().Msgf("飞书通知失败了: %v", err)
				}
			}

		}
	}
	return true
}

// 拼接慢查询通知的消息内容
func create_slow_query_alert_content(slowList []*alertmodel.MysqlSlowLog) (string, error) {
	log.Log().Msgf("查出来的结果是%d", len(slowList))
	if len(slowList) == 0 {
		return "", nil
	}
	content, err := ioutil.ReadFile(slow_alert_tplPath)
	if err != nil {
		fmt.Printf("读取文件失败：%v\n", err)
		return "", err
	}
	tpl := string(content)
	t, err := template.New("feishu").Funcs(template.FuncMap{"ToSeconds": alertmodel.ToSeconds}).Parse(tpl)
	if err != nil {
		log.Error().Msgf("模板加载错误:%v", err)
		return "", err
	}
	// 解析模板
	var buf bytes.Buffer
	// 应用模板并输出结果
	err = t.Execute(&buf, slowList)
	if err != nil {
		log.Error().Msgf("模板执行错误:%v", err)
		return "", err
	}
	return buf.String(), nil
}

// 拼接通知的消息内容
func create_common_alert_content(alertData alertmodel.Alert) (string, error) {
	content, err := ioutil.ReadFile(feishu_tplPath)
	if err != nil {
		fmt.Printf("读取文件失败：%v\n", err)
		return "", err
	}
	tpl := string(content)
	t, err := template.New("feishu").Funcs(template.FuncMap{"ToSeconds": alertmodel.ToSeconds}).Parse(tpl)
	if err != nil {
		log.Error().Msgf("模板加载错误:%v", err)
		return "", err
	}
	// 解析模板
	var buf bytes.Buffer
	// 应用模板并输出结果
	err = t.Execute(&buf, alertData)
	if err != nil {
		log.Error().Msgf("模板执行错误:%v", err)
		return "", err
	}
	return buf.String(), nil
}
