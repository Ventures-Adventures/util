package redis

import (
	//"fmt"
	//"time"
	"sync"
	"errors"
	"context"
	"github.com/go-redis/redis"
	"common/logs"
	"common/lb/hash/chash"
)

var (
	initCHasherOnce sync.Once
	cHasher *chash.Consistent
	connMap map[string]*redis.Client
	client0 *redis.Client
)

func init(){
	connMap = make(map[string]*redis.Client)
}
func GetRedis0CHash()*redis.Client{
	return client0
}
func MembersCHash()[]string{
	return cHasher.Members()
}
func GetRedisCHash(key string)*redis.Client{
	k := cHasher.Get(key)
	return connMap[k]
}
func InitRedisCHash(vnode int, opts []Options)(err error){
	if 0==len(opts){
		return errors.New("InitRedis(), param opts invalid")
	}
	if vnode<=0{
		return errors.New("InitRedis(), param vnode<=0")
	}

	initCHasherOnce.Do(func(){
		cHasher = chash.NewCHash(vnode)
	})

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

		cHasher.Add(opt.Addr)
		connMap[opt.Addr] = client
		if 0==i{
			client0 = client
		}
		logs.Info(context.Background(), "InitRedis success: %s", opt.Addr)
	}
	return nil
}
