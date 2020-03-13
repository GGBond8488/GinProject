package gredis

import (
	"My-gin-Project/pkg/logging"
	"My-gin-Project/pkg/setting"
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"time"
)

//Dial：提供创建和配置应用程序连接的一个函数
//TestOnBorrow：可选的应用程序检查健康功能
//MaxIdle：最大空闲连接数
//MaxActive：在给定时间内，允许分配的最大连接数（当为零时，没有限制）
//IdleTimeout：在给定时间内将会保持空闲状态，若到达时间限制则关闭连接（当为零时，没有限制
var ConnPool *redis.Pool

func Setup()error {
	ConnPool = &redis.Pool{
		MaxIdle: setting.RedisSetting.MaxIdle,
		MaxActive:setting.RedisSetting.MaxActive,
		IdleTimeout:setting.RedisSetting.IdleTimeout,
		Dial: func() (conn redis.Conn, err error) {
			c ,err := redis.Dial("tcp",setting.RedisSetting.Host)
			if err != nil{
				logging.Warn(err)
				return nil,err
			}else if _ , err := c.Do("AUTH",setting.RedisSetting.Password);err != nil{
				logging.Warn(err)
				return nil,err
			}
			return c,err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return nil
}
//（1）RedisConn.Get()：在连接池中获取一个活跃连接
//
//（2）conn.Do(commandName string, args ...interface{})：向 Redis 服务器发送命令并返回收到的答复
//
//（3）redis.Bool(reply interface{}, err error)：将命令返回转为布尔值
//
//（4）redis.Bytes(reply interface{}, err error)：将命令返回转为 Bytes
//
//（5）redis.Strings(reply interface{}, err error)：将命令返回转为 []string

func Set(key string,data interface{},time int) error {
	conn := ConnPool.Get()
	defer conn.Close()

	value,err := json.Marshal(data)
	if err !=nil {
		return err
	}
	_, err = conn.Do("SET",key,value)
	if err != nil {
		return err
	}

	_ , err = conn.Do("EXPIRE",key,time)
	if err != nil{
		return err
	}
	return nil
}

func Exists(key string)bool{
	conn := ConnPool.Get()
	defer conn.Close()
	exists,err := redis.Bool(conn.Do("EXISTS",key))
	if err != nil {
		//logging.Warn()
		return false
	}
	return exists
}

func Get(key string)([]byte,error)  {
	conn := ConnPool.Get()
	defer conn.Close()
	reply,err := redis.Bytes(conn.Do("GET",key))
	if err != nil{
		return nil,err
	}
	return reply,nil
}

func Delete(key string)(bool,error)  {
	conn := ConnPool.Get()
	defer conn.Close()
	return redis.Bool(conn.Do("DEL",key))
}

func LikeDeletes(key string)error  {
	conn := ConnPool.Get()
	defer conn.Close()
	keys,err := redis.Strings(conn.Do("keys","*"+key+"*"))
	if err != nil{
		for _,key := range keys{
			_,err = Delete(key)
			if err != nil{
				return err
			}
		}
	}
	return nil
}