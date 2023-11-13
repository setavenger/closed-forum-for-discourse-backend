package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

type Daemon struct {
	DB *gorm.DB
}

func RunServer(api *Daemon) {

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:3000", "http://localhost:3000", "*"},
		AllowMethods:     []string{"GET", "PUT", "POST", "PATCH"},
		AllowHeaders:     []string{"Origin", "Authorization", "content-type", "Content-Range"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		//AllowOriginFunc: func(origin string) bool {
		//	return origin == "https://github.com"
		//},
		MaxAge: 12 * time.Hour,
	}))

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	authorized := router.Group("/api", api.BasicAuth())
	{
		authorized.GET("/login", api.GetCheckToken)
		authorized.GET("/keystones", api.GetKeystones)
		authorized.GET("/keystone/:id", api.GetKeystone)
		authorized.GET("/reflections/:id", api.GetReflections)

		authorized.POST("/new-keystone", api.PostPublishKeystone)
		authorized.POST("/new-reflection", api.PostPublishReflection)
	}

	if err := router.Run(":8000"); err != nil {
		log.Fatal(err)
	}
}
