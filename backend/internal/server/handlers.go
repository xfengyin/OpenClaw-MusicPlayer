package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xfengyin/OpenClaw-MusicPlayer/backend/internal/parser"
)

// Handler HTTP处理器
type Handler struct {
	parserService *parser.ParserService
}

// NewHandler 创建处理器
func NewHandler() *Handler {
	return &Handler{
		parserService: parser.NewParserService(),
	}
}

// SearchRequest 搜索请求
type SearchRequest struct {
	Keyword string `form:"keyword" binding:"required"`
	Limit   int    `form:"limit,default=20"`
	Offset  int    `form:"offset,default=0"`
}

// SearchResponse 搜索响应
type SearchResponse struct {
	Code    int                 `json:"code"`
	Message string              `json:"message"`
	Data    *parser.SearchResult `json:"data"`
}

// HandleSearch 处理搜索请求
func (h *Handler) HandleSearch(c *gin.Context) {
	var req SearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, SearchResponse{
			Code:    400,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	result, err := h.parserService.Search(req.Keyword, req.Limit, req.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, SearchResponse{
			Code:    500,
			Message: "Search failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SearchResponse{
		Code:    200,
		Message: "success",
		Data:    result,
	})
}

// SongURLRequest 歌曲URL请求
type SongURLRequest struct {
	ID      string `uri:"id" binding:"required"`
	Quality string `form:"quality,default=standard"`
}

// SongURLResponse 歌曲URL响应
type SongURLResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		ID      string `json:"id"`
		URL     string `json:"url"`
		Quality string `json:"quality"`
	} `json:"data"`
}

// HandleGetSongURL 处理获取歌曲URL请求
func (h *Handler) HandleGetSongURL(c *gin.Context) {
	var req SongURLRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, SongURLResponse{
			Code:    400,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	// 解析音质
	qualityMap := map[string]int{
		"standard": 128000,
		"medium":   192000,
		"high":     320000,
		"lossless": 999000,
	}
	quality := qualityMap[c.DefaultQuery("quality", "standard")]

	url, err := h.parserService.GetSongURL("netease", req.ID, quality)
	if err != nil {
		c.JSON(http.StatusInternalServerError, SongURLResponse{
			Code:    500,
			Message: "Failed to get URL: " + err.Error(),
		})
		return
	}

	resp := SongURLResponse{
		Code:    200,
		Message: "success",
	}
	resp.Data.ID = req.ID
	resp.Data.URL = url
	resp.Data.Quality = c.DefaultQuery("quality", "standard")

	c.JSON(http.StatusOK, resp)
}

// LyricsResponse 歌词响应
type LyricsResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		ID         string `json:"id"`
		Lyrics     string `json:"lyrics"`
		Translated string `json:"translated,omitempty"`
	} `json:"data"`
}

// HandleGetLyrics 处理获取歌词请求
func (h *Handler) HandleGetLyrics(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, LyricsResponse{
			Code:    400,
			Message: "ID is required",
		})
		return
	}

	lyrics, translated, err := h.parserService.GetLyrics("netease", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, LyricsResponse{
			Code:    500,
			Message: "Failed to get lyrics: " + err.Error(),
		})
		return
	}

	resp := LyricsResponse{
		Code:    200,
		Message: "success",
	}
	resp.Data.ID = id
	resp.Data.Lyrics = lyrics
	resp.Data.Translated = translated

	c.JSON(http.StatusOK, resp)
}

// PlaylistParseRequest 歌单解析请求
type PlaylistParseRequest struct {
	URL string `form:"url" binding:"required"`
}

// PlaylistParseResponse 歌单解析响应
type PlaylistParseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		URL    string             `json:"url"`
		Title  string             `json:"title"`
		Songs  []parser.SongInfo  `json:"songs"`
	} `json:"data"`
}

// HandleParsePlaylist 处理歌单解析请求
func (h *Handler) HandleParsePlaylist(c *gin.Context) {
	var req PlaylistParseRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, PlaylistParseResponse{
			Code:    400,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	// TODO: 实现歌单解析逻辑
	c.JSON(http.StatusOK, PlaylistParseResponse{
		Code:    200,
		Message: "success",
		Data: struct {
			URL    string            `json:"url"`
			Title  string            `json:"title"`
			Songs  []parser.SongInfo `json:"songs"`
		}{
			URL:   req.URL,
			Title: "歌单",
			Songs: []parser.SongInfo{},
		},
	})
}

// SongDetailResponse 歌曲详情响应
type SongDetailResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    *parser.SongInfo `json:"data"`
}

// HandleGetSongDetail 处理获取歌曲详情请求
func (h *Handler) HandleGetSongDetail(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, SongDetailResponse{
			Code:    400,
			Message: "ID is required",
		})
		return
	}

	// TODO: 实现歌曲详情获取逻辑
	c.JSON(http.StatusOK, SongDetailResponse{
		Code:    200,
		Message: "success",
		Data: &parser.SongInfo{
			ID:   id,
			Name: "歌曲名",
		},
	})
}

// DownloadRequest 下载请求
type DownloadRequest struct {
	ID      string `json:"id" binding:"required"`
	Quality string `json:"quality,default=standard"`
}

// DownloadResponse 下载响应
type DownloadResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		TaskID string `json:"taskId"`
		Status string `json:"status"`
	} `json:"data"`
}

// HandleDownload 处理下载请求
func (h *Handler) HandleDownload(c *gin.Context) {
	var req DownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, DownloadResponse{
			Code:    400,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	// TODO: 实现下载逻辑
	c.JSON(http.StatusOK, DownloadResponse{
		Code:    200,
		Message: "success",
		Data: struct {
			TaskID string `json:"taskId"`
			Status string `json:"status"`
		}{
			TaskID: "task_" + req.ID,
			Status: "pending",
		},
	})
}

// DownloadProgressResponse 下载进度响应
type DownloadProgressResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		TaskID    string  `json:"taskId"`
		Status    string  `json:"status"`
		Progress  float64 `json:"progress"`
		Speed     string  `json:"speed"`
		TotalSize int64   `json:"totalSize"`
	} `json:"data"`
}

// HandleGetDownloadProgress 处理获取下载进度请求
func (h *Handler) HandleGetDownloadProgress(c *gin.Context) {
	taskID := c.Param("taskId")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, DownloadProgressResponse{
			Code:    400,
			Message: "Task ID is required",
		})
		return
	}

	// TODO: 实现获取下载进度逻辑
	c.JSON(http.StatusOK, DownloadProgressResponse{
		Code:    200,
		Message: "success",
		Data: struct {
			TaskID    string  `json:"taskId"`
			Status    string  `json:"status"`
			Progress  float64 `json:"progress"`
			Speed     string  `json:"speed"`
			TotalSize int64   `json:"totalSize"`
		}{
			TaskID:    taskID,
			Status:    "downloading",
			Progress:  50.0,
			Speed:     "1.5MB/s",
			TotalSize: 10240000,
		},
	})
}
