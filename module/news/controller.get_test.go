// Copyright 2019 Axetroy. All rights reserved. MIT license.
package news_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/module/news"
	"github.com/axetroy/go-server/module/news/news_model"
	"github.com/axetroy/go-server/module/news/news_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetNews(t *testing.T) {
	// 获取一篇不存在的新闻公告
	{
		r := news.GetNews("123123")

		assert.Equal(t, schema.StatusFail, r.Status)
		assert.Equal(t, news.ErrNewsNotExist.Error(), r.Message)
	}

	// 获取一篇存在的新闻公告
	{
		var (
			newsId string
		)

		adminInfo, _ := tester.LoginAdmin()

		// 2. 先创建一篇新闻作为测试
		{
			var (
				title    = "test"
				content  = "test"
				newsType = news_model.NewsType_News
			)

			r := news.Create(schema.Context{
				Uid: adminInfo.Id,
			}, news.CreateNewParams{
				Title:   title,
				Content: content,
				Type:    newsType,
				Tags:    []string{},
			})

			assert.Equal(t, schema.StatusSuccess, r.Status)
			assert.Equal(t, "", r.Message)

			n := news_schema.News{}

			assert.Nil(t, tester.Decode(r.Data, &n))

			newsId = n.Id

			defer news.DeleteNewsById(n.Id)
		}

		// 3. 获取文章公告
		{
			r := news.GetNews(newsId)

			assert.Equal(t, schema.StatusSuccess, r.Status)
			assert.Equal(t, "", r.Message)

			newsInfo := news_schema.News{}

			assert.Nil(t, tester.Decode(r.Data, &newsInfo))

			assert.Equal(t, "test", newsInfo.Title)
			assert.Equal(t, "test", newsInfo.Content)
			assert.Equal(t, news_model.NewsType_News, newsInfo.Type)
		}
	}
}

func TestGetNewsRouter(t *testing.T) {
	var (
		newsId string
	)

	adminInfo, _ := tester.LoginAdmin()

	// 先创建一篇新闻作为测试
	{
		var (
			title    = "test"
			content  = "test"
			newsType = news_model.NewsType_News
		)

		r := news.Create(schema.Context{
			Uid: adminInfo.Id,
		}, news.CreateNewParams{
			Title:   title,
			Content: content,
			Type:    newsType,
			Tags:    []string{},
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		n := news_schema.News{}

		assert.Nil(t, tester.Decode(r.Data, &n))

		newsId = n.Id

		defer news.DeleteNewsById(n.Id)
	}

	// 获取详情
	{
		header := mocker.Header{}

		r := tester.HttpUser.Get("/v1/news/n/"+newsId, nil, &header)
		res := schema.Response{}

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

		n := news_schema.News{}

		assert.Nil(t, tester.Decode(res.Data, &n))
	}
}
