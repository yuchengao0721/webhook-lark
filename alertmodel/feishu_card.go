package alertmodel

const (
	MsgType    = "interactive"
	TitleTag   = "plain_text"
	ElementTag = "div"
	TextTag    = "lark_md"
)

// 创建一个飞书卡片对象
func CreateFsCard(title, content, level string) FeishuCard {
	colorMap := map[string]string{
		"S1": "red",
		"S2": "yellow",
		"S3": "purple",
	}
	var feishu_card FeishuCard
	feishu_card.Msg_type = MsgType
	feishu_card.Card.Header.Title.Tag = TitleTag
	feishu_card.Card.Header.Title.Content = title
	feishu_card.Card.Header.Template = colorMap[level]
	feishu_card.Card.Elements = append(feishu_card.Card.Elements, Element{
		Tag: ElementTag, Text: Text{
			Tag:     TextTag,
			Content: content,
		}})
	return feishu_card
}

type FeishuCard struct {
	Msg_type string    `json:"msg_type"`
	Card     CardBlock `json:"card"`
}

type CardBlock struct {
	Elements []Element `json:"elements"`
	Header   Header    `json:"header"`
}

type Element struct {
	Tag  string `json:"tag"`
	Text Text   `json:"text"`
}

type Text struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type Header struct {
	Title    Text   `json:"title"`
	Template string `json:"template"`
}
