// Copyright 2019 Axetroy. All rights reserved. MIT license.
package banner_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/module/admin"
	"github.com/axetroy/go-server/module/auth"
	"github.com/axetroy/go-server/module/banner"
	"github.com/axetroy/go-server/module/banner/banner_model"
	"github.com/axetroy/go-server/module/banner/banner_schema"
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

	// 创建一个 Banner
	{
		var (
			image    = "test.png"
			href     = "https://example.com"
			platform = banner_model.BannerPlatformApp
		)

		r := banner.Create(schema.Context{
			Uid: adminInfo.Id,
		}, banner.CreateParams{
			Image:    image,
			Href:     href,
			Platform: platform,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := banner_schema.Banner{}

		assert.Nil(t, tester.Decode(r.Data, &n))

		defer banner.DeleteBannerById(n.Id)

		assert.Equal(t, image, n.Image)
		assert.Equal(t, href, n.Href)
		assert.Equal(t, platform, n.Platform)
	}

	// 非管理员的uid去创建，应该报错
	{

		userInfo, _ := tester.CreateUser()

		defer auth.DeleteUserByUserName(userInfo.Username)

		var (
			image    = "test.png"
			href     = "https://example.com"
			platform = banner_model.BannerPlatformApp
		)

		r := banner.Create(schema.Context{
			Uid: userInfo.Id,
		}, banner.CreateParams{
			Image:    image,
			Href:     href,
			Platform: platform,
		})

		assert.Equal(t, schema.StatusFail, r.Status)
		assert.Equal(t, admin.ErrAdminNotExist.Error(), r.Message)
	}
}

func TestCreateRouter(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	// 创建 banner
	{
		var (
			image    = "test.png"
			href     = "https://example.com"
			platform = banner_model.BannerPlatformApp
		)

		header := mocker.Header{
			"Authorization": token.Prefix + " " + adminInfo.Token,
		}

		body, _ := json.Marshal(&banner.CreateParams{
			Image:    image,
			Href:     href,
			Platform: platform,
		})

		r := tester.HttpAdmin.Post("/v1/banner", body, &header)
		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

		n := banner_schema.Banner{}

		assert.Nil(t, tester.Decode(res.Data, &n))

		defer banner.DeleteBannerById(n.Id)

		assert.Equal(t, image, n.Image)
		assert.Equal(t, href, n.Href)
		assert.Equal(t, platform, n.Platform)
	}
}
