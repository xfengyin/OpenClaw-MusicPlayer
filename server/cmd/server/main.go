package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xfengyin/OpenClaw-MusicPlayer/server/handler"
	"github.com/xfengyin/OpenClaw-MusicPlayer/server/middleware"
)

func main() {
	r := gin.Default()

	// 使用中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	// 路由
	api := r.Group("/api")
	{
		// 用户相关
		api.POST("/login", gin.WrapF(handler.UserLogin))
		api.GET("/user/profile", gin.WrapF(handler.GetUserProfile))
		api.GET("/user/playlists", gin.WrapF(handler.GetUserPlaylists))

		// 播放相关
		api.POST("/playback/play", gin.WrapF(handler.PlayMusic))
		api.POST("/playback/pause", gin.WrapF(handler.PauseMusic))
		api.POST("/playback/next", gin.WrapF(handler.NextTrack))
		api.POST("/playback/prev", gin.WrapF(handler.PrevTrack))
		api.POST("/playback/volume", gin.WrapF(handler.SetVolume))

		// 下载相关
		api.POST("/download", gin.WrapF(handler.DownloadSong))
		api.GET("/download/status", gin.WrapF(handler.GetDownloadStatus))

		// 解析相关
		api.POST("/parse/playlist", gin.WrapF(handler.ParsePlaylist))
		api.POST("/parse/song", gin.WrapF(handler.ParseSong))

		// 歌词相关
		api.GET("/lyrics", gin.WrapF(handler.GetLyrics))

		// 远程控制
		api.POST("/remote/control", gin.WrapF(handler.RemoteControl))
		api.GET("/remote/status", gin.WrapF(handler.GetRemoteStatus))

		// 同步相关
		api.POST("/sync", gin.WrapF(handler.SyncLibrary))
		api.GET("/sync/status", gin.WrapF(handler.GetSyncStatus))
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	port := ":8080"
	log.Printf("Server starting on %s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
