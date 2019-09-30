// Copyright 2019 Axetroy. All rights reserved. MIT license.
package downloader

import (
	"fmt"
	"github.com/axetroy/go-fs"
	"github.com/axetroy/go-server/src/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
)

func Image(context *gin.Context) {
	filename := context.Param("filename")
	originImagePath := path.Join(config.Upload.Path, config.Upload.Image.Path, filename)
	if fs.PathExists(originImagePath) == false {
		// if the path not found
		http.NotFound(context.Writer, context.Request)
		return
	}
	context.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%v", filename))
	http.ServeFile(context.Writer, context.Request, originImagePath)
}
