package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/penguin-statistics/db-benchmark/utils"
	"github.com/struCoder/pidusage"
)

type Redis struct {
	DB *redis.Client
	pid int
	testCandidates []string
}

func New(pid int) *Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	if err := rdb.FlushDB(context.Background()).Err(); err != nil {
		panic(err)
	}
	return &Redis{DB: rdb, pid: pid}
}

func (r *Redis) Init(rowCount int) (err []error) {
	var errors []error
	for i := 0; i < rowCount; i++ {
		inst := utils.RandomInstallation()
		c := context.Background()
		thisE := r.DB.HSet(c, inst.DeviceToken,
			"server", inst.Server,
			"locale", inst.Locale,
			"preference", inst.ClientPreferences,
			"version", inst.ClientVersion).Err()

		if thisE != nil {
			fmt.Println("error occurred when #Init postgresql:", thisE)
			errors = append(errors, thisE)
			continue
		}

		if utils.SampleHit(0.05) {
			r.testCandidates = append(r.testCandidates, inst.DeviceToken)
		}
	}
	return errors
}

func (r *Redis) Count() int64 {
	res, err := r.DB.Keys(context.Background(), "*").Result()
	if err != nil {
		return -1
	}
	return int64(len(res))
}

func (r *Redis) Quota() uint64 {
	stat, err := pidusage.GetStat(r.pid)
	if err != nil {
		return 0
	}
	return uint64(stat.Memory)
}

func (r *Redis) Read(times int) {
	already := 0
	OUTER:
	for {
		for _, candidate := range r.testCandidates {
			err := r.DB.HGetAll(context.Background(), candidate).Err()
			if err != nil {
				panic(err)
			}
			if already >= times {
				break OUTER
			}
			already += 1
		}
	}
}

func (r *Redis) PrintStats() {
	fmt.Println("===>> Redis current status" +
		"\n	- count of records:", r.Count(),
		"\n	- process resident set size", r.Quota(), "bytes")
}
