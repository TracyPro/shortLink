package main

import (
	"log"
	"net/http"
	"time"
)

type Middleware struct {
}

// log middleware
// 记录请求消耗时间
func (m Middleware) LoggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next.ServeHTTP(w, r)
		endTime := time.Now()
		log.Printf("[%s] %q %v", r.Method, r.URL.String(), endTime.Sub(startTime))
	}
	return http.HandlerFunc(fn)
}

// 恢复Panic
func (m Middleware) RecoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recover from panic: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
