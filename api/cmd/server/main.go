package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/parwin-pp/todo-application/internal/auth"
	"github.com/parwin-pp/todo-application/internal/config"
	"github.com/parwin-pp/todo-application/internal/middleware"
	"github.com/parwin-pp/todo-application/internal/postgres"
	"github.com/parwin-pp/todo-application/internal/todo"
	todotask "github.com/parwin-pp/todo-application/internal/todo_task"
	"github.com/rs/cors"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"
)

func main() {
	conf := config.NewConfig()
	isProduction := conf.App.Env == "production"

	database := postgres.GetConnection(conf.Database)
	defer database.Close()

	if err := database.Ping(); err != nil {
		log.Fatalf("could not ping DB: %v", err)
	}
	if !isProduction {
		database.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
			bundebug.FromEnv("BUNDEBUG"),
		))
	}
	db := postgres.NewDB(database)
	encrypter := auth.NewAuthExcrption("HS256", []byte("secret"), "1h")

	authServer := auth.NewServer(db, encrypter, conf.Auth)
	todoServer := todo.NewServer(db)
	taskServer := todotask.NewServer(db)

	requestLogger := reqlog.NewMiddleware(reqlog.WithEnabled(!isProduction))
	router := bunrouter.New(
		bunrouter.Use(requestLogger),
	)

	router.POST("/login", authServer.HandleLogin)

	authRouter := router.Use(middleware.NewAuthMiddleware(encrypter))
	authRouter.GET("/me", authServer.HandleGetMe)
	authRouter.POST("/logout", authServer.HandleLogout)
	authRouter.GET("/todos", todoServer.HandleGetTodos)
	authRouter.GET("/todos/:todoId", todoServer.HandleGetTodo)
	authRouter.POST("/todos", todoServer.HandleCreateTodo)
	authRouter.GET("/todos/:todoId/tasks", taskServer.HandleGetTasks)
	authRouter.POST("/todos/:todoId/tasks", taskServer.HandleCreateTask)
	authRouter.PATCH("/todos/:todoId/tasks/:taskId", taskServer.HandlePartialUpdateTask)
	authRouter.DELETE("/todos/:todoId/tasks/:taskId", taskServer.HandleDeleteTask)

	handler := http.Handler(router)
	// TODO: FOR TEST ONLY MUST BE CHANGE IN PRODUCTION
	handler = cors.New(cors.Options{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler(handler)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", conf.App.Port),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       60 * time.Second,
		Handler:           handler,
	}

	log.Printf("listening on 0.0.0.0:%d\n", conf.App.Port)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("could not listen on 0.0.0.0:%d: %v", conf.App.Port, err)
		}
	}()

	waitForShutdown(server)
}

func waitForShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("could not gracefully shutdown the server: %v", err)
	}
	log.Printf("server shutdown completed")
}
