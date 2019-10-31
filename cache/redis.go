package cache

import (
    "github.com/garyburd/redigo/redis"
    "log"
)

var RedisPool *redis.Pool

func InitRedis(host string, port string, auth string) {
    redisAddress := host + ":" + port
    RedisPool = &redis.Pool{
        //MaxIdle:     16,
        //MaxActive:   0,
        IdleTimeout: 300,
        Wait:        true,
        Dial: func() (redis.Conn, error) {
            r, err := redis.Dial("tcp", redisAddress)
            if auth != "none" {
                _, err := r.Do("AUTH", auth)
                if err != nil {
                    log.Fatalln("redis鉴权失败")
                }
            }
            return r, err
        },
    }
}
