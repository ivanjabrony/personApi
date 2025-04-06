package app

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ivanjabrony/personApi/internal/client"
	"github.com/ivanjabrony/personApi/internal/client/client_impl"
	"github.com/ivanjabrony/personApi/internal/controller"
	"github.com/ivanjabrony/personApi/internal/repository"
	"github.com/ivanjabrony/personApi/internal/repository/pg"
	"github.com/ivanjabrony/personApi/internal/service"
	"github.com/ivanjabrony/personApi/internal/service/service_impl"
	"github.com/jmoiron/sqlx"
)

type App struct {
	Router *gin.Engine
	db     *sqlx.DB
}

func New(db *sqlx.DB) *App {
	repositories := initRepositories(db)
	clients := initClients()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: getLogLevel()}))
	services := initServices(repositories, clients, logger)

	router := controller.SetupRouter(
		logger,
		services.person,
	)

	return &App{
		Router: router,
		db:     db,
	}
}

func (a *App) Run(addr string) error {
	return a.Router.Run(addr)
}

type repositories struct {
	person repository.PersonRepository
}

type clients struct {
	ageClient         client.AgeClient
	genderClient      client.GenderClient
	nationalityClient client.NationalityClient
}

type services struct {
	person service.PersonService
}

func initRepositories(db *sqlx.DB) *repositories {
	return &repositories{
		person: pg.NewPgPersonRepository(db),
	}
}

func initClients() *clients {
	return &clients{
		ageClient:         client_impl.NewAgifyClient(),
		genderClient:      client_impl.NewGenderizeClient(),
		nationalityClient: client_impl.NewNationalityClient(),
	}
}

func initServices(r *repositories, cl *clients, logger *slog.Logger) *services {
	return &services{
		person: service_impl.NewPersonService(r.person, cl.ageClient, cl.genderClient, cl.nationalityClient, logger),
	}
}

func getLogLevel() slog.Level {
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
