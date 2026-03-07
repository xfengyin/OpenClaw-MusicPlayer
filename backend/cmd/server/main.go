package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/xfengyin/OpenClaw-MusicPlayer/backend/internal/server"
	"github.com/xfengyin/OpenClaw-MusicPlayer/backend/internal/utils"
)

func main() {
	// 初始化日志
	logger := utils.NewLogger()
	logger.Info("Starting OpenClaw Music Player Server...")

	// 创建HTTP服务器
	srv := server.NewServer(logger)

	// 启动服务器（非阻塞）
	go func() {
		if err := srv.Start(":8080"); err != nil {
			logger.Fatal("Server failed to start: " + err.Error())
		}
	}()

	logger.Info("Server started on http://localhost:8080")

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	if err := srv.Stop(); err != nil {
		logger.Error("Server forced to shutdown: " + err.Error())
	}
	logger.Info("Server exited")
}
