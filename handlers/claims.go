package handlers

import (
	"github.com/golang-jwt/jwt/v4"
)

type config struct {
	ApiUrl   string                `yaml:"apiurl"`
	Username string                `yaml:"username"`
	Password string                `yaml:"password"`
	Devices  map[string]deviceInfo `yaml:"devices"`
}

type deviceInfo struct {
	Newid   string            `yaml:"newid"`
	Sensors map[string]string `yaml:"sensors"`
}

type tokenAuthReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type accesstoken struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Token   string `json:"token"`
}

type usernameClaims struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
}

type sensordataReq struct {
	Sensorinfo    string `json:"sensorinfo"`
	Concentration int32  `json:"concentration"`
	Unit          string `json:"unit"`
}

type realtimeDataReq struct {
	Updatetime int64           `json:"updatetime"`
	Device     string          `json:"device"`
	Sensordata []sensordataReq `json:"sensordata"`
}

type restructuredData struct {
	Did     string        `json:"did"`   // 主机id----可配置
	Utime   string        `json:"utime"` //采集时间----"updatetime": new Date(), //实时数据的时间戳，自 Epoch 以来的毫秒数
	Content []dataContent `json:"content"`
}

type dataContent struct {
	Pid   string `json:"pid"`   //默认：1
	Type  string `json:"type"`  //type=0代表模拟量，type=1代表开关量
	Addr  string `json:"addr"`  //iot传感id------HY001LEL01---可配置
	Addrv string `json:"addrv"` //采集值-----"concentration": 7190,
	Ctime string `json:"ctime"` //采集时间---"updatetime": new Date(), //实时数据的时间戳，自 Epoch 以来的毫秒数
}
