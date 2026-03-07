package middleware

import (
    "log"
    "net/http"
    "time"
)

// Logger 日志中间件
func Logger() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            log.Printf("[%s] %s %s", r.Method, r.URL.Path, r.RemoteAddr)

            next.ServeHTTP(w, r)

            log.Printf("[%s] Completed in %v", r.Method, time.Since(start))
        })
    }
}
