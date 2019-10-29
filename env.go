package main

import (
	"log"
	"os"
	"strconv"
)

type Env struct {
	S Storage
}

func getEnv() *Env {
	addr := os.Getenv("APP_REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	password := os.Getenv("APP_REDIS_PASSWD")
	if password == "" {
		password = ""
	}

	dbStr := os.Getenv("APP_REDIS_DB")
	if dbStr == "" {
		dbStr = "0"
	}

	db, err := strconv.Atoi(dbStr);
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("connect to redis (addr: %s, password: %s, db: %d)", addr, password, db)

	r := NewRedisClient(addr, password, db)

	return &Env{S: r}
}
