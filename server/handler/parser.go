package handler

import (
	"encoding/json"
	"net/http"

	"openclaw-music-player/server/service"
)

// HealthCheck 健康检查
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"service": "openclaw-music-server",
	})
}

// SearchMusic 搜索音乐
func SearchMusic(parserService *service.ParserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyword := r.URL.Query().Get("keyword")
		page := r.URL.Query().Get("page")
		size := r.URL.Query().Get("size")

		if keyword == "" {
			http.Error(w, "keyword is required", http.StatusBadRequest)
			return
		}

		result, err := parserService.SearchMusic(keyword, 1, 20)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code": 500,
				"msg":  err.Error(),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code": 200,
			"msg":  "success",
			"data": result,
		})
	}
}

// GetMusicDetail 获取音乐详情
func GetMusicDetail(parserService *service.ParserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		
		detail, err := parserService.GetMusicDetail(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code": 500,
				"msg":  err.Error(),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code": 200,
			"msg":  "success",
			"data": detail,
		})
	}
}

// GetMusicUrl 获取音乐播放地址
func GetMusicUrl(parserService *service.ParserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		
		music, err := parserService.GetMusicUrl(id, "netease")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code": 500,
				"msg":  err.Error(),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code": 200,
			"msg":  "success",
			"data": music,
		})
	}
}

// GetLyrics 获取歌词
func GetLyrics(lyricsService *service.LyricsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		
		lyrics, err := lyricsService.GetLyricsForSong(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code": 500,
				"msg":  err.Error(),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code": 200,
			"msg":  "success",
			"data": lyrics,
		})
	}
}

// StartDownload 开始下载
func StartDownload(downloadService *service.DownloadService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			MusicID string `json:"music_id"`
			Title   string `json:"title"`
			URL     string `json:"url"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		task, err := downloadService.StartDownload(req.MusicID, req.Title, req.URL)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code": 500,
				"msg":  err.Error(),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code": 200,
			"msg":  "success",
			"data": task,
		})
	}
}

// GetDownloadStatus 获取下载状态
func GetDownloadStatus(downloadService *service.DownloadService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		
		task, err := downloadService.GetDownloadStatus(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code": 404,
				"msg":  err.Error(),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code": 200,
			"msg":  "success",
			"data": task,
		})
	}
}

// ListDownloads 列出下载
func ListDownloads(downloadService *service.DownloadService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tasks := downloadService.ListDownloads()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code": 200,
			"msg":  "success",
			"data": tasks,
		})
	}
}

// CancelDownload 取消下载
func CancelDownload(downloadService *service.DownloadService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		
		if err := downloadService.CancelDownload(id); err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code": 404,
				"msg":  err.Error(),
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code": 200,
			"msg":  "success",
		})
	}
}