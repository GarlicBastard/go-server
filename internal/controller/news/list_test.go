// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package news_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/controller"
	"github.com/axetroy/go-server/internal/controller/news"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetList(t *testing.T) {
	// 现在没有任何文章，获取到的应该是0个长度的
	{
		var (
			data = make([]model.News, 0)
		)
		query := schema.Query{
			Limit: 20,
		}
		r := news.GetNewsListByUser(news.Query{
			Query: query,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		assert.Nil(t, tester.Decode(r.Data, &data))
		assert.Equal(t, query.Limit, r.Meta.Limit)
		assert.Equal(t, schema.DefaultPage, r.Meta.Page)
	}

	adminInfo, _ := tester.LoginAdmin()

	// 2. 先创建一篇新闻作为测试
	{
		var (
			title    = "test"
			content  = "test"
			newsType = model.NewsTypeNews
		)

		r := news.Create(controller.Context{
			Uid: adminInfo.Id,
		}, news.CreateNewParams{
			Title:   title,
			Content: content,
			Type:    newsType,
			Tags:    []string{},
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := schema.News{}

		assert.Nil(t, tester.Decode(r.Data, &n))

		defer news.DeleteNewsById(n.Id)
	}

	// 3. 获取列表
	{
		var (
			data = make([]model.News, 0)
		)
		query := schema.Query{
			Limit: 20,
		}
		r := news.GetNewsListByUser(news.Query{
			Query: query,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		assert.Nil(t, tester.Decode(r.Data, &data))
		assert.Equal(t, query.Limit, r.Meta.Limit)
		assert.Equal(t, schema.DefaultPage, r.Meta.Page)

		assert.True(t, len(data) >= 1)
	}
}

func TestGetListRouter(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	{
		var (
			title    = "test"
			content  = "test"
			newsType = model.NewsTypeNews
		)

		r := news.Create(controller.Context{
			Uid: adminInfo.Id,
		}, news.CreateNewParams{
			Title:   title,
			Content: content,
			Type:    newsType,
			Tags:    []string{},
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := schema.News{}

		assert.Nil(t, tester.Decode(r.Data, &n))

		defer news.DeleteNewsById(n.Id)
	}

	{
		r := tester.HttpUser.Get("/v1/news", nil, nil)

		res := schema.Response{}

		if !assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res)) {
			return
		}

		if !assert.Equal(t, schema.StatusSuccess, res.Status) {
			return
		}

		if !assert.Equal(t, "", res.Message) {
			return
		}

		list := make([]schema.News, 0)

		assert.Nil(t, tester.Decode(res.Data, &list))

		for _, b := range list {
			assert.IsType(t, "string", b.Title)
			assert.IsType(t, "string", b.Content)
			assert.IsType(t, "string", b.Author)
			assert.IsType(t, model.NewsTypeAnnouncement, b.Type)
			assert.IsType(t, []string{""}, b.Tags)
			assert.IsType(t, model.NewsStatusActive, b.Status)
			assert.IsType(t, "string", b.CreatedAt)
			assert.IsType(t, "string", b.UpdatedAt)
		}
	}
}
