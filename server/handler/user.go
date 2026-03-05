package handler

import (
    "encoding/json"
    "net/http"
)

// UserLogin 用户登录
func UserLogin(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Phone    string `json:"phone"`
        Password string `json:"password"`
    }
    
    json.NewDecoder(r.Body).Decode(&req)
    
    // 实现登录逻辑
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "code": 200,
        "msg":  "登录成功",
        "data": map[string]interface{}{
            "token": "xxxx",
        },
    })
}

// GetUserProfile 获取用户信息
func GetUserProfile(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "code": 200,
        "msg":  "success",
        "data": map[string]interface{}{
            "nickname": "用户",
        },
    })
}

// GetUserPlaylists 获取用户歌单
func GetUserPlaylists(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "code": 200,
        "msg":  "success",
        "data": map[string]interface{}{
            "playlists": []string{},
        },
    })
}
