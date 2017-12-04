package config

import (
	"fmt"
	"time"

	"goaccount/mysql"

	redis "github.com/garyburd/redigo/redis"
)

type Config struct {
	Server ServerConfig `json:"server"`
	Mysql  MysqlConfig  `json:"mysql"`
	Redis  RedisConfig  `json:"redis"`
	Jwt    JwtConfig    `json:"jwt"`
}

type ServerConfig struct {
	Name      string `json:"server_name" yaml:"server_name"`
	UrlPrefix string `json:"url_prefix" yaml:"url_prefix"`
	ConnLimit int    `json:"connlimit" yaml:"connlimit"`
	ReteLimit int    `json:"ratelimit" yaml:"ratelimit"`
	CertFile  string `json:"cert_file" yaml:"cert_file"`
	KeyFile   string `json:"key_file" yaml:"key_file"`
}

type MysqlConfig struct {
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	DbName   string `json:"db_name" yaml:"db_name"`
}

type RedisConfig struct {
	Password string `json:"password" yaml:"password"`
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
}

type JwtConfig struct {
	Secret    string `json:"secret" yaml:"secret"`
	ExpiresAt int64  `json:"expires_at" yaml:"expires_at"` // Seconds
}

func (this *Config) Load(file string) (err error) {
	return this.loadFromFile(file)
}

func (this *Config) MysqlConfig() *mysql.MysqlConfig {
	return this.dbConfig(this.Mysql.DbName)
}

func (this *Config) dbConfig(dbName string) *mysql.MysqlConfig {
	return &mysql.MysqlConfig{
		User:   this.Mysql.User,
		Pw:     this.Mysql.Password,
		Host:   this.Mysql.Host,
		Port:   this.Mysql.Port,
		DBName: dbName,
	}
}

func (this *Config) RedisPool() *redis.Pool {
	pool := &redis.Pool{
		MaxActive:   50,
		MaxIdle:     5,
		Wait:        true,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", this.Redis.Host, this.Redis.Port))
			if err != nil {
				return nil, err
			}
			if this.Redis.Password != "" {
				if _, err := c.Do("AUTH", this.Redis.Password); err != nil {
					c.Close()
					return nil, err
				}
			}

			if _, err := c.Do("SELECT", 1); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
	}

	return pool
}
