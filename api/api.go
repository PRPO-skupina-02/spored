package api

import (
	"net/http"

	"github.com/PRPO-skupina-02/common/middleware"
	_ "github.com/PRPO-skupina-02/spored/api/docs"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

//	@title			Spored API
//	@version		1.0
//	@description	API za upravljanje z kinodvoranami in njihovim sporedom

//	@host		localhost:8080
//	@BasePath	/api/v1/spored

func Register(router *gin.Engine, db *gorm.DB, trans ut.Translator, authHost string) {
	// Healthcheck
	router.GET("/healthcheck", healthcheck)

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// REST API
	v1 := router.Group("/api/v1/spored")
	v1.Use(middleware.TransactionMiddleware(db))
	v1.Use(middleware.TranslationMiddleware(trans))
	v1.Use(middleware.ErrorMiddleware)

	// Theaters
	v1.GET("/theaters", TheatersList)
	theaters := v1.Group("/theaters/:theaterID")
	theaters.Use(TheaterContextMiddleware)
	theaters.GET("", TheatersShow)

	theatersAdmin := v1.Group("/theaters")
	theatersAdmin.Use(middleware.UserMiddleware(authHost))
	theatersAdmin.Use(middleware.RequireAdmin())
	theatersAdmin.POST("", TheatersCreate)

	theatersAdminWithID := theaters.Group("")
	theatersAdminWithID.Use(middleware.UserMiddleware(authHost))
	theatersAdminWithID.Use(middleware.RequireAdmin())
	theatersAdminWithID.PUT("", TheatersUpdate)
	theatersAdminWithID.DELETE("", TheatersDelete)

	// Rooms
	theaters.GET("/rooms", RoomsList)
	theaters.GET("/rooms/:roomID", RoomsShow)

	roomsAdmin := theaters.Group("/rooms")
	roomsAdmin.Use(middleware.UserMiddleware(authHost))
	roomsAdmin.Use(middleware.RequireAdmin())
	roomsAdmin.POST("", RoomsCreate)
	roomsAdmin.PUT("/:roomID", RoomsUpdate)
	roomsAdmin.DELETE("/:roomID", RoomsDelete)

	// Movies
	v1.GET("/movies", MoviesList)
	movies := v1.Group("/movies/:movieID")
	movies.Use(MovieContextMiddleware)
	movies.GET("", MoviesShow)

	moviesAdmin := v1.Group("/movies")
	moviesAdmin.Use(middleware.UserMiddleware(authHost))
	moviesAdmin.Use(middleware.RequireAdmin())
	moviesAdmin.POST("", MoviesCreate)

	moviesAdminWithID := movies.Group("")
	moviesAdminWithID.Use(middleware.UserMiddleware(authHost))
	moviesAdminWithID.Use(middleware.RequireAdmin())
	moviesAdminWithID.PUT("", MoviesUpdate)
	moviesAdminWithID.DELETE("", MoviesDelete)

	// TimeSlots
	theaters.GET("/rooms/:roomID/timeslots", TimeSlotsList)
	theaters.GET("/rooms/:roomID/timeslots/:timeSlotID", TimeSlotsShow)
}

func healthcheck(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
