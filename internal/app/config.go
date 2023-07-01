package app

import (
	"github.com/core-go/core"
	"github.com/core-go/log"
	mid "github.com/core-go/log/middleware"
)

type Config struct {
	Server     core.ServerConf `mapstructure:"server"`
	Hive       DBConfig        `mapstructure:"hive"`
	Log        log.Config      `mapstructure:"log"`
	MiddleWare mid.LogConfig   `mapstructure:"middleware"`
}

type DBConfig struct {
	Driver string `mapstructure:"driver" json:"driver,omitempty" gorm:"column:driver" bson:"driver,omitempty" dynamodbav:"driver,omitempty" firestore:"driver,omitempty"`
	Host   string `mapstructure:"host" json:"host,omitempty" gorm:"column:host" bson:"host,omitempty" dynamodbav:"host,omitempty" firestore:"host,omitempty"`
	Port   int    `mapstructure:"port" json:"port,omitempty" gorm:"column:port" bson:"port,omitempty" dynamodbav:"port,omitempty" firestore:"port,omitempty"`
	Auth   string `mapstructure:"auth" json:"auth,omitempty" gorm:"column:auth" bson:"auth,omitempty" dynamodbav:"auth,omitempty" firestore:"auth,omitempty"`
}
