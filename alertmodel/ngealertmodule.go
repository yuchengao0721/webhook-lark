package alertmodel

type Annotation struct {
	FeishuAts string `json:"FeishuAts"`
}

type NotifyUser struct {
	Contacts FsToken `json:"contacts"`
}
type FsToken struct {
	Token string `json:"feishu_robot_token"`
}

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type N9eAlert struct {
	Annotations    Annotation   `json:"annotations"`
	IsRecovered    bool         `json:"is_recovered"`
	LastEvalTime   int64        `json:"last_eval_time"`
	NotifyUsersObj []NotifyUser `json:"notify_users_obj"`
	RuleName       string       `json:"rule_name"`
	RuleNote       string       `json:"rule_note"`
	Tags           []string     `json:"tags"`
}
