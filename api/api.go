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
	theaters := v1.Group("/theaters/:theaterID")
	theaters.Use(TheaterContextMiddleware)

	theatersRestricted := theaters.Group("")
	theatersRestricted.Use(TheaterPermissionsMiddleware)

	v1.GET("/theaters", TheatersList)
	theaters.GET("", TheatersShow)
	v1.POST("/theaters", TheatersCreate)

	theatersRestricted.PUT("", TheatersUpdate)
	theatersRestricted.DELETE("", TheatersDelete)

	// Rooms
	theaters.GET("/rooms", RoomsList)
	theaters.GET("/rooms/:roomID", RoomsShow)
	theatersRestricted.POST("/rooms", RoomsCreate)
	theatersRestricted.PUT("/rooms/:roomID", RoomsUpdate)
	theatersRestricted.DELETE("/rooms/:roomID", RoomsDelete)

}

func healthcheck(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
