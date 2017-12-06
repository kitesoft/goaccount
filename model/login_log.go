package model

import (
	"encoding/json"
	"goaccount/mixin"

	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"
)

//使用redis 的list
type LoginLog struct {
	CreatedAt int64           `json:"created_at"`
	Name      string          `json:"name"`
	Ip        string          `json:"ip"`
	Ua        string          `json:"ua"`
	Result    mixin.ErrorCode `json:"result"`
}

func AddLoginLog(loginLog LoginLog) {
	conn := Pool.Get()

	logJSON, _ := json.Marshal(loginLog)

	key := "loginLog:" + loginLog.Name

	_, err := conn.Do("LPUSH", key, logJSON)
	if err != nil {
		logrus.Debugf("[model.AddLoginLog] redis.do %s", err.Error())
		return
	}

	_, err = conn.Do("LTRIM", key, 0, 29)
	if err != nil {
		logrus.Debugf("[model.AddLoginLog] redis.do %s", err.Error())
		return
	}

}

func QueryLoginLog(userName string) []string {
	conn := Pool.Get()
	key := "loginLog:" + userName

	resp, err := redis.Strings(conn.Do("LRANGE", key, 0, "end"))
	if err != nil {
		logrus.Debugf("[model.AddLoginLog] redis.do %s", err.Error())
		return nil
	}

	return resp
}
