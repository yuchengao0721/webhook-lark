package alertmodel

import (
	"strings"
	"time"
)

type GrafanaAlert struct {
	Alerts []Alert `json:"alerts"`
}
type Alert struct {
	Annotations struct {
		Description string `json:"description"`
		Summary     string `json:"summary"`
		Labels      string `json:"labels"`
		Title       string `json:"title"`
		Value       string `json:"value"`
	} `json:"annotations"`
	Labels struct {
		AlertTag  string `json:"AlertTag"`
		Level     string `json:"Level"`
		Namespace string `json:"namespace"`
		Pod       string `json:"pod"`
	} `json:"labels"`
	StartsAt   time.Time `json:"startsAt"`
	StartsTime string
}

// 结构转换
func Convert(ga GrafanaAlert) []Alert {
	result := []Alert{}
	for _, al := range ga.Alerts {
		// 将时间戳转换为中国标准时间
		tm := al.StartsAt.Add(8 * time.Hour)
		formattedTime := tm.Format("2006-01-02 15:04:05")
		al.StartsAt = tm
		al.StartsTime = formattedTime
		// 去掉前面的 "map" 并将 "[]" 替换为 "{}"
		al.Annotations.Labels = strings.TrimPrefix(al.Annotations.Labels, "map")
		al.Annotations.Labels = strings.Replace(al.Annotations.Labels, "[", "{", 1)
		al.Annotations.Labels = strings.Replace(al.Annotations.Labels, "]", "}", strings.LastIndex(al.Annotations.Labels, "]"))
		result = append(result, al)
	}
	return result
}
