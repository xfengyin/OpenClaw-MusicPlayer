package middleware

import (
    "log"
    "net/http"
    "time"
)

// LoggerToFile 日志中间件
func LoggerToFile() func(http.Handler) {
    return func(next http.Handler) {
        http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            log.Printf("[%s] %s %s", r.Method, r.URL.Path, r.RemoteAddr)
            
            next.ServeHTTP(w, r)
            
            log.Printf("[%s] Completed in %v", r.Method, time.Since(start))
        })
    }
}
