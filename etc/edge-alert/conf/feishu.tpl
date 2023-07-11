**级别状态**:  {{.Labels.Level}}
**规则名称**: {{.Labels.AlertName}}
**规则备注**: {{.Annotations.Summary}}
**实例标签**: {{ "{" }}{{.Metric.Labelstr}}{{ "}" }}
**触发时间**: {{.StartsTime}}
**触发时值**: {{.Metric.Value}}
**查看详情**: [告警平台]({{.HosrUrl}})