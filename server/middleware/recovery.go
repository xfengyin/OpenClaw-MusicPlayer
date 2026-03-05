package middleware

import (
    "net/http"
    "runtime/debug"
    "log"
)

// Recovery 恢复中间件
func Recovery() func(http.Handler) {
    return func(next http.Handler) {
        http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            defer func() {
                if err := recover(); err != nil {
                    log.Printf(" panicked: %v
%s", err, debug.Stack())
                    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                }
            }()
            next.ServeHTTP(w, r)
        })
    }
}
