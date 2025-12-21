package api

import (
	"net/http"

	"github.com/PRPO-skupina-02/common/middleware"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	_ "github.com/orgs/PRPO-skupina-02/Spored/api/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

//	@title			Spored API
//	@version		1.0
//	@description	API za upravljanje z kinodvoranami in njihovim sporedom

//	@host		localhost:8080
//	@BasePath	/api/v1

func Register(router *gin.Engine, db *gorm.DB, trans ut.Translator) {
	// Healthcheck
	router.GET("/healthcheck", healthcheck)

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// REST API
	v1 := router.Group("/api/v1")
	v1.Use(middleware.TransactionMiddleware(db))
	v1.Use(middleware.TranslationMiddleware(trans))
	v1.Use(middleware.ErrorMiddleware)

	// Theaters
	theaters := v1.Group("/theaters")
	theaters.GET("", TheatersList)
	theaters.GET("/:id", TheatersShow)
	theaters.POST("", TheatersCreate)
	theaters.PUT("/:id", TheatersUpdate)
	theaters.DELETE("/:id", TheatersDelete)

}

func healthcheck(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
