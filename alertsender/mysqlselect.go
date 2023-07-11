package alertsender

import (
	"database/sql"
	"edge-alert/alertinit"
	"edge-alert/alertmodel"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	_ "github.com/go-sql-driver/mysql"
)

var Pool *MySQLConnectionPools

type MySQLConnectionPools struct {
	pools map[string]*sql.DB
	mu    sync.Mutex
}

func NewMySQLConnectionPools() *MySQLConnectionPools {
	return &MySQLConnectionPools{
		pools: make(map[string]*sql.DB),
		mu:    sync.Mutex{},
	}
}

func (cp *MySQLConnectionPools) GetMySQLConnectionPool(instanceName string) (*sql.DB, error) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	pool, ok := cp.pools[instanceName]
	if ok {
		return pool, nil
	}

	return nil, fmt.Errorf("connection pool '%s' not found", instanceName)
}

func InitializeConnectionPools() {
	var config = alertinit.MysqlConf
	connectionPools := NewMySQLConnectionPools()
	for _, instance := range config.Instances {
		if instance.Labels.Instance != "" {
			db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/performance_schema", instance.Username, instance.Password, instance.Address))
			if err != nil {
				log.Error().Msgf("failed to open MySQL connection for instance '%s': %v", fmt.Sprintf("%s-%s", instance.Labels.Cluster, instance.Labels.Instance), err)
			}
			// 设置连接池最大连接数为 5
			db.SetMaxOpenConns(5)
			connectionPools.pools[fmt.Sprintf("%s-%s", instance.Labels.Cluster, instance.Labels.Instance)] = db
		}
	}
	Pool = connectionPools
}

func GetSlowList(data *alertmodel.Alert) []*alertmodel.MysqlSlowLog {
	slowList := []*alertmodel.MysqlSlowLog{}
	cluster := data.Metric.Labels["cluster"]
	instance := data.Metric.Labels["instance"]
	//不存在cluster和instance的话，就没法找到对应的数据库
	// if !cexists || !iexists {
	// 	return slowList
	// }
	conn, err := Pool.GetMySQLConnectionPool(fmt.Sprintf("%s-%s", cluster, instance))
	if err != nil {
		log.Error().Msgf("Failed to get connection:%s", err)
		return slowList
	}
	tm := data.StartsAt.Add(-(time.Duration(alertinit.Conf.MysqlSlowQuery.LongQueryTime)) * time.Minute)
	formattedTime := tm.Format("2006-01-02 15:04:05")
	sql := fmt.Sprintf("SELECT `SCHEMA_NAME` as 'db',`QUERY_SAMPLE_TIMER_WAIT` as 'query_time',`QUERY_SAMPLE_TEXT` as 'query',`LAST_SEEN` as 'last_query_time'FROM events_statements_summary_by_digest  where `QUERY_SAMPLE_TIMER_WAIT` > %d *1000000000000 AND `LAST_SEEN` > '%s' ORDER BY LAST_SEEN DESC", alertinit.Conf.MysqlSlowQuery.LongQueryTime, formattedTime)
	logData := map[string]interface{}{
		"tag":          fmt.Sprintf("%s-%s", cluster, instance),
		"lastEvalTime": data.StartsAt,
		"tm":           tm,
		"sql":          sql,
	}
	log.Info().Fields(logData).Msg("慢查询信息")
	rows, err := conn.Query(sql)
	if err != nil {
		log.Error().Msgf("Failed to execute query:%s", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var slow alertmodel.MysqlSlowLog
			err := rows.Scan(&slow.Db, &slow.Query_time, &slow.Query, &slow.Last_query_time)
			if err != nil {
				log.Error().Msgf("Failed to scan row:%s", err)
			} else {
				slow.Instance = instance
				slowList = append(slowList, &slow)
			}
		}
	}
	return slowList
}
