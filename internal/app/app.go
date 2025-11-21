package app

import (
	"log"

	"TestTaskInAvito/internal/config"
	"TestTaskInAvito/internal/infra/db"
	"TestTaskInAvito/internal/infra/tx"
	"TestTaskInAvito/internal/repo/postgres"
	"TestTaskInAvito/internal/service"
	transport "TestTaskInAvito/internal/transport/http"

	stdhttp "net/http"
)

type App struct {
	server *stdhttp.Server
}

func New(cfg config.Config) (*App, error) {
	// --- DB ---
	database, err := db.NewPostgres(cfg.DBDSN)
	if err != nil {
		return nil, err
	}

	// --- Tx manager + CtxGetter ---
	txManager := tx.NewManager(database)
	ctxGetter := tx.NewCtxGetter()

	// --- Repositories ---
	teamRepo := postgres.NewTeamRepository(database, ctxGetter)
	userRepo := postgres.NewUserRepository(database, ctxGetter)
	prRepo := postgres.NewPRRepository(database, ctxGetter)

	// --- Services ---
	teamSvc := service.NewTeamService(teamRepo, userRepo, txManager)
	userSvc := service.NewUserService(userRepo, prRepo, txManager)
	prSvc := service.NewPRService(prRepo, userRepo, teamRepo, txManager)

	// --- HTTP router ---
	router := transport.NewRouter(teamSvc, userSvc, prSvc)

	server := &stdhttp.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: router,
	}

	return &App{
		server: server,
	}, nil
}

func (a *App) Run() error {
	log.Printf("HTTP server listening on %s", a.server.Addr)
	return a.server.ListenAndServe()
}
