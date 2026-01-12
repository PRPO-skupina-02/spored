package api

import (
	"testing"

	"github.com/PRPO-skupina-02/common/clients/auth/models"
	"github.com/PRPO-skupina-02/common/middleware"
	"github.com/PRPO-skupina-02/common/validation"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func MockAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := &models.APIUserResponse{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001").String(),
			Email:     "admin@example.com",
			FirstName: "Admin",
			LastName:  "User",
			Role:      models.ModelsUserRoleAdmin,
			Active:    true,
		}
		middleware.SetContextUser(c, user)
		c.Next()
	}
}

func TestingRouter(t *testing.T, db *gorm.DB) *gin.Engine {
	router := gin.Default()
	trans, err := validation.RegisterValidation()
	require.NoError(t, err)

	router.Use(MockAdminMiddleware())

	registerTestRoutes(router, db, trans)

	return router
}

func registerTestRoutes(router *gin.Engine, db *gorm.DB, trans ut.Translator) {
	// Healthcheck
	router.GET("/healthcheck", healthcheck)

	// REST API
	v1 := router.Group("/api/v1/spored")
	v1.Use(middleware.TransactionMiddleware(db))
	v1.Use(middleware.TranslationMiddleware(trans))
	v1.Use(middleware.ErrorMiddleware)

	// Public routes
	v1.GET("/theaters", TheatersList)
	v1.GET("/movies", MoviesList)

	// Theaters
	theaters := v1.Group("/theaters/:theaterID")
	theaters.Use(TheaterContextMiddleware)
	theaters.GET("", TheatersShow)
	v1.POST("/theaters", TheatersCreate)
	theaters.PUT("", TheatersUpdate)
	theaters.DELETE("", TheatersDelete)

	// Rooms
	theaters.GET("/rooms", RoomsList)
	theaters.GET("/rooms/:roomID", RoomsShow)
	theaters.POST("/rooms", RoomsCreate)
	theaters.PUT("/rooms/:roomID", RoomsUpdate)
	theaters.DELETE("/rooms/:roomID", RoomsDelete)

	// Movies
	movies := v1.Group("/movies/:movieID")
	movies.Use(MovieContextMiddleware)
	movies.GET("", MoviesShow)
	v1.POST("/movies", MoviesCreate)
	movies.PUT("", MoviesUpdate)
	movies.DELETE("", MoviesDelete)

	// TimeSlots
	theaters.GET("/rooms/:roomID/timeslots", TimeSlotsList)
	theaters.GET("/rooms/:roomID/timeslots/:timeSlotID", TimeSlotsShow)
}
