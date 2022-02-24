package app

import (
	"context"
	"log"
	"myapp/internal/domain/data/controller"
	"myapp/internal/domain/data/service"
	"myapp/internal/domain/data/source"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	source := source.New()
	usecase := service.New(source)
	handler := controller.New(usecase)
	port := os.Getenv("PORT")
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: handler.Init(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	// максимальное время выполнения запросов после фиксации остановки приложения
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Print("Сервер остановлен") // корректно
}
