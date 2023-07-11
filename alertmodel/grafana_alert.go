package alertmodel

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type GrafanaAlert struct {
	Alerts     []Alert `json:"alerts"`
	Message    string  `json:"message"`
	Title      string  `json:"title"`
	MessageObj Message
}
type Message struct {
	HosrUrl      string `json:"host_url"`
	FSRebotToken string `json:"fs_rebot_token"`
}
type Alert struct {
	Annotations struct {
		Description string `json:"description"`
		Summary     string `json:"summary"`
	} `json:"annotations"`
	Labels struct {
		AlertTag  string `json:"AlertTag"`
		Level     string `json:"Level"`
		AlertName string `json:"alertname"`
		Instance  string `json:"instance"`
	} `json:"labels"`
	StartsAt        time.Time `json:"startsAt"`
	StartsTime      string
	Status          string `json:"status"`
	ValueString     string `json:"valueString"`
	Metric          PrometheusMetric
	HosrUrl         string
	PrometheusLabel string
	FSRebotToken    string
}
type PrometheusMetric struct {
	Metric   string
	Labels   map[string]string
	Labelstr string
	Value    string
}

func parsePrometheusMetric(input string) (PrometheusMetric, error) {
	metric := PrometheusMetric{
		Labels: make(map[string]string),
	}

	// 提取labels和value的正则表达式
	labelsRegex := regexp.MustCompile(`labels={([^}]*)}`)
	valueRegex := regexp.MustCompile(`value=([\d.]+)`)

	// 提取labels和value的值
	labelsMatches := labelsRegex.FindStringSubmatch(input)
	valueMatches := valueRegex.FindStringSubmatch(input)

	if len(labelsMatches) != 2 || len(valueMatches) != 2 {
		log.Error().Msg("告警标签缺失")
	}
	labelsStr := labelsMatches[1]
	valueStr := valueMatches[1]
	value, _ := strconv.ParseFloat(valueStr, 64)
	metric.Value = fmt.Sprintf("%.2f", value)
	metric.Labelstr = labelsStr
	// 解析标签
	labelParts := strings.Split(labelsStr, ",")
	for _, part := range labelParts {
		kv := strings.Split(part, "=")
		key := strings.TrimSpace(kv[0])
		value := strings.Trim(kv[1], "\"")
		metric.Labels[key] = value
	}
	return metric, nil
}

// 结构转换
func Convert(ga GrafanaAlert) []*Alert {
	result := []*Alert{}
	for _, al := range ga.Alerts {
		var metric, _ = parsePrometheusMetric(al.ValueString)
		// 将时间戳转换为中国标准时间
		tm := al.StartsAt.Add(8 * time.Hour)
		formattedTime := tm.Format("2006-01-02 15:04:05")
		al.StartsAt = tm
		al.StartsTime = formattedTime
		al.Metric = metric
		al.HosrUrl = strings.ReplaceAll(ga.MessageObj.HosrUrl, " ", "")
		result = append(result, &al)
	}
	return result
}
