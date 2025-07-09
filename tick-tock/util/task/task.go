package task

import (
	"fmt"
	"strconv"
	"strings"
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
