// Copyright 2019 Axetroy. All rights reserved. MIT license.
package message_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/module/admin"
	"github.com/axetroy/go-server/module/auth"
	"github.com/axetroy/go-server/module/message"
	"github.com/axetroy/go-server/module/message/message_model"
	"github.com/axetroy/go-server/module/message/message_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCreate(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	// 创建一个消息
	{
		var (
			title   = "test"
			content = "test"
		)

		r := message.Create(schema.Context{
			Uid: adminInfo.Id,
		}, message.CreateMessageParams{
			Uid:     userInfo.Id,
			Title:   title,
			Content: content,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := message_model.Message{}

		assert.Nil(t, tester.Decode(r.Data, &n))

		defer message.DeleteMessageById(n.Id)

		assert.Equal(t, title, n.Title)
		assert.Equal(t, content, n.Content)
	}

	// 非管理员的uid去创建，应该报错
	{
		var (
			title   = "test"
			content = "test"
		)

		r := message.Create(schema.Context{
			Uid: userInfo.Id,
		}, message.CreateMessageParams{
			Uid:     userInfo.Id,
			Title:   title,
			Content: content,
		})

		assert.Equal(t, schema.StatusFail, r.Status)
		assert.Equal(t, admin.ErrAdminNotExist.Error(), r.Message)
	}
}

func TestCreateRouter(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	// 创建一条消息
	{
		var (
			title   = "test"
			content = "test"
		)

		header := mocker.Header{
			"Authorization": token.Prefix + " " + adminInfo.Token,
		}

		body, _ := json.Marshal(&message.CreateMessageParams{
			Uid:     userInfo.Id,
			Title:   title,
			Content: content,
		})

		r := tester.HttpAdmin.Post("/v1/message", body, &header)
		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

		n := message_schema.Message{}

		assert.Nil(t, tester.Decode(res.Data, &n))

		defer message.DeleteMessageById(n.Id)

		assert.Equal(t, title, n.Title)
		assert.Equal(t, content, n.Content)
	}
}
