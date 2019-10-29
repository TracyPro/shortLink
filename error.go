package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// 自定义错误接口 包括error接口
type Error interface {
	error
	Status() int
}

type StatusError struct {
	Code int   // 错误码
	Err  error // error接口变量 表示只要实现error接口的值，就可以赋值给这个变量
}

func (se StatusError) Error() string {
	return se.Err.Error()
}

func (se StatusError) Status() int {
	return se.Code
}

func respondWithError(w http.ResponseWriter, err error) {
	// 结构体实现了error接口和Error接口，此处可以进行判断
	switch e := err.(type) {
	case Error:
		log.Printf("HTTP %d - %s", e.Status(), e.Error())
		respondWithJSON(w, e.Status(), e.Error())
	default:
		respondWithJSON(w, http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError))
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// 序列化payload
	// 写入header
	resp, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(code)
	w.Write(resp)
}
