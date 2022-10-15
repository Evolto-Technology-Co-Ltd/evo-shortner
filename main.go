package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ariebrainware/evo-shortner/config"
	"github.com/ariebrainware/evo-shortner/endpoint"
	"github.com/ariebrainware/evo-shortner/model"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	// load app.env file data to struct
	config, err := config.LoadConfig(".")
	// handle errors
	if err != nil {
		log.Fatalf("can't load environment app.env: %v", err)
	}

	r := gin.Default()

	allowedOrigins := "*"
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowWildcard = true
	corsConfig.AllowOrigins = strings.Split(allowedOrigins, ",") // contain whitelist domain
	corsConfig.AllowHeaders = []string{"*", "Content-Type", "Accept"}
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowMethods("OPTIONS")
	r.Use(cors.New(corsConfig))

	r.GET("/", func(c *gin.Context) {
		var success string = fmt.Sprintf("Server listening with version %s", config.Version)
		c.JSON(http.StatusOK, &model.Response{
			Success: true,
			Error:   nil,
			Msg:     success,
			Data:    nil,
		})
	})
	r.POST("/short", endpoint.ShortURL)
	r.GET("/:key", endpoint.GetURL)

	port, _ := strconv.Atoi(config.Port)
	log.Infof("Service version: %s", config.Version)
	err = r.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Error(err)
	}
}
