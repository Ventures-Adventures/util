package redis

import (
	"errors"
	"context"
	"strings"
	"strconv"
	"github.com/go-redis/redis"
	"common/logs"
	"common/lb/hash"
)

var conns []*redis.Client

func GetRedis(key int64) (client *redis.Client) {
	if 0==len(conns){
		return nil
	}
	idx := key % int64(len(conns))
	return conns[idx]
}

func GetRedisInt(key string) (client *redis.Client) {
	if 0==len(conns){
		return nil
	}

	k, _ := strconv.ParseInt(strings.TrimSpace(key), 10, 64)
	return GetRedis(k)
}

func GetRedisString(key string) (client *redis.Client) {
	if 0==len(conns){
		return nil
	}

	idx := hash.StrToUint32(key, uint32(len(conns)))
	return conns[idx]
}

func InitRedis(opts []Options)(err error){
	if 0==len(opts){
		return errors.New("InitRedis(), param opt invalid")
	}

	conns = make([]*redis.Client, len(opts))
	for i, opt := range opts{
		client := redis.NewClient(&redis.Options{
			Addr:         opt.Addr,
			Password:     opt.Password,
			DB:           opt.DB,
			DialTimeout:  opt.DialTimeout,
			ReadTimeout:  opt.ReadTimeout,
			WriteTimeout: opt.WriteTimeout,
			PoolSize:     opt.PoolSize,
			PoolTimeout:  opt.PoolTimeout,
		})

		if _, err := client.Ping().Result(); err != nil {
			return err
		}

		logs.Info(context.Background(), "InitRedis success: %s", opt.Addr)
		conns[i] = client
	}
	return nil
}
