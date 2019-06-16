// Copyright 2019 Axetroy. All rights reserved. MIT license.
package model

import (
	"github.com/axetroy/go-server/src/util"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"time"
)

type NewsType string
type NewsStatus int

const (
	NewsTypeNews         NewsType = "news"         // 新闻资讯
	NewsTypeAnnouncement NewsType = "announcement" // 官方公告

	NewsStatusInActive NewsStatus = -1 // 未启用的状态
	NewsStatusActive                   // 启用的状态
)

var (
	NewsTypes = []NewsType{NewsTypeNews, NewsTypeAnnouncement}
)

func IsValidNewsType(t NewsType) bool {
	for _, v := range NewsTypes {
		if v == t {
			return true
		}
	}
	return false
}

type News struct {
	Id        string         `gorm:"primary_key;unique;not null;index;type:varchar(32)" json:"id"` // 新闻公告类ID
	Author    string         `gorm:"not null;index;type:varchar(32)" json:"author"`                // 公告的作者ID
	Title     string         `gorm:"not null;index;type:varchar(32)" json:"title"`                 // 公告标题
	Content   string         `gorm:"not null;type:text" json:"content"`                            // 公告内容
	Type      NewsType       `gorm:"not null;type:varchar(32)" json:"type"`                        // 公告类型
	Tags      pq.StringArray `gorm:"type:varchar(32)[]" json:"tags"`                               // 公告的标签
	Status    NewsStatus     `gorm:"not null;type:integer" json:"status"`                          // 公告状态
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (news *News) TableName() string {
	return "news"
}

func (news *News) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("id", util.GenerateId())
}
