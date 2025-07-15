package sql

import (
	"fmt"
)

func GetDSN(username, password, addr, dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=True&loc=UTC", username, password, addr, dbName)
}
