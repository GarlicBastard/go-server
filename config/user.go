// Copyright 2019 Axetroy. All rights reserved. MIT license.
package config

import (
	"github.com/axetroy/go-server/service/dotenv"
)

type user struct {
	Domain string `json:"domain"` // 用户端 API 绑定的域名, 例如 https://example.com
	Port   string `json:"port"`   // 用户端 API 监听的端口
	Secret string `json:"secret"` // 用户端密钥，用于加密/解密 token
}

var User user

func init() {
	if User.Port = dotenv.Get("USER_HTTP_PORT"); User.Port == "" {
		User.Port = "8080"
	}
	if User.Domain = dotenv.Get("USER_HTTP_DOMAIN"); User.Domain == "" {
		User.Domain = "http://127.0.0.1:" + User.Port
	}
	if User.Secret = dotenv.Get("USER_TOKEN_SECRET_KEY"); User.Secret == "" {
		User.Secret = "user"
	}
}
