package handler

import (
	"encoding/json"
	"net/http"

	"openclaw-music-player/server/service"
)

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