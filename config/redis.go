// Copyright 2019 Axetroy. All rights reserved. MIT license.
package config

import (
	"github.com/axetroy/go-server/service/dotenv"
)

type redis struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
}

var Redis redis

func init() {
	if Redis.Host = dotenv.Get("REDIS_SERVER"); Redis.Host == "" {
		Redis.Host = "127.0.0.1"
	}
	if Redis.Port = dotenv.Get("REDIS_PORT"); Redis.Port == "" {
		Redis.Port = "6379"
	}
	if Redis.Password = dotenv.Get("REDIS_PASSWORD"); Redis.Password == "" {
		Redis.Password = "password"
	}
}
