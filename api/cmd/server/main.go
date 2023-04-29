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

	db := MustGetDBConnection(conf.Database, isProduction)
	defer db.Close()

	encrypter := MustGetEncrypter(conf.Auth)

	authServer := auth.NewServer(db, encrypter, conf.Auth)
	todoServer := todo.NewServer(db)
	taskServer := todotask.NewServer(db)

	requestLogger := reqlog.NewMiddleware(reqlog.WithEnabled(!isProduction))
	router := bunrouter.New(bunrouter.Use(requestLogger))

	// No auth required
	{
		router.POST("/login", authServer.HandleLogin)
		router.POST("/logout", authServer.HandleLogout)
	}

	// Auth required
	{
		authRouter := router.Use(middleware.NewAuthMiddleware(encrypter))
		authRouter.GET("/me", authServer.HandleGetMe)
		authRouter.GET("/todos", todoServer.HandleGetTodos)
		authRouter.GET("/todos/:todoId", todoServer.HandleGetTodo)
		authRouter.POST("/todos", todoServer.HandleCreateTodo)
		authRouter.GET("/todos/:todoId/tasks", taskServer.HandleGetTasks)
		authRouter.POST("/todos/:todoId/tasks", taskServer.HandleCreateTask)
		authRouter.PATCH("/todos/:todoId/tasks/:taskId", taskServer.HandlePartialUpdateTask)
		authRouter.DELETE("/todos/:todoId/tasks/:taskId", taskServer.HandleDeleteTask)
	}

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

	server := StartServer(conf.App.Port, handler)

	fmt.Println(WaitExitSignal())
	fmt.Println(Shutdown(server))
}

func Shutdown(server *http.Server) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return server.Shutdown(ctx)
}

func WaitExitSignal() os.Signal {
	ch := make(chan os.Signal, 3)
	signal.Notify(
		ch,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	return <-ch
}

func StartServer(port int, handler http.Handler) *http.Server {
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       60 * time.Second,
		Handler:           handler,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("ListenAndServe failed: %v", err)
		}
	}()
	return server
}

func MustGetEncrypter(conf config.AuthConfig) *auth.AuthEncryption {
	return auth.NewAuthEncryption(
		"HS256",
		[]byte(conf.SecretKey),
		conf.ExpireDuration,
	)
}

func MustGetDBConnection(conf config.DatabaseConfig, isProduction bool) *postgres.DB {
	database := postgres.GetConnection(conf)

	if err := database.Ping(); err != nil {
		log.Fatalf("could not ping DB: %v", err)
	}
	if !isProduction {
		database.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
			bundebug.FromEnv("BUNDEBUG"),
		))
	}
	return postgres.NewDB(database)
}
