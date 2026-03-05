package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"openclaw-music-player/server/handler"
	"openclaw-music-player/server/middleware"
	"openclaw-music-player/server/service"
	"openclaw-music-player/server/utils"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	// 加载配置
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/openclaw/")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("警告: 未找到配置文件，使用默认配置: %v", err)
	}

	// 初始化日志
	logger := utils.NewLogger()
	logger.Info("OpenClaw Music Player 服务启动中...")

	// 初始化服务
	parserService := service.NewParserService()
	downloadService := service.NewDownloadService()
	lyricsService := service.NewLyricsService()
	syncService := service.NewSyncService()
	remoteService := service.NewRemoteService()

	// 初始化 Gin
	gin.SetMode(viper.GetString("server.mode"))
	r := gin.New()

	// 中间件
	r.Use(middleware.LoggerToFile())
	r.Use(middleware.CORS())
	r.Use(middleware.Recovery())

	// 配置路由组
	api := r.Group("/api/v1")
	{
		// 健康检查
		api.GET("/health", handler.HealthCheck)

		// 音乐解析
		parser := api.Group("/parser")
		{
			parser.GET("/search", handler.SearchMusic(parserService))
			parser.GET("/detail/:id", handler.GetMusicDetail(parserService))
			parser.GET("/url/:id", handler.GetMusicUrl(parserService))
			parser.GET("/lyrics/:id", handler.GetLyrics(lyricsService))
		}

		// 下载管理
		download := api.Group("/download")
		{
			download.POST("/start", handler.StartDownload(downloadService))
			download.GET("/status/:id", handler.GetDownloadStatus(downloadService))
			download.GET("/list", handler.ListDownloads(downloadService))
			download.POST("/cancel/:id", handler.CancelDownload(downloadService))
		}

		// 播放服务
		playback := api.Group("/playback")
		{
			playback.POST("/play", handler.PlayMusic)
			playback.POST("/pause", handler.PauseMusic)
			playback.POST("/next", handler.NextTrack)
			playback.POST("/prev", handler.PrevTrack)
			playback.POST("/volume/:level", handler.SetVolume)
		}

		// 用户和数据同步
		user := api.Group("/user")
		{
			user.POST("/login", handler.UserLogin(syncService))
			user.GET("/profile", handler.GetUserProfile(syncService))
			user.GET("/playlists", handler.GetUserPlaylists(syncService))
			user.POST("/sync", handler.SyncData(syncService))
		}

		// 远程控制
		remote := api.Group("/remote")
		{
			remote.POST("/control", handler.RemoteControl(remoteService))
			remote.GET("/status", handler.GetRemoteStatus(remoteService))
		}
	}

	// 启动服务
	port := viper.GetString("server.port")
	if port == "" {
		port = "8080"
	}

	address := fmt.Sprintf(":%s", port)
	if err := r.Run(address); err != nil {
		logger.Fatal("服务启动失败: ", err)
	}
}