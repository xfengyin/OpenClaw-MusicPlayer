package handler

import (
	"encoding/json"
	"net/http"

	"github.com/xfengyin/OpenClaw-MusicPlayer/server/service"
)

var lyricsService = service.NewLyricsService()

// GetLyrics 获取歌词
func GetLyrics(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code": 400,
			"msg":  "missing id parameter",
		})
		return
	}

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
