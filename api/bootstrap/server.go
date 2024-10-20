package bootstrap

import (
	"context"
	"fmt"
	"time"

	"git.pmx.cn/hci/microservice-app/pkg/redis"
	"git.pmx.cn/hci/microservice-app/pkg/session"

	"git.pmx.cn/hci/microservice-app/api"
	"git.pmx.cn/hci/microservice-app/dao"
	"git.pmx.cn/hci/microservice-app/pkg/micro"
	"git.pmx.cn/hci/microservice-app/pkg/utils/ginutil"
	"git.pmx.cn/hci/microservice-app/pkg/utils/gormutil"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A cloud disk base on the cloud service.",
	Run: func(cmd *cobra.Command, args []string) {
		serverRun()
	},
}

var ctx = context.Background()

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().Int("port", 8080, "server port")

	viper.BindPFlags(serverCmd.Flags())
}

func serverRun() {
	//gin.SetMode(gin.ReleaseMode)
	if viper.IsSet("installed") {
		gormutil.Init(
			gormutil.Config{
				Driver: viper.GetString("database.driver"),
				DSN:    viper.GetString("database.dsn"),
			},
			true,
		)
		// init db
		dao.InitDao()
	}

	rdb := redis.GetRedis(
		viper.GetString("redis.host"),
		viper.GetString("redis.pass"),
		viper.GetInt("redis.systemDb"),
	) // 连接redis 服务
	err := rdb.Set(ctx, "start", time.Now(), 0).Err()
	if err != nil {
		panic(err)
	}
	micro.Run() // 加载微服务

	g := gin.Default()
	session.Start(g) // 加载 session
	api.Register(g)

	addr := fmt.Sprintf(":%d", viper.GetInt("port"))
	ginutil.Startup(g, addr)
}
