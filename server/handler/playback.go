package handler

import (
    "encoding/json"
    "net/http"
)

// PlayMusic 播放音乐
func PlayMusic(w http.ResponseWriter, r *http.Request) {
    var req struct {
        MusicID string `json:"music_id"`
    }
    
    json.NewDecoder(r.Body).Decode(&req)
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "code": 200,
        "msg":  "success",
    })
}

// PauseMusic 暂停播放
func PauseMusic(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "code": 200,
        "msg":  "success",
    })
}

// NextTrack 下一首
func NextTrack(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "code": 200,
        "msg":  "success",
    })
}

// PrevTrack 上一首
func PrevTrack(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "code": 200,
        "msg":  "success",
    })
}

// SetVolume 设置音量
func SetVolume(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "code": 200,
        "msg":  "success",
    })
}
