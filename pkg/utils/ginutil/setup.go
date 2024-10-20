package ginutil

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"git.pmx.cn/hci/microservice-app/pkg/utils/httputil"
)

type Resource interface {
	Register(router *gin.RouterGroup)
}

func SetupResource(rg *gin.RouterGroup, resources ...Resource) {
	for _, resource := range resources {
		resource.Register(rg)
	}
}

func SetupPing(e *gin.Engine) {
	pingHandler := func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	}

	e.HEAD("/ping", pingHandler)
	e.GET("/ping", pingHandler)
}

func SetupSwagger(engine *gin.Engine) {
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func Startup(e *gin.Engine, addr string) {
	srv := &http.Server{
		Addr:    addr,
		Handler: e,
	}

	go func() {
		log.Printf("[rest server listen at %s]", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()

	httputil.SetupGracefulStop(srv)
}
