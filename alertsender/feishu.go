package alertsender

import (
	"bytes"
	"edge-alert/alertmodel"
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
const tplPath = "./etc/edge-alert/conf/feishu.tpl"

type FeishuSender struct{}

// 飞书发送消息
func (s *FeishuSender) SendMsg(alertData alertmodel.N9eAlert) bool {
	slowList := GetSlowList(alertData)
	if len(slowList) == 0 {
		return true
	}
	content, err := create_content(slowList)
	if err != nil {
		log.Error().Msgf("模板解析错误:%v", err)
		return false
	}
	// 发送飞书消息
	client := req.C().DevMode()
	var feishu_card alertmodel.FeishuCard
	feishu_card.Msg_type = MsgType
	feishu_card.Card.Header.Title.Tag = TitleTag
	feishu_card.Card.Header.Title.Content = TitleContent
	feishu_card.Card.Header.Template = Template
	feishu_card.Card.Elements = append(feishu_card.Card.Elements, alertmodel.Element{
		Tag: ElementTag, Text: alertmodel.Text{
			Tag:     TextTag,
			Content: content,
		}})
	// 发送给多个飞书机器人
	for _, user := range alertData.NotifyUsersObj {
		token := user.Contacts.Token
		if len(token) > 0 && strings.TrimSpace(token) != "" {
			{
				feishu_url := fmt.Sprintf("https://open.feishu.cn/open-apis/bot/v2/hook/%s", token)
				resp, err := client.R().
					SetHeader(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8).
					SetHeader(fiber.HeaderHost, "open.feishu.cn").
					SetBody(feishu_card).
					Post(feishu_url)
				if err != nil {
					log.Err(err)
				}
				if !resp.IsSuccessState() {
					log.Error().Msgf("飞书通知失败了: %v", err)
					return false
				}
			}

		}
	}
	return true
}

// 拼接通知的消息内容
func create_content(slowList []*alertmodel.MysqlSlowLog) (string, error) {

	content, err := ioutil.ReadFile(tplPath)
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
