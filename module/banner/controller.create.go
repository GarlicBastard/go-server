// Copyright 2019 Axetroy. All rights reserved. MIT license.
package banner

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/axetroy/go-server/common_error"
	"github.com/axetroy/go-server/middleware"
	"github.com/axetroy/go-server/module/admin"
	"github.com/axetroy/go-server/module/admin/admin_model"
	"github.com/axetroy/go-server/module/banner/banner_model"
	"github.com/axetroy/go-server/module/banner/banner_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type CreateParams struct {
	Image       string                      `json:"image" valid:"required~请填写图片URL"` // 图片 URL
	Href        string                      `json:"href" valid:"required~请填写图片跳转链接"` // 图片跳转的 URL
	Platform    banner_model.BannerPlatform `json:"platform" valid:"required~请选择平台"` // 用于哪个平台, web/app
	Description *string                     `json:"description"`                     // Banner 描述
	Priority    *int                        `json:"priority"`                        // 优先级，用于排序
	Identifier  *string                     `json:"identifier"`                      // APP 跳转标识符
	FallbackUrl *string                     `json:"fallback_url"`                    // APP 跳转标识符的备选方案
}

func Create(context schema.Context, input CreateParams) (res schema.Response) {
	var (
		err          error
		data         banner_schema.Banner
		tx           *gorm.DB
		isValidInput bool
	)

	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = common_error.ErrUnknown
			}
		}

		if tx != nil {
			if err != nil {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}

		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
		}
	}()

	// 参数校验
	if isValidInput, err = govalidator.ValidateStruct(input); err != nil {
		return
	} else if isValidInput == false {
		err = common_error.ErrInvalidParams
		return
	}

	tx = database.Db.Begin()

	adminInfo := admin_model.Admin{
		Id: context.Uid,
	}

	if err = tx.First(&adminInfo).Error; err != nil {
		// 没有找到管理员
		if err == gorm.ErrRecordNotFound {
			err = admin.ErrAdminNotExist
		}
		return
	}

	if !adminInfo.IsSuper {
		err = admin.ErrAdminNotSuper
		return
	}

	if input.Platform == banner_model.BannerPlatformPc {
		// PC 端
	} else if input.Platform == banner_model.BannerPlatformApp {
		// 移动端
	} else {
		err = ErrBannerInvalidPlatform
		return
	}

	bannerInfo := banner_model.Banner{
		// require
		Image:    input.Image,
		Href:     input.Href,
		Platform: input.Platform,
		// optional
		Description: input.Description,
		Priority:    input.Priority,
		Identifier:  input.Identifier,
		FallbackUrl: input.FallbackUrl,
	}

	if err = tx.Create(&bannerInfo).Error; err != nil {
		return
	}

	if er := mapstructure.Decode(bannerInfo, &data.BannerPure); er != nil {
		err = er
		return
	}

	data.CreatedAt = bannerInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = bannerInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func CreateRouter(ctx *gin.Context) {
	var (
		input CreateParams
		err   error
		res   = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		ctx.JSON(http.StatusOK, res)
	}()

	if err = ctx.ShouldBindJSON(&input); err != nil {
		err = common_error.ErrInvalidParams
		return
	}

	res = Create(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, input)
}
