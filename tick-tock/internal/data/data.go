package data

import (
	"context"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"tick-tock/internal/biz"
	"tick-tock/internal/conf"
	"tick-tock/internal/data/gen"
	"tick-tock/pkg/log"
	"tick-tock/util/sql"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewQuery, NewRedis, NewData, NewTaskRepo, NewTaskDefineRepo, NewTransaction,
	NewLock, NewTaskCache)

// 传递 transaction db
const contextVal = "transaction"

// Data .
type Data struct {
	query *gen.Query
	redis *redis.Client
}

// NewData .
func NewData(query *gen.Query, redis *redis.Client) (*Data, func(), error) {
	cleanup := func() {
		log.Info(nil, "closing the data resources.")
	}
	return &Data{
		query: query,
		redis: redis,
	}, cleanup, nil
}

func NewTransaction(query *gen.Query, redis *redis.Client) biz.Transaction {
	return &Data{
		query: query,
		redis: redis,
	}
}

func (d *Data) Txn(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.query.Transaction(func(tx *gen.Query) error {
		ctx = context.WithValue(ctx, contextVal, tx)
		return fn(ctx)
	})
}

func (d *Data) DB(ctx context.Context) *gen.Query {
	q, ok := ctx.Value(contextVal).(*gen.Query)
	if ok {
		return q
	}
	return d.query
}

func NewQuery(data *conf.Data) *gen.Query {
	dsn := sql.GetDSN(data.Database.Username, data.Database.Password, data.Database.Addr, data.Database.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: log.Logger(),
	})
	if err != nil {
		log.Fatal(nil, "fail to connect db.", "error", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(nil, "fail to get db.", "error", err)
	}
	// 最大（长）连接数
	sqlDB.SetMaxOpenConns(int(data.Database.MaxOpenConn))
	// 最大空闲链接，默认为 2
	sqlDB.SetMaxIdleConns(int(data.Database.MaxIdleConn))
	// 最大存活时间
	sqlDB.SetConnMaxLifetime(data.Database.ConnMaxIdleTime.AsDuration())

	log.Info(nil, "connect db success.", "maxConn", data.Database.MaxOpenConn,
		"maxIdleConn", data.Database.MaxIdleConn, "maxIdleTime_s", data.Database.ConnMaxIdleTime.AsDuration())
	return gen.Use(db)
}

func NewRedis(data *conf.Data) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:        data.Redis.Addr,
		ReadTimeout: data.Redis.ReadTimeout.AsDuration(),
	})
	return client
}
