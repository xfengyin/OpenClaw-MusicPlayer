package handler

import (
	"encoding/json"
	"net/http"

	"github.com/xfengyin/OpenClaw-MusicPlayer/server/service"
)

var syncService = service.NewSyncService()

// SyncLibrary 同步音乐库
func SyncLibrary(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID     string `json:"user_id"`
		PlaylistID string `json:"playlist_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := syncService.SyncData(req.UserID, req.PlaylistID); err != nil {
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
	})
}

// GetSyncStatus 获取同步状态
func GetSyncStatus(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code": 400,
			"msg":  "missing user_id parameter",
		})
		return
	}

	playlists, err := syncService.GetUserPlaylists(userID)
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
		"data": map[string]interface{}{
			"playlists": playlists,
		},
	})
}
