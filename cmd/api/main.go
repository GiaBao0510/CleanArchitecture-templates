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

	"github.com/GiaBao0510/project-templates/internal/adapter/handler"
	pgRepo "github.com/GiaBao0510/project-templates/internal/adapter/repository/postgres"
	"github.com/GiaBao0510/project-templates/internal/infrastructure/config"
	"github.com/GiaBao0510/project-templates/internal/infrastructure/database"
	"github.com/GiaBao0510/project-templates/internal/infrastructure/logger"
	"github.com/GiaBao0510/project-templates/internal/infrastructure/router"
	"github.com/GiaBao0510/project-templates/internal/usecase"
)

func main() {
	// 1. Load config
	cfg := config.Load()

	// 2. Init logger
	appLogger, err := logger.NewLogger(cfg.App.LogLevel)
	if err != nil {
		log.Fatalf("Không thể khởi tạo logger: %v", err)
	}
	defer appLogger.Sync()

	// 3. Init database
	ctx := context.Background()
	dbPool, err := database.NewPostgresPool(ctx, cfg.Database)
	if err != nil {
		log.Fatalf("Không thể kết nối database: %v", err)
	}
	defer dbPool.Close()

	// 4. Dependency Injection (thủ công, từ trong ra ngoài)
	//    Domain (repository interface) ← Adapter (postgres impl) ← Use Case ← Handler
	userRepo := pgRepo.NewUserRepository(dbPool)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)

	// 5. Setup router
	r := router.NewRouter(appLogger, userHandler)

	// 6. Start HTTP server với graceful shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.App.Port,
		Handler: r,
	}

	go func() {
		fmt.Printf("🚀 Server đang chạy tại http://localhost:%s\n", cfg.App.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server lỗi: %v", err)
		}
	}()

	// 7. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\n⏳ Đang tắt server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown lỗi: %v", err)
	}
	fmt.Println("✅ Server đã tắt thành công")
}
