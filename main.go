package main

import (
	"cmd/internal/myweb"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("%w", err))
	}
	router := gin.Default()
	router.Use(myweb.CorsMiddleware())
	myweb.InitUserRouter(router)
	myweb.InitPostRouter(router)
	myweb.InitCommentRouter(router)
	router.Run(":" + fmt.Sprintf(viper.GetString("port")))
}
