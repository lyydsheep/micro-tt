package cron

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"testing"
	"time"
)

func TestExp(t *testing.T) {
	expr := "*/15 * * * *"
	start := time.Now()
	end := start.Add(time.Hour)
	for nextTime := start; nextTime.Before(end); {
		fmt.Println(nextTime)
		nextTime = cronexpr.MustParse(expr).Next(nextTime)
	}
}
