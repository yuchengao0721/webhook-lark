package alertmodel

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
