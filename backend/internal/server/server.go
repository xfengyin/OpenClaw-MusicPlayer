package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xfengyin/OpenClaw-MusicPlayer/backend/internal/utils"
)

// Server HTTP服务器
type Server struct {
	router *gin.Engine
	srv    *http.Server
	logger *utils.Logger
}

// NewServer 创建新的服务器
func NewServer(logger *utils.Logger) *Server {
	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())
	router.Use(requestLogger(logger))

	s := &Server{
		router: router,
		logger: logger,
	}

	s.setupRoutes()
	return s
}

// setupRoutes 设置路由
func (s *Server) setupRoutes() {
	// 健康检查
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"time":   time.Now().Unix(),
		})
	})

	// API路由组
	api := s.router.Group("/api/v1")
	{
		// 音乐相关
		music := api.Group("/music")
		{
			music.GET("/search", s.handleSearch)
			music.GET("/detail/:id", s.handleDetail)
			music.GET("/url/:id", s.handleGetURL)
		}

		// 歌单相关
		playlist := api.Group("/playlist")
		{
			playlist.GET("/parse", s.handleParsePlaylist)
		}

		// 歌词相关
		lyrics := api.Group("/lyrics")
		{
			lyrics.GET("/:id", s.handleGetLyrics)
		}
	}
}

// Start 启动服务器
func (s *Server) Start(addr string) error {
	s.srv = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}
	return s.srv.ListenAndServe()
}

// Stop 停止服务器
func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.srv.Shutdown(ctx)
}

// corsMiddleware CORS中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// requestLogger 请求日志中间件
func requestLogger(logger *utils.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		logger.Info("[%s] %s %s %d %v", clientIP, method, path, statusCode, latency)
	}
}

// 处理器方法（占位符）
func (s *Server) handleSearch(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		c.JSON(400, gin.H{"error": "keyword is required"})
		return
	}
	c.JSON(200, gin.H{
		"keyword": keyword,
		"results": []interface{}{},
	})
}

func (s *Server) handleDetail(c *gin.Context) {
	id := c.Param("id")
	c.JSON(200, gin.H{"id": id})
}

func (s *Server) handleGetURL(c *gin.Context) {
	id := c.Param("id")
	c.JSON(200, gin.H{"id": id, "url": ""})
}

func (s *Server) handleParsePlaylist(c *gin.Context) {
	url := c.Query("url")
	c.JSON(200, gin.H{"url": url})
}

func (s *Server) handleGetLyrics(c *gin.Context) {
	id := c.Param("id")
	c.JSON(200, gin.H{"id": id, "lyrics": ""})
}
