package task

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"tick-tock/internal/conf"
	"tick-tock/internal/constant"
	"time"
)

func UnionTimerIDAndRunTime(timerID string, runTime int64) string {
	return fmt.Sprintf("%s_%d", timerID, runTime)
}

func SplitTimerIDAndRunTime(key string) (timerID string, runTime int64, err error) {
	strs := strings.Split(key, "_")
	if len(strs) != 2 {
		return "", 0, fmt.Errorf("invalid key: %s", key)
	}
	if runTime, err = strconv.ParseInt(strs[1], 10, 64); err != nil {
		return "", 0, fmt.Errorf("invalid runTime: %s", strs[1])
	}
	return
}

func GetTableName(ctx context.Context, confData *conf.Data, runTime time.Time, taskID int64) string {
	// key 格式：颗粒度是一分钟、“横向”分桶
	// eg: 2025-07-08 16:36:53_7
	mod := int64(16)
	if confData.Scheduler.BucketCount != 0 {
		mod = int64(confData.Scheduler.BucketCount)
	}
	prefix := runTime.Format(constant.BucketPrefixFormat)
	return fmt.Sprintf("%s_%d", prefix, taskID%mod)
}

func GetLockKey(ctx context.Context, confData *conf.Data, tableName string) string {
	// key 锁格式：prefix_tableName
	// eg: scheduler_2025-07-08 16:36:53_7
	return fmt.Sprintf("%s_%s", confData.Scheduler.LockPrefix, tableName)
}
