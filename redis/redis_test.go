package redis

import(
	"fmt"
	"time"
	"testing"
)

func init(){
	//opt := Options{Addr:"xxxxx:6379",
	opt := Options{Addr:"www..com:6379",
		Password:"ksd_f123Au29u3p03",
		DB : 0,
		DialTimeout :2*time.Second,
		ReadTimeout :2*time.Second,
		WriteTimeout :2*time.Second,
		PoolSize :1,
		PoolTimeout :60*time.Second,
	}
	err := InitRedis([]Options{opt})
	if err!=nil{
		panic(err)
	}
}

func Test_ttl(t *testing.T){
	key := "aaa"
	r := GetRedisString(key)

	resTTL := r.TTL(key)
	v := resTTL.Val()
	if v==-1*time.Second{
		fmt.Printf("no1\n")
	}
	if v==-2*time.Second{
		fmt.Printf("no2\n")
	}
	fmt.Printf("haha\n")
}


func Test_expire(t *testing.T){
	key := "10086"
	r := GetRedisInt(key)

	timezone:= 60
	r.Incr(key)
	r.Expire(key, time.Duration(timezone)*time.Second)

	resTTL := r.TTL(key)
	v := resTTL.Val()

	fmt.Printf("v:%v\n", v)
}

