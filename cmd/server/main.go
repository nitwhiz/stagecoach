package main

import (
	"crypto/sha512"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nitwhiz/stagecoach/pkg/config"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

func main() {
	if err := config.Load(); err != nil {
		log.Fatalln(err)
	}

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization, Origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 64 MiB
	r.MaxMultipartMemory = 64 << 20

	r.POST("/upload", func(ctx *gin.Context) {
		authHeader := strings.TrimSpace(ctx.GetHeader("Authorization"))
		tokenData := strings.SplitN(authHeader, " ", 2)

		if !strings.HasPrefix(authHeader, "Token ") || len(tokenData) < 2 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "missing auth token",
			})
			return
		}

		if fmt.Sprintf("%x", sha512.Sum512([]byte(tokenData[1]))) != config.C.AuthorizationToken {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "incorrect auth token",
			})
			return
		}

		file, err := ctx.FormFile("file")

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		name, ok := ctx.GetPostForm("name")

		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "missing file name",
			})
			return
		}

		suffix := 0
		ext := path.Ext(name)
		base := path.Join(config.C.DestinationDirectory, strings.TrimSuffix(name, ext))

		allowedFilePath := name
		lastError := error(nil)

		for suffix < 1000000 {
			if suffix > 0 {
				allowedFilePath = fmt.Sprintf("%s-%d%s", base, suffix, ext)
			} else {
				allowedFilePath = fmt.Sprintf("%s%s", base, ext)
			}

			_, err := os.Stat(allowedFilePath)

			if err == nil {
				suffix += 1
				continue
			}

			if os.IsNotExist(err) {
				break
			}

			// stat failed somehow we don't expect it to fail
			allowedFilePath = ""
			lastError = err
			break
		}

		if allowedFilePath == "" {
			if lastError != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": lastError.Error(),
				})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "unknown error",
				})
			}

			return
		}

		if err := ctx.SaveUploadedFile(file, allowedFilePath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"name": allowedFilePath,
		})
	})

	_ = r.Run("0.0.0.0:4444")
}
