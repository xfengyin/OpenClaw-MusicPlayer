package handler

import (
	"encoding/json"
	"net/http"

	"github.com/xfengyin/OpenClaw-MusicPlayer/server/service"
)

var downloadService = service.NewDownloadService()

// DownloadSong 下载歌曲
func DownloadSong(w http.ResponseWriter, r *http.Request) {
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

// GetDownloadStatus 获取下载状态
func GetDownloadStatus(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code": 400,
			"msg":  "missing id parameter",
		})
		return
	}

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

// ListDownloads 列出下载
func ListDownloads(w http.ResponseWriter, r *http.Request) {
	tasks := downloadService.ListDownloads()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": 200,
		"msg":  "success",
		"data": tasks,
	})
}

// CancelDownload 取消下载
func CancelDownload(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code": 400,
			"msg":  "missing id parameter",
		})
		return
	}

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
