package controller

import (
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/ivanjabrony/personApi/docs"

	"github.com/gin-gonic/gin"
	"github.com/ivanjabrony/personApi/internal/controller/middleware"
	"github.com/ivanjabrony/personApi/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(logger *slog.Logger, personService service.PersonService) *gin.Engine {
	r := gin.Default()

	timeoutTime := os.Getenv("TIMEOUT_TIME")
	if timeoutTime == "" {
		timeoutTime = "3"
	}
	timeoutParsed, err := strconv.Atoi(timeoutTime)
	if err != nil {
		timeoutParsed = 3
	}

	r.Use(middleware.LoggerMiddleware(logger))
	r.Use(middleware.TimeoutMiddleware(time.Duration(timeoutParsed) * time.Second))

	personCotroller := NewPersonController(personService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	docs.SwaggerInfo.Host = "localhost:" + port
	docs.SwaggerInfo.BasePath = "/api"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := r.Group("/api/persons")

	api.POST("/", personCotroller.CreatePerson)
	api.PUT("/", personCotroller.UpdatePerson)
	api.GET("/:id", personCotroller.GetPerson)
	api.DELETE("/:id", personCotroller.DeletePersonById)
	api.GET("/", personCotroller.GetAllPersons)
	api.GET("/filtered", personCotroller.GetFilteredPesons)

	return r
}
