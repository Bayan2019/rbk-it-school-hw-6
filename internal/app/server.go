package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/Bayan2019/rbk-it-school-hw-6/internal/auth"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/client"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/config"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/handler"
	middle "github.com/Bayan2019/rbk-it-school-hw-6/internal/middleware"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/repository/postgres"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/service"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

// Ch 2. Logging Lv 4. Global Logger vs. Dependency Injection
// Add a logger field to the server struct,
// and update server logging to use that injected logger.
type Server struct {
	httpServer *http.Server
	cancel     context.CancelFunc
	logger     *slog.Logger
}

func NewServer(
	// store store.Store,
	port int,
	cancel context.CancelFunc,
	accessLogger *slog.Logger,
	db *sqlx.DB,
) *Server {

	userRepo := postgres.NewUserRepository(db)
	cityRepo := postgres.NewCityRepository(db)
	weatherRepo := postgres.NewWeatherRepository(db)

	osmClient := client.NewOsmClient(config.Cfg.Api)
	weatherClient := client.NewWeatherClient()

	userService := service.NewUserService(userRepo)
	cityService := service.NewCityService(cityRepo, osmClient)
	weatherService := service.NewWeatherService(weatherRepo, weatherClient)

	jwtManager := auth.NewJWTManager([]byte(config.Cfg.App.JwtSecret), accessLogger)

	userHandler := handler.NewUserHandler(userService, jwtManager)
	cityHandler := handler.NewCityHandler(cityService)
	weatherHandler := handler.NewWeatherHandler(cityService, weatherService)

	r := chi.NewRouter()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: middle.RequestLogger(accessLogger)(r),
	}

	s := &Server{
		httpServer: srv,
		// store:      store,
		cancel: cancel,
		logger: accessLogger,
	}

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	// 1. Аутентификация
	userHandler.RegisterCommonRoutes(r)
	// 4. Защита маршрутов
	r.Group(func(r chi.Router) {
		// Все операции должны работать через текущего пользователя из JWT.
		r.Use(middle.AuthMiddleware(userHandler.JwtManager))
		// Убрать user_id из URL.
		r.Post("/cities", cityHandler.Add2User)
		r.Get("/cities", cityHandler.ListOfUser)
		r.Delete("/cities/{city_id}", cityHandler.DeleteFromUser)
		r.Get("/weather", weatherHandler.GetWeatherOfUserCities)
		r.Get("/weather/history", weatherHandler.GetWeatherHistoryOfUser)

		userHandler.RegisterAuthRoutes(r)

		// 5. Авторизация (Roles)
		r.Group(func(r chi.Router) {
			// Использовать middleware RequireRole("admin")
			r.Use(middle.RequireRole(auth.RolesAdmin))
			// Только admin может:
			userHandler.RegisterAdminRoutes(r)
			r.HandleFunc("POST /admin/shutdown", s.HandlerShutdown)
		})
	})

	return s
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		return err
	}

	// Ch 1. Observability Lv 3. What Is Observability?
	// When the server starts, print the following message to the console,
	// where %d is the port number:
	// ln.Addr() returns a net.Addr interface.
	if addr, ok := ln.Addr().(*net.TCPAddr); ok {
		httpPort := addr.Port
		s.logger.Info(fmt.Sprintf("WeatherAPI is running on http://localhost:%d\n", httpPort))
	}

	if err := s.httpServer.Serve(ln); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) HandlerShutdown(w http.ResponseWriter, r *http.Request) {
	// s.logger.Info("shutting down", "ENV", config.Cfg.App.Environment)
	if config.Cfg.App.Environment == "production" {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	go s.cancel()
}
