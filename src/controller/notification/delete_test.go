// Copyright 2019 Axetroy. All rights reserved. MIT license.
package notification_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/admin"
	"github.com/axetroy/go-server/src/controller/notification"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/database"
	"github.com/axetroy/go-server/src/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestDelete(t *testing.T) {
	var (
		adminUid string
	)
	// 先登陆获取管理员的Token
	{
		// 登陆超级管理员-成功

		r := admin.Login(admin.SignInParams{
			Username: "admin",
			Password: "admin",
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		adminInfo := schema.AdminProfileWithToken{}

		if err := tester.Decode(r.Data, &adminInfo); err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, "admin", adminInfo.Username)
		assert.True(t, len(adminInfo.Token) > 0)

		if c, er := token.Parse(token.Prefix+" "+adminInfo.Token, true); er != nil {
			t.Error(er)
		} else {
			adminUid = c.Uid
		}
	}

	context := controller.Context{
		Uid: adminUid,
	}

	var testNotification schema.Notification

	// 创建一篇系统通知
	{
		var (
			title   = "TestUpdate"
			content = "TestUpdate"
		)

		r := notification.Create(context, notification.CreateParams{
			Title:   title,
			Content: content,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		testNotification = schema.Notification{}

		assert.Nil(t, tester.Decode(r.Data, &testNotification))

		defer notification.DeleteNotificationById(testNotification.Id)

		assert.Equal(t, title, testNotification.Title)
		assert.Equal(t, content, testNotification.Content)
	}

	// 获取通知
	{
		n := model.Notification{
			Id: testNotification.Id,
		}

		assert.Nil(t, database.Db.Model(&n).Where(&n).First(&n).Error)
	}

	// 删除通知
	{
		r := notification.Delete(context, testNotification.Id)

		assert.Equal(t, "", r.Message)
		assert.Equal(t, schema.StatusSuccess, r.Status)
	}

	// 再次获取通知，这时候通知应该已经被删除了
	{
		n := model.Notification{
			Id: testNotification.Id,
		}

		if err := database.Db.Model(&n).Where(&n).First(&n).Error; err != nil {
			assert.Equal(t, gorm.ErrRecordNotFound.Error(), err.Error())
		} else {
			assert.Fail(t, "通知应该已被删除")
		}
	}
}

func TestDeleteRouter(t *testing.T) {
	var (
		adminToken       string
		notificationInfo = schema.Notification{}
	)
	// 先登陆获取管理员的Token
	{
		// 登陆超级管理员-成功

		r := admin.Login(admin.SignInParams{
			Username: "admin",
			Password: "admin",
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		adminInfo := schema.AdminProfileWithToken{}

		if err := tester.Decode(r.Data, &adminInfo); err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, "admin", adminInfo.Username)
		assert.True(t, len(adminInfo.Token) > 0)

		if _, er := token.Parse(token.Prefix+" "+adminInfo.Token, true); er != nil {
			t.Error(er)
		} else {
			adminToken = adminInfo.Token
		}
	}

	header := mocker.Header{
		"Authorization": token.Prefix + " " + adminToken,
	}

	// 创建一条系统通知
	{
		var (
			title   = "test title"
			content = "test content"
		)

		body, _ := json.Marshal(&notification.CreateParams{
			Title:   title,
			Content: content,
		})

		r := tester.HttpAdmin.Post("/v1/notification", body, &header)
		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal(r.Body.Bytes()), &res))

		assert.Nil(t, tester.Decode(res.Data, &notificationInfo))

		defer notification.DeleteNotificationById(notificationInfo.Id)

		assert.Equal(t, title, notificationInfo.Title)
		assert.Equal(t, content, notificationInfo.Content)
	}

	// 删除这条通知
	{
		r := tester.HttpAdmin.Delete("/v1/notification/n/"+notificationInfo.Id, nil, &header)

		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal(r.Body.Bytes()), &res))

		// 再查找这条记录，应该是空的

		n := model.Notification{Id: notificationInfo.Id}

		err := database.Db.Where(&n).First(&n).Error

		assert.NotNil(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound.Error(), err.Error())
	}
}
