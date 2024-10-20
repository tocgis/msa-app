package session

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)
func init() {

}

func Start(r *gin.Engine){
	store, _ := redis.NewStoreWithDB(10, "tcp", "localhost:6379", "rcyj123456", "10",  []byte("aSecret"))
	r.Use(sessions.Sessions("ginSess", store))
}