package redis

import (
	"fmt"
	"time"
	"errors"
	"strings"
	"strconv"
	"github.com/go-redis/redis"

	"context"
	"common/logs"
	"common/lb/hash"
)


type SentielOptions struct {
	Addrs []string
	Password string
	DB int
	HashBucket int
	GroupPre string
	DialTimeout time.Duration
	ReadTimeout time.Duration
	WriteTimeout time.Duration
	PoolSize int
	PoolTimeout time.Duration
}

var sentinelConns []*redis.Client
var hashBucket int


func GetSentials(key int64) (client *redis.Client) {
	if 0==len(sentinelConns){
		return nil
	}
	idx := key % int64(len(sentinelConns))
	return sentinelConns[idx]
}

func GetSentialInt(key string) (client *redis.Client) {
	if 0==len(sentinelConns){
		return nil
	}

	k, _ := strconv.ParseInt(strings.TrimSpace(key), 10, 64)
	return GetSentials(k)
}

func GetSentialString(key string) (client *redis.Client) {
	if 0==len(sentinelConns){
		return nil
	}

	idx := hash.StrToUint32(key, uint32(len(sentinelConns)))
	return sentinelConns[idx]
}

func InitSential(sOpt SentielOptions) (err error) {
	sentialNum := int64(len(sOpt.Addrs))
	if sentialNum < 1 {
		return errors.New("sentialNum < 1")
	}
	if sOpt.DB < 0 {
		return errors.New("db < 0")
	}
	if sOpt.HashBucket <= 0 {
		return errors.New("sOpt.HashBucket <= 0")
	}
	if "" == sOpt.GroupPre {
		return errors.New("sOpt.GroupPre is empty")
	}

	sentinelConns = make([]*redis.Client, sOpt.HashBucket)
	for i := 0; i < sOpt.HashBucket; i++ {
		failoverOpt := &redis.FailoverOptions{
			MasterName:    sOpt.GroupPre + fmt.Sprintf("%d", i+1),
			SentinelAddrs: sOpt.Addrs,
			OnConnect: func(cn *redis.Conn) error {
				//return cn.ClientSetName("on_connect").Err()
				return nil
			},

			Password: sOpt.Password,
			DB:       sOpt.DB,

			MaxRetries:   2,
			DialTimeout:  sOpt.DialTimeout,
			ReadTimeout:  sOpt.ReadTimeout,
			WriteTimeout: sOpt.WriteTimeout,

			PoolSize:    sOpt.PoolSize,
			PoolTimeout: sOpt.PoolTimeout,
			//IdleTimeout        time.Duration
			//IdleCheckFrequency time.Duration
		}

		sentinelConns[i] = redis.NewFailoverClient(failoverOpt)
		if _, err := sentinelConns[i].Ping().Result(); err != nil {
			logs.Error(context.Background(), "sentinel %d Ping() err:%v", i+1, err)
			return err
		}
	}

	return nil
}
