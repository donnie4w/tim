package cluster

import (
	"runtime/debug"
	"time"

	"github.com/donnie4w/go-logger/logger"
	"github.com/garyburd/redigo/redis"
)

var Redis *RedisClient = new(RedisClient)

type RedisClient struct {
	RedisClient *redis.Pool
}

func (this *RedisClient) initPool() {
	this.RedisClient = &redis.Pool{
		MaxIdle:     5,
		MaxActive:   30,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redisConn redis.Conn, err error) {
			redisConn, err = redis.DialTimeout("tcp", Cluster.RedisAddr, 5*time.Second, 2*time.Second, 2*time.Second)
			if err != nil {
				return nil, err
			}
			// 选择db
			if Cluster.RedisPwd != "" {
				redisConn.Do("AUTH", Cluster.RedisPwd)
			}
			redisConn.Do("SELECT", 1)
			return
		},
	}
}

func (this *RedisClient) Do(cmd string, args ...interface{}) (reply interface{}, err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Do", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	logger.Debug("Do ", cmd, " ", args)
	c := this.RedisClient.Get()
	defer c.Close()
	if c != nil {
		reply, err = c.Do(cmd, args...)
		if err != nil {
			logger.Error("Do error ", err.Error(), " ", cmd, " ", args)
		}
	}
	return
}

func (this *RedisClient) RedisCmd(cmd string, retType string, args ...interface{}) (reply interface{}, err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("RedisCmd:", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	logger.Debug("RedisCmd:", cmd, " ", retType, " ", args)
	switch retType {
	case "bool":
		reply, err = redis.Bool(this.Do(cmd, args...))
	case "string":
		reply, err = redis.String(this.Do(cmd, args...))
	case "strings":
		reply, err = redis.Strings(this.Do(cmd, args...))
	case "int":
		reply, err = redis.Int(this.Do(cmd, args...))
	case "int64":
		reply, err = redis.Int64(this.Do(cmd, args...))
	case "float64":
		reply, err = redis.Float64(this.Do(cmd, args...))
	default:
		reply, err = this.Do(cmd, args...)
	}
	return
}

func (this *RedisClient) Ping() bool {
	v, err := redis.String(this.Do("PING"))
	if v == "PONG" && err == nil {
		return true
	}
	return false
}

func (this *RedisClient) IsKeyExsit(key string) (b bool, err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("IsKeyExsit:", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	logger.Debug("IsKeyExsit:", key)
	b, err = redis.Bool(this.Do("EXISTS", key))
	if err != nil {
		b, err = redis.Bool(this.Do("EXISTS", key))
	}
	return
}

//redis setex
func (this *RedisClient) Setex(key string, seconds int, value interface{}) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Setex error:", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	logger.Debug("Setex:", key, " ", seconds, " ", value)
	_, err := this.Do("SETEX", key, seconds, value)
	if err != nil {
		this.Do("SETEX", key, seconds, value)
	}
}

func (this *RedisClient) Set(key string, value interface{}) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Set:", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	logger.Debug("Set:", key, " ", value)
	_, err := this.Do("SET", key, value)
	if err != nil {
		this.Do("SET", key, value)
	}
}

func (this *RedisClient) GetInt(key string) (i int, err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("GetInt error:", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	logger.Debug("Get:", key)
	i, err = redis.Int(this.Do("GET", key))
	if err != nil {
		this.Do("GET", key)
	}
	return
}

func (this *RedisClient) GetString(key string) (s string) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("GetString:", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	logger.Debug("Get:", key)
	var err error
	s, err = redis.String(this.Do("GET", key))
	if err != nil {
		this.Do("GET", key)
	}
	return
}

//自增
func (this *RedisClient) INCR(key string) (i int) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("INCR:", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	logger.Debug("INCR:", key)
	var err error
	i, err = redis.Int(this.Do("INCR", key))
	if err != nil {
		this.Do("INCR", key)
	}
	return
}

//
func (this *RedisClient) Del(key string) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Del:", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	logger.Debug("Del:", key)
	this.Do("DEL", key)
}
