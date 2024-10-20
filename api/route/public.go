package route

import "github.com/gin-gonic/gin"

func RegisterPublic(router *gin.RouterGroup) {

	r := router.Group("/v1")
	r.GET("/bootstrappers", Dev)
	r.GET("/advertisingspace", Dev)
}

func Dev(c *gin.Context) {}