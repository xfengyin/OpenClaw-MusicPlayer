package handler

import (
    "encoding/json"
    "net/http"
)

// RemoteControl 远程控制
func RemoteControl(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Command string          `json:"command"`
        Data    json.RawMessage `json:"data"`
    }
    
    json.NewDecoder(r.Body).Decode(&req)
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "code": 200,
        "msg":  "success",
    })
}

// GetRemoteStatus 获取远程状态
func GetRemoteStatus(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "code": 200,
        "msg":  "success",
        "data": map[string]interface{}{
            "connected": 0,
        },
    })
}
