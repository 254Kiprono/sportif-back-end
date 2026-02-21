package worker

import (
	"context"

	"webuye-sportif/app/loggers"

	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
)

type RedisWorker struct {
	rdb  *redis.Client
	cron *cron.Cron
}

func NewRedisWorker(rdb *redis.Client) *RedisWorker {
	return &RedisWorker{
		rdb:  rdb,
		cron: cron.New(),
	}
}

func (w *RedisWorker) Start(ctx context.Context) {
	// Schedule to flush Redis every 30 minute for debugging
	_, err := w.cron.AddFunc("@every 30m", func() {
		err := w.rdb.FlushAll(ctx).Err()
		if err != nil {
			loggers.Log.WithField("service", "RedisWorker").Errorf("Failed to flush redis: %v", err)
		} else {
			loggers.Log.WithField("service", "RedisWorker").Info("Redis cleared successfully (Debug mode)")
		}
	})

	if err != nil {
		loggers.Log.Fatalf("Failed to schedule Redis flusher: %v", err)
	}

	w.cron.Start()
	loggers.Log.Info("Redis Debug Worker started (clearing every 30m)")
}

func (w *RedisWorker) Stop() {
	w.cron.Stop()
}
