/**
 * donnie4w@gmail.com  tim server
 */
package cluster

import (
	"errors"
	"runtime/debug"
	"time"

	"github.com/donnie4w/go-logger/logger"
	"github.com/garyburd/redigo/redis"
	. "tim.common"
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
			redisConn, err = redis.DialTimeout("tcp", ClusterConf.RedisAddr, 5*time.Second, 2*time.Second, 2*time.Second)
			if err != nil {
				return nil, err
			}
			if ClusterConf.RedisPwd != "" {
				redisConn.Do("AUTH", ClusterConf.RedisPwd)
			}
			redisConn.Do("SELECT", ClusterConf.RedisDB)
			return
		},
	}
}

func (this *RedisClient) initRedisClient(addr, pwd string) {
	this.RedisClient = &redis.Pool{
		MaxIdle:     5,
		MaxActive:   30,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redisConn redis.Conn, err error) {
			redisConn, err = redis.DialTimeout("tcp", addr, 5*time.Second, 2*time.Second, 2*time.Second)
			if err != nil {
				return nil, err
			}
			if ClusterConf.RedisPwd != "" {
				redisConn.Do("AUTH", pwd)
			}
			redisConn.Do("SELECT", ClusterConf.RedisDB)
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

func (this *RedisClient) Ping() (bool, error) {
	v, err := redis.String(this.Do("PING"))
	if err != nil {
		return false, err
	}
	if v == "PONG" && err == nil {
		return true, nil
	} else {
		return false, errors.New("ping failed")
	}
}

func (this *RedisClient) IsKeyExsit(key string) (b bool, err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("IsKeyExsit:", err)
			logger.Error(string(debug.Stack()))
		}
	}()
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
	var err error
	i, err = redis.Int(this.Do("INCR", key))
	if err != nil {
		this.Do("INCR", key)
	}
	return
}

//
func (this *RedisClient) Del(key string) (err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Del:", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	_, err = this.Do("DEL", key)
	return
}

func (this *RedisClient) Hkeys(key string) (reply []string, err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Hkeys:", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	reply, err = redis.Strings(this.Do("HKEYS", key))
	return
}

func (this *RedisClient) Hget(key, field string) (reply string, err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Hget:", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	reply, err = redis.String(this.Do("HGET", key, field))
	return
}

func (this *RedisClient) Hvals(key string) (reply []string, err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Hvals:", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	reply, err = redis.Strings(this.Do("HVALS", key))
	return
}

/**
 */

func (this *RedisClient) Hset(key, field, value string) (err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Hset:", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	_, err = this.Do("HSET", key, field, value)
	return
}

/**
 */

func (this *RedisClient) Hdel(key, field string) (err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("HDEL:", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	_, err = this.Do("HDEL", key, field)
	return
}

type ScriptCmd struct {
	script *redis.Script
}

func (this *ScriptCmd) EvalSha(arg ...interface{}) (reply interface{}, err error) {
	c := Redis.RedisClient.Get()
	defer c.Close()
	reply, err = this.script.Do(c, arg...)
	return
}

func (this *ScriptCmd) EvalShaStrings(arg ...interface{}) (reply []string, err error) {
	c := Redis.RedisClient.Get()
	defer c.Close()
	reply, err = redis.Strings(this.script.Do(c, arg...))
	return
}

func (this *RedisClient) NewScript(scriptStr string, keycount int) (scriptcmd *ScriptCmd) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("NewScript:", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	c := this.RedisClient.Get()
	defer c.Close()
	if c != nil {
		scriptcmd = &ScriptCmd{script: redis.NewScript(keycount, scriptStr)}
	}
	return
}

var scriptAddCmd = `
redis.call('hset',KEYS[1],KEYS[2],0)
redis.call('setex',KEYS[2],ARGV[1],ARGV[2])
return 1
`

var scriptGetCmd = `
 local list = redis.call('Hkeys',KEYS[1])
 local array = {}
 for _,v in pairs(list) do
        local vv = redis.call('get',v)
        if vv == nil or vv == '(nil)' or vv == false then 
           redis.call('hdel',KEYS[1],v)
        else
           table.insert(array,vv)
        end 
 end
 if table.maxn(array) == 0 then
    return 0
 else 	
    return array
 end
`
