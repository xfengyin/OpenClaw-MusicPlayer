package server

import (
	"time"

	"github.com/gin-gonic/gin"
)

// setupRoutes 设置路由
func (s *Server) setupRoutes() {
	handler := NewHandler()

	// 健康检查
	s.router.GET("/health", s.handleHealth)

	// API v1 路由组
	api := s.router.Group("/api/v1")
	{
		// 音乐相关
		music := api.Group("/music")
		{
			music.GET("/search", handler.HandleSearch)
			music.GET("/detail/:id", handler.HandleGetSongDetail)
			music.GET("/url/:id", handler.HandleGetSongURL)
		}

		// 歌单相关
		playlist := api.Group("/playlist")
		{
			playlist.GET("/parse", handler.HandleParsePlaylist)
		}

		// 歌词相关
		lyrics := api.Group("/lyrics")
		{
			lyrics.GET("/:id", handler.HandleGetLyrics)
		}

		// 下载相关
		download := api.Group("/download")
		{
			download.POST("", handler.HandleDownload)
			download.GET("/progress/:taskId", handler.HandleGetDownloadProgress)
		}
	}
}

// handleHealth 健康检查
func (s *Server) handleHealth(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":    "ok",
		"version":   "1.0.0",
		"timestamp": timeNow().Unix(),
	})
}
