**级别状态**:  {{.Labels.Level}} [alertmanager]
**规则名称**: {{.Annotations.Title}}
**规则备注**: {{.Annotations.Summary}}
**实例标签**: {{.Annotations.Labels}}
**触发时间**: {{.StartsTime}}
**触发时值**: {{.Annotations.Value}}
**查看详情**: [告警平台](http://alertmanager:9093)