package alertmodel

import "fmt"

type MysqlSlowLog struct {
	Instance        string
	Db              string
	Query           string
	Query_time      int64
	Last_query_time string
}

func ToSeconds(nanoseconds int64) string {
	seconds := float64(nanoseconds) / 1e12
	return fmt.Sprintf("%.2f", seconds)
}
